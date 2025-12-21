package responses

import (
	. "github.com/radityprtama/proxygate/v6/internal/constant"
	"github.com/radityprtama/proxygate/v6/internal/interfaces"
	"github.com/radityprtama/proxygate/v6/internal/translator/translator"
)

func init() {
	translator.Register(
		OpenaiResponse,
		Gemini,
		ConvertOpenAIResponsesRequestToGemini,
		interfaces.TranslateResponse{
			Stream:    ConvertGeminiResponseToOpenAIResponses,
			NonStream: ConvertGeminiResponseToOpenAIResponsesNonStream,
		},
	)
}
