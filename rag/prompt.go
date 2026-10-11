package rag

import (
	"fmt"
	"rag-agent/vector"
	"strings"
)

const (
	PREAMBLE = `Use the following excerpts from the document collection to answer the question.
	Cite sources by file name when you draw from them. If the excerpts do not address the question, say so before before
	answering from general knowledge.`
	UNKNOWN_SOURCE = "(unknown source)"
)

// takes the retrieves chuncks from the vector store and formatts from for the llm to receive
func formatContext(hits []vector.Result) string {
	if len(hits) == 0 {
		return ""
	}
	var stringBuilder strings.Builder
	stringBuilder.WriteString(PREAMBLE)
	stringBuilder.WriteString("\n\n--- Excerpts ---\n\n")
	for i, h := range hits {
		source := h.Metadata["source"]
		if source == "" {
			source = UNKNOWN_SOURCE
		}
		fmt.Fprintf(&stringBuilder, "[%d] Source: %s (similarity %.2f) \n%s\n\n",
			i+1, source, h.Score, h.Content)
	}
	return strings.TrimSpace(stringBuilder.String())
}
