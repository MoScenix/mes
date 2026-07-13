package designer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/MoScenix/mes/app/ai/agent"
	"github.com/MoScenix/mes/common/aievent"
)

const answerWait = 60 * time.Second

type answerEvent struct {
	TargetID string
	Answer   agent.AssistantAnswer
}

func waitAnswer(ctx context.Context, answers <-chan answerEvent, targetID string) (agent.AssistantAnswer, bool, error) {
	if answers == nil || targetID == "" {
		return agent.AssistantAnswer{}, false, nil
	}

	timer := time.NewTimer(answerWait)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return agent.AssistantAnswer{}, false, ctx.Err()
		case <-timer.C:
			return agent.AssistantAnswer{}, false, nil
		case event := <-answers:
			if answer, ok := answerForTarget(event.Answer, targetID); ok {
				return answer, true, nil
			}
		}
	}
}

func agentAnswer(event aievent.TaskEvent) agent.AssistantAnswer {
	return agent.AssistantAnswer{
		Content: event.Content,
		Payload: event.Payload,
		Answers: parseAnswerMap(event.Payload),
	}
}

func answerForTarget(answer agent.AssistantAnswer, targetID string) (agent.AssistantAnswer, bool) {
	if targetID == "" || len(answer.Answers) == 0 {
		return agent.AssistantAnswer{}, false
	}
	value, ok := answer.Answers[targetID]
	return value, ok
}

func parseAnswerMap(payload map[string]any) map[string]agent.AssistantAnswer {
	raw, ok := payload["answers"]
	if !ok || raw == nil {
		return nil
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	var answers map[string]agent.AssistantAnswer
	if err := json.Unmarshal(data, &answers); err != nil {
		return nil
	}
	return answers
}
