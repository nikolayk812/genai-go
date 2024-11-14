package main

import (
	"strings"
	"testing"
)

func Test1(t *testing.T) {
	chatModel, err := buildChatModel()
	if err != nil {
		t.Fatalf("build chat model: %s", err)
	}

	t.Run("straight-answer", func(t *testing.T) {
		answer, err := straightAnswer(chatModel)
		if err != nil {
			t.Fatalf("straight chat: %s", err)
		}

		if !strings.Contains(answer, "cloud.logs.verbose = true") {
			t.Fatalf("straight chat: %s", answer)
		}
	})

	t.Run("ragged-answer", func(t *testing.T) {
		answer, err := raggedAnswer(chatModel)
		if err != nil {
			t.Fatalf("straight chat: %s", err)
		}

		if !strings.Contains(answer, "cloud.logs.verbose = true") {
			t.Fatalf("ragged chat: %s", answer)
		}
	})
}
