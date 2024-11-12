package ai

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

type Chatter interface {
	Chat(userMessage string) error
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

func New(model llms.Model, opts ...ChatServiceOption) *ChatService {
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

func (s *ChatService) Chat(userMessage string) error {
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

	_, err := s.chatModel.GenerateContent(
		ctx, content,
		llms.WithTemperature(0.00),
		llms.WithTopK(1),
		llms.WithSeed(42),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("llm generate content: %w", err)
	}

	return nil
}
