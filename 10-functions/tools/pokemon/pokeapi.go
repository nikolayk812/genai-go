package pokemon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/tools"
)

// pokemonResponse is the struct that represents the response from the PokeAPI.
// We are only interested in the id, name, moves and types.
type pokemonResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Moves []struct {
		Move struct {
			Name string `json:"name"`
		} `json:"move"`
	} `json:"moves"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

// Tool is an implementation of the tool interface that finds information using PokeAPI (https://pokeapi.co/).
type Tool struct {
	callbacksHandler callbacks.Handler
}

// Implement the tools.Tool interface.
var _ tools.Tool = Tool{}

// New creates a new PokeAPI tool.
func New() Tool {
	return Tool{
		callbacksHandler: callbacks.LogHandler{},
	}
}

// Name returns the name of the tool. This is used by the LLM to identify the tool.
func (t Tool) Name() string {
	return "PokeAPI"
}

// Description returns a description of the tool. This is used by the LLM to understand what the tool does.
func (t Tool) Description() string {
	return `A wrapper around PokeAPI. 
	Useful for when you need to answer general questions about pokemons. 
	Input should be a pokemon name in lowercase, without quotes.`
}

// Call sums up the numbers in the input and returns the result.
func (t Tool) Call(ctx context.Context, input string) (string, error) {
	if t.callbacksHandler != nil {
		t.callbacksHandler.HandleToolStart(ctx, input)
	}

	result, err := FetchAPI(input)
	if err != nil {
		if t.callbacksHandler != nil {
			t.callbacksHandler.HandleToolError(ctx, err)
		}
		return "", err
	}

	if t.callbacksHandler != nil {
		t.callbacksHandler.HandleToolEnd(ctx, result)
	}

	return result, nil
}

// FetchAPI fetches the pokemon information from PokeAPI. It returns a string with the pokemon information,
// including the ID, the number of moves, the moves and the types.
func FetchAPI(pokemon string) (s string, err error) {
	ctx := context.Background()

	baseApiUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", strings.ToLower(pokemon))

	req, err := http.NewRequestWithContext(ctx, "GET", baseApiUrl, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "pokemon-tool")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return "", fmt.Errorf("copying data in pokeapi: %w", err)
	}

	var p pokemonResponse
	err = json.Unmarshal(buf.Bytes(), &p)
	if err != nil {
		return "", fmt.Errorf("unmarshalling data in pokeapi: %w", err)
	}
	defer resp.Body.Close()

	var moveNames []string
	for _, m := range p.Moves {
		moveNames = append(moveNames, m.Move.Name)
	}

	var typeNames []string
	for _, t := range p.Types {
		typeNames = append(typeNames, t.Type.Name)
	}

	return fmt.Sprintf("ID: %d, MovesCount: %d, Moves: [%s], Types: [%s]", p.Id, len(moveNames), strings.Join(moveNames, ", "), strings.Join(typeNames, ", ")), nil
}
