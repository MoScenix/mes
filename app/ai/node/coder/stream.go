package coder

import (
	"context"
	"strings"
	"time"

	"github.com/MoScenix/mes/app/ai/utils"
	"github.com/MoScenix/mes/common/aievent"
	"github.com/MoScenix/mes/common/redisstate"
	"github.com/MoScenix/mes/common/redisstream"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

func HandlePush(ctx context.Context, stateStore *redisstate.Store, store redisstream.Store, projectID string, agent string, rawContent string, loop *adk.TurnLoop[[]*schema.Message, *schema.Message], lastEventID *string, assistantOutput *utils.StringBuffer) {
	content := strings.TrimSpace(rawContent)
	if content == "" {
		return
	}
	if buffer, ok := utils.StringBufferFromContext(ctx); ok {
		buffer.Clear()
	}

	accepted, done := loop.Push([]*schema.Message{schema.UserMessage(content)})
	if done != nil {
		go func() { <-done }()
	}
	if !accepted {
		return
	}

	if err := persistPushHistory(ctx, projectID, content, assistantOutput); err != nil {
		setLastEventID(lastEventID, publishError(ctx, store, projectID, agent, err))
		return
	}

	acceptedID := publishAccepted(ctx, store, projectID, agent)
	setLastEventID(lastEventID, acceptedID)
	_ = updateProjectLastEventID(ctx, stateStore, projectID, acceptedID)
}

func persistPushHistory(ctx context.Context, projectID string, content string, assistantOutput *utils.StringBuffer) error {
	if assistantOutput != nil {
		if err := utils.AddProjectAssistantMessage(ctx, projectID, assistantOutput.String()); err != nil {
			return err
		}
	}
	if err := utils.AddProjectUserMessage(ctx, projectID, content); err != nil {
		return err
	}
	if assistantOutput != nil {
		assistantOutput.Clear()
	}
	return nil
}

func publishAccepted(ctx context.Context, store redisstream.Store, projectID string, agent string) string {
	id, _ := publishTaskEvent(ctx, store, aievent.TaskEvent{
		ProjectID: projectID,
		Type:      aievent.EventAccepted,
		Agent:     agent,
		Content:   "push accepted",
		CreatedAt: time.Now().UnixMilli(),
	})
	return id
}

func publishError(ctx context.Context, store redisstream.Store, projectID string, agent string, err error) string {
	id, _ := publishTaskEvent(ctx, store, aievent.TaskEvent{
		ProjectID: projectID,
		Type:      aievent.EventError,
		Agent:     agent,
		Content:   aievent.TrimEventContent(err.Error()),
		CreatedAt: time.Now().UnixMilli(),
	})
	return id
}

func setLastEventID(lastEventID *string, id string) {
	if lastEventID != nil && id != "" {
		*lastEventID = id
	}
}
