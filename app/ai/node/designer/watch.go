package designer

import (
	"context"
	"strings"
	"time"

	taskrunner "github.com/MoScenix/mes/app/ai/node/coder"
	"github.com/MoScenix/mes/app/ai/node/control"
	"github.com/MoScenix/mes/app/ai/utils"
	"github.com/MoScenix/mes/common/aievent"
	"github.com/MoScenix/mes/common/redisstate"
	"github.com/MoScenix/mes/common/redisstream"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/kitex/pkg/klog"
)

func watchPushes(ctx context.Context, stateStore *redisstate.Store, store redisstream.Store, projectID string, answers chan<- answerEvent, loop *adk.TurnLoop[[]*schema.Message, *schema.Message], lastEventID *string, assistantOutput *utils.StringBuffer) {
	control.Watch(ctx, store, projectID, controlCursor(ctx, projectID), control.Handler{
		OnPush: func(ctx context.Context, msg redisstream.Message, event aievent.TaskEvent) {
			utils.SetControlCursor(ctx, msg.ID)
			taskrunner.HandlePush(ctx, stateStore, store, projectID, agentName, event.Content, loop, lastEventID, assistantOutput)
		},
		OnCancel: func(ctx context.Context, msg redisstream.Message, event aievent.TaskEvent) {
			utils.SetControlCursor(ctx, msg.ID)
			utils.CancelRuntime(ctx)
			if loop != nil {
				reason := strings.TrimSpace(event.Content)
				if reason == "" {
					reason = "cancelled"
				}
				loop.Stop(adk.WithImmediate(), adk.WithStopCause(reason), adk.WithSkipCheckpoint())
			}
		},
		OnAnswer: func(ctx context.Context, msg redisstream.Message, event aievent.TaskEvent) {
			if answers == nil {
				return
			}
			select {
			case answers <- answerEvent{
				TargetID: strings.TrimSpace(event.TargetID),
				Answer:   agentAnswer(event),
			}:
				utils.SetControlCursor(ctx, msg.ID)
			case <-ctx.Done():
				return
			}
			if err := markAnswerAccepted(ctx, projectID, strings.TrimSpace(event.TargetID), msg.ID); err != nil {
				klog.CtxErrorf(ctx, "accept assistant answer failed: project_id=%s target_id=%s err=%v", projectID, strings.TrimSpace(event.TargetID), err)
			}
		},
	})
}

func markAnswerAccepted(ctx context.Context, projectID string, targetID string, eventID string) error {
	stateStore, ok := utils.StateStoreFromContext(ctx)
	if !ok || stateStore == nil || projectID == "" || eventID == "" {
		return nil
	}

	var state aievent.ProjectState
	ok, err := stateStore.Get(ctx, aievent.RunningStateKey(projectID), &state)
	if err != nil || !ok || state.Status != aievent.ProjectStatusWaitingAnswer {
		return err
	}
	if targetID != "" && !aievent.PendingInterruptsMatch(state.PendingInterrupts, targetID) {
		return nil
	}

	state.Status = aievent.ProjectStatusRunning
	state.PendingInterrupts = nil
	state.UpdatedAt = time.Now().UnixMilli()
	return stateStore.Set(ctx, aievent.RunningStateKey(projectID), state)
}

func controlCursor(ctx context.Context, projectID string) string {
	if cursor := strings.TrimSpace(utils.ControlCursor(ctx)); cursor != "" {
		return cursor
	}
	return "$"
}
