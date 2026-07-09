package coder

import (
	"context"
	"strings"
	"time"

	"github.com/MoScenix/mes/app/ai/utils"
	"github.com/MoScenix/mes/common/aievent"
	"github.com/MoScenix/mes/common/redisstate"
	"github.com/MoScenix/mes/common/redisstream"
	"github.com/cloudwego/kitex/pkg/klog"
)

const terminalTaskTTL = 10 * time.Second

type Committer interface {
	Commit(context.Context) error
}

func FinishCancelled(ctx context.Context, streamStore redisstream.Store, stateStore *redisstate.Store, projectID string, agent string, reason string) {
	if reason == "" {
		reason = "cancelled"
	}
	klog.CtxWarnf(ctx, "ai task cancelled: project_id=%s agent=%s cause=%s", projectID, agent, reason)
	publishTerminalEvent(ctx, streamStore, stateStore, projectID, agent, aievent.EventCancelled, aievent.ProjectStatusCancelled, reason)
}

func FinishError(ctx context.Context, streamStore redisstream.Store, stateStore *redisstate.Store, projectID string, agent string, err error) error {
	if err == nil {
		return nil
	}
	klog.CtxErrorf(ctx, "ai task failed: project_id=%s agent=%s err=%v", projectID, agent, err)
	publishTerminalEvent(ctx, streamStore, stateStore, projectID, agent, aievent.EventError, aievent.ProjectStatusError, err.Error())
	return err
}

func FinishDone(ctx context.Context, streamStore redisstream.Store, stateStore *redisstate.Store, projectID string, agent string, output string, lastEventID string, committer Committer) error {
	if committer != nil {
		if err := committer.Commit(ctx); err != nil {
			klog.CtxErrorf(ctx, "commit ai task changes failed: project_id=%s agent=%s err=%v", projectID, agent, err)
			publishTerminalEvent(ctx, streamStore, stateStore, projectID, agent, aievent.EventError, aievent.ProjectStatusError, err.Error())
			return err
		}
	}

	if err := utils.AddProjectAssistantMessage(ctx, projectID, output); err != nil {
		klog.CtxErrorf(ctx, "persist assistant message failed: project_id=%s agent=%s err=%v", projectID, agent, err)
		publishTerminalEvent(ctx, streamStore, stateStore, projectID, agent, aievent.EventError, aievent.ProjectStatusError, err.Error())
		return err
	}
	if strings.TrimSpace(output) != "" {
		_ = updateProjectLastEventID(ctx, stateStore, projectID, lastEventID)
	}

	publishTerminalEvent(ctx, streamStore, stateStore, projectID, agent, aievent.EventDone, aievent.ProjectStatusDone, "")
	klog.CtxInfof(ctx, "ai task completed: project_id=%s agent=%s", projectID, agent)
	return nil
}
