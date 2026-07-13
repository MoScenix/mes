package service

import (
	"context"
	"encoding/json"
	"fmt"

	lapp "github.com/MoScenix/mes/app/bff/hertz_gen/bff/app"
	"github.com/cloudwego/hertz/pkg/app"
)

type AnswerAIService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAnswerAIService(Context context.Context, RequestContext *app.RequestContext) *AnswerAIService {
	return &AnswerAIService{RequestContext: RequestContext, Context: Context}
}

func (h *AnswerAIService) Run(req *lapp.AIControlRequest) (resp *lapp.BaseResponseBoolean, err error) {
	appID, answers, err := h.answerRequest(req)
	if err != nil {
		return &lapp.BaseResponseBoolean{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	submitted, err := answerAIQuestion(h.Context, appID, answers)
	if err != nil {
		return &lapp.BaseResponseBoolean{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	return &lapp.BaseResponseBoolean{
		Code:    0,
		Data:    submitted,
		Message: "success",
	}, nil
}

func (h *AnswerAIService) answerRequest(req *lapp.AIControlRequest) (int64, map[string]aiAnswerPayload, error) {
	var body struct {
		AppID   int64                      `json:"appId"`
		Answers map[string]aiAnswerPayload `json:"answers"`
	}
	if h.RequestContext != nil {
		data := h.RequestContext.Request.Body()
		if len(data) > 0 {
			if err := json.Unmarshal(data, &body); err != nil {
				return 0, nil, err
			}
		}
	}
	if body.AppID <= 0 && req != nil {
		body.AppID = req.GetAppId()
	}
	if len(body.Answers) == 0 && req != nil {
		answers, err := generatedAnswers(req.GetAnswers())
		if err != nil {
			return 0, nil, err
		}
		body.Answers = answers
	}
	if body.AppID <= 0 {
		return 0, nil, fmt.Errorf("appId is required")
	}
	if len(body.Answers) == 0 {
		return 0, nil, fmt.Errorf("answers are required")
	}
	return body.AppID, body.Answers, nil
}

func generatedAnswers(values map[string]*lapp.AIAnswer) (map[string]aiAnswerPayload, error) {
	if len(values) == 0 {
		return nil, nil
	}
	answers := make(map[string]aiAnswerPayload, len(values))
	for id, value := range values {
		if value == nil {
			continue
		}
		answer := aiAnswerPayload{Content: value.GetContent()}
		if value.GetPayloadJson() != "" {
			var payload map[string]any
			if err := json.Unmarshal([]byte(value.GetPayloadJson()), &payload); err != nil {
				return nil, fmt.Errorf("invalid answer payloadJson for %s: %w", id, err)
			}
			answer.Payload = payload
		}
		answers[id] = answer
	}
	return answers, nil
}
