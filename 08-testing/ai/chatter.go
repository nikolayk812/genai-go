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

// ragContext is the context for the RAG in the form of a list of relevant documents
type ragContext struct {
	relevantDocs []schema.Document
}

type ChatService struct {
	systemMessage string
	chatModel     llms.Model
	ragCtx        *ragContext
}

// ChatServiceOption is a functional option for ChatService
type ChatServiceOption func(*ChatService)

// WithRAGContext sets the relevant documents for the RAG
func WithRAGContext(docs []schema.Document) ChatServiceOption {
	return func(s *ChatService) {
		s.ragCtx = &ragContext{
			relevantDocs: docs,
		}
	}
}

// NewChat creates a new ChatService.
// It defines a default system message for the chat model to answer questions
// in a structured way:
// - Your answer should be clear and concise, maximum 3-4 sentences
// - If you do not know the answer, you can say so
// - Use the information provided to answer, do not make up information
// - Important: Do not mention that you have been provided with additional information or documents
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

// Chat creates a chat response from the user message.
// If there is a RAG context in the form of relevant documents, it will be added to the prompt
// as system messages.
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
