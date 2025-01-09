package main

// This file uses https://github.com/tmc/langchaingo/tree/main/examples/ollama-functions-example as a reference.

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"

	"github.com/tmc/langchaingo/llms"

	"github.com/mdelapenya/genai-testcontainers-go/functions/tools/pokemon"
)

// functions is a list of tools that the model can use.
var functions = []llms.FunctionDefinition{
	{
		Name: "fetchPokeAPI",
		Description: `A wrapper around PokeAPI. 
	Useful for when you need to answer general questions about pokemons. 
	Input should be a pokemon name in lowercase, without quotes.`,
		Parameters: json.RawMessage(`{
			"type": "object", 
			"properties": {
				"pokemon": {"type": "string", "description": "The pokemon name in lowercase, without quotes. E.g. pikachu"}
			}, 
			"required": ["pokemon"]
		}`),
	},
	{
		// I found that providing a tool for Ollama to give the final response significantly
		// increases the chances of success.
		Name:        "finalResponse",
		Description: "Provide the final response to the user query",
		Parameters: json.RawMessage(`{
			"type": "object", 
			"properties": {
				"response": {"type": "string", "description": "The final response to the user query"}
			}, 
			"required": ["response"]
		}`),
	},
}

// Call is a struct that represents a tool call.
type Call struct {
	Tool  string         `json:"tool"`
	Input map[string]any `json:"tool_input"`
}

// unmarshalCall unmarshals the input string into a list of Call objects.
func unmarshalCall(input string) []Call {
	var calls []Call
	if err := json.Unmarshal([]byte(input), &calls); err != nil {
		return nil
	}

	return calls
}

// dispatchCall dispatches a call to the appropriate tool.
// It returns a message and a boolean indicating if the model should try again.
func dispatchCall(c *Call) (llms.MessageContent, bool) {
	// ollama doesn't always respond with a *valid* function call. As we're using prompt
	// engineering to inject the tools, it may hallucinate.
	if !validTool(c.Tool) {
		log.Printf("invalid function call: %#v, prompting model to try again", c)
		return llms.TextParts(llms.ChatMessageTypeHuman,
			"Tool does not exist, please try again."), true
	}

	// we could make this more dynamic, by parsing the function schema.
	switch c.Tool {
	case "fetchPokeAPI":
		pokemonName, ok := c.Input["pokemon"].(string)
		if !ok {
			log.Fatal("invalid input")
		}

		p, err := pokemon.FetchAPI(pokemonName)
		if err != nil {
			log.Fatal(err)
		}

		return llms.TextParts(llms.ChatMessageTypeTool, p), true
	case "finalResponse":
		resp, ok := c.Input["response"].(string)
		if !ok {
			log.Fatal("invalid input")
		}

		return llms.TextParts(llms.ChatMessageTypeHuman, resp), false
	default:
		// we already checked above if we had a valid tool.
		panic("unreachable")
	}
}

// systemMessage returns a system message with the available tools using a tool schema.
// It also defines how to call the tools based on that schema.
func systemMessage() string {
	bs, err := json.Marshal(functions)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(`You have access to the following tools:

%s

To use a tool, respond with a JSON object with the following structure:
[
	{
		"tool": <name of the called tool>,
		"tool_input": <parameters for the tool matching the above JSON schema>
	}
]

If you need to call multiple tools, append more tool objects to the array.
The tool input should be a JSON object with the following structure:
- when calling fetchPokeAPI: {"pokemon": <pokemon name>}
- when calling finalResponse, include the reasons in the final response: {"response": <final response to the user query including the reasons>}
`, string(bs))
}

// validTool checks if a tool name is valid.
func validTool(name string) bool {
	var valid []string
	for _, v := range functions {
		valid = append(valid, v.Name)
	}
	return slices.Contains(valid, name)
}
