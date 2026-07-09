package designer

import (
	"context"
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
			if event.TargetID != "" && event.TargetID != targetID {
				continue
			}
			return event.Answer, true, nil
		}
	}
}

func agentAnswer(event aievent.TaskEvent) agent.AssistantAnswer {
	return agent.AssistantAnswer{
		Content: event.Content,
		Payload: event.Payload,
	}
}
