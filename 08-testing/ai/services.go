package ai

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

type Chatter interface {
	Chat(userMessage string) (string, error)
}

type Summarizer interface {
	Summarize(message string) error
}

type Validator interface {
	Validate(question string, answer string, reference string) (string, error)
}

type ChatService struct {
	systemMessage string
	chatModel     llms.Model
	ragCtx        *ragContext
}

type ragContext struct {
	relevantDocs []schema.Document
}

// ChatServiceOption is a functional option for ChatService
type ChatServiceOption func(*ChatService)

// WithStore sets the VectorStore for the ChatService
func WithRAGContext(docs []schema.Document) ChatServiceOption {
	return func(s *ChatService) {
		s.ragCtx = &ragContext{
			relevantDocs: docs,
		}
	}
}

func NewChat(model llms.Model, opts ...ChatServiceOption) *ChatService {
	system := `You are a helpful assistant.
Your task is to answer questions by providing clear and concise answers.

Follow these instructions:
- Your answer should be clear and concise, maximum 3-4 sentences
- If you do not know the answer, you can say so
- Use the information provided to answer, do not make up information
- Important: Do not mention that you have been provided with additional information or documents`

	cs := &ChatService{
		systemMessage: system,
		chatModel:     model,
	}

	for _, opt := range opts {
		opt(cs)
	}

	return cs
}

func (s *ChatService) Chat(userMessage string) (string, error) {
	ctx := context.Background()
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, s.systemMessage),
	}

	if s.ragCtx != nil {
		for _, doc := range s.ragCtx.relevantDocs {
			content = append(content, llms.TextParts(llms.ChatMessageTypeSystem, doc.PageContent))
		}
	}

	content = append(content, llms.TextParts(llms.ChatMessageTypeHuman, userMessage))

	completion, err := s.chatModel.GenerateContent(
		ctx, content,
		llms.WithTemperature(0.00),
		llms.WithTopK(1),
		llms.WithSeed(42),
	)
	if err != nil {
		return "", fmt.Errorf("llm generate content: %w", err)
	}

	response := ""
	for _, choice := range completion.Choices {
		response += choice.Content
	}

	return response, nil
}

type ValidatorAgent struct {
	systemMessage string
	chatModel     llms.Model
	userMessage   string
}

func (v *ValidatorAgent) Validate(question string, answer string, reference string) (string, error) {
	ctx := context.Background()
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, v.systemMessage),
		llms.TextParts(llms.ChatMessageTypeHuman, fmt.Sprintf(v.userMessage, question, answer, reference)),
	}

	completion, err := v.chatModel.GenerateContent(
		ctx, content,
		llms.WithTemperature(0.00),
		llms.WithTopK(1),
		llms.WithSeed(42),
	)
	if err != nil {
		return "", fmt.Errorf("llm generate content: %w", err)
	}

	response := ""
	for _, choice := range completion.Choices {
		response += choice.Content
	}

	return response, nil
}

func NewValidatorAgent(model llms.Model) *ValidatorAgent {
	system := `
### Instructions
You are a strict validator.
You will be provided with a question, an answer, and a reference.
Your task is to validate whether the answer is correct for the given question, based on the reference.

Follow these instructions:
- Respond only 'yes', 'no' or 'unsure' and always include the reason for your response
- Respond with 'yes' if the answer is correct
- Respond with 'no' if the answer is incorrect
- If you are unsure, simply respond with 'unsure'
- Respond with 'no' if the answer is not clear or concise
- Respond with 'no' if the answer is not based on the reference

Your response must be a json object with the following structure:
{
	"response": "yes",
	"reason": "The answer is correct because it is based on the reference provided."
}

### Example
Question: Is Madrid the capital of Spain?
Answer: No, it's Barcelona.
Reference: The capital of Spain is Madrid
###
Response: {
	"response": "no",
	"reason": "The answer is incorrect because the reference states that the capital of Spain is Madrid."
}
""")`

	user := `
###
Question: %s
###
Answer: %s
###
Reference: %s
###
`

	v := &ValidatorAgent{
		chatModel:     model,
		systemMessage: system,
		userMessage:   user,
	}

	return v
}
