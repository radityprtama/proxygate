// Package openai provides translation between OpenAI Chat Completions and Kiro formats.
package openai

import (
	. "github.com/radityprtama/proxygate/v6/internal/constant"
	"github.com/radityprtama/proxygate/v6/internal/interfaces"
	"github.com/radityprtama/proxygate/v6/internal/translator/translator"
)

func init() {
	translator.Register(
		OpenAI, // source format
		Kiro,   // target format
		ConvertOpenAIRequestToKiro,
		interfaces.TranslateResponse{
			Stream:    ConvertKiroStreamToOpenAI,
			NonStream: ConvertKiroNonStreamToOpenAI,
		},
	)
}