package agent

import (
	"github.com/MoScenix/mes/app/ai/conf"
	"github.com/cloudwego/eino/adk"
)

func modelRetryConfig() *adk.ModelRetryConfig {
	maxRetries := conf.GetConf().LLM.MaxRetries
	if maxRetries <= 0 {
		return nil
	}
	return &adk.ModelRetryConfig{MaxRetries: maxRetries}
}
