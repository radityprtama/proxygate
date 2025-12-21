package claude

import (
	. "github.com/radityprtama/proxygate/v6/internal/constant"
	"github.com/radityprtama/proxygate/v6/internal/interfaces"
	"github.com/radityprtama/proxygate/v6/internal/translator/translator"
)

func init() {
	translator.Register(
		Claude,
		OpenAI,
		ConvertClaudeRequestToOpenAI,
		interfaces.TranslateResponse{
			Stream:     ConvertOpenAIResponseToClaude,
			NonStream:  ConvertOpenAIResponseToClaudeNonStream,
			TokenCount: ClaudeTokenCount,
		},
	)
}
