package llms

import (
	lc "github.com/tmc/langchaingo/llms"
)

func ContentResponseToStrings(response *lc.ContentResponse) []string {
	if response == nil {
		return nil
	}

	result := make([]string, 0, len(response.Choices))

	for _, choice := range response.Choices {
		if choice == nil {
			continue
		}
		result = append(result, choice.Content)
	}

	return result
}

func TotalTokens(response *lc.ContentResponse) int {
	if response == nil {
		return 0
	}

	result := 0

	for _, choice := range response.Choices {
		if choice == nil {
			continue
		}

		if totalTokens, ok := choice.GenerationInfo["TotalTokens"].(int); ok {
			result += totalTokens
		}
	}

	return result
}
