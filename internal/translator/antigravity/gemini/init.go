package gemini

import (
	. "github.com/radityprtama/proxygate/v6/internal/constant"
	"github.com/radityprtama/proxygate/v6/internal/interfaces"
	"github.com/radityprtama/proxygate/v6/internal/translator/translator"
)

func init() {
	translator.Register(
		Gemini,
		Antigravity,
		ConvertGeminiRequestToAntigravity,
		interfaces.TranslateResponse{
			Stream:     ConvertAntigravityResponseToGemini,
			NonStream:  ConvertAntigravityResponseToGeminiNonStream,
			TokenCount: GeminiTokenCount,
		},
	)
}
