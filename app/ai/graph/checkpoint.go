package graph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/MoScenix/mes/app/ai/agent"
	"github.com/MoScenix/mes/app/ai/utils"
	"github.com/MoScenix/mes/common/aievent"
	"github.com/MoScenix/mes/common/redisstate"
	"github.com/MoScenix/mes/common/redisstream"
	"github.com/cloudwego/eino/compose"
)

var ErrInterrupted = errors.New("graph interrupted")
var ErrNoInterruptedCheckpoint = errors.New("graph has no interrupted checkpoint")
var ErrResumeAnswerNotFound = errors.New("graph resume answer not found")

type graphCheckpointStore struct {
	state *redisstate.Store
}

func newGraphCheckpointStore(ctx context.Context) (*graphCheckpointStore, bool) {
	state, ok := utils.StateStoreFromContext(ctx)
	if !ok || state == nil {
		return nil, false
	}
	return &graphCheckpointStore{state: state}, true
}

func (s *graphCheckpointStore) Get(ctx context.Context, checkpointID string) ([]byte, bool, error) {
	var data []byte
	ok, err := s.state.Get(ctx, checkpointID, &data)
	return data, ok, err
}

func (s *graphCheckpointStore) Set(ctx context.Context, checkpointID string, checkpoint []byte) error {
	return s.state.Set(ctx, checkpointID, checkpoint)
}

func Run(ctx context.Context) error {
	r, err := Buildaicode(ctx)
	if err != nil {
		return err
	}
	opts := make([]compose.Option, 0, 2)
	if _, ok := newGraphCheckpointStore(ctx); ok {
		opts = append(opts, compose.WithCheckPointID(graphCheckpointID(ctx)), compose.WithForceNewRun())
	}
	_, err = r.Invoke(ctx, map[string]any{}, opts...)
	return handleGraphResult(ctx, err)
}

func Resume(ctx context.Context) error {
	state, err := loadInterruptedGraphState(ctx)
	if err != nil {
		return err
	}
	if len(state.PendingInterrupts) == 0 || state.PendingInterrupts[0].ID == "" {
		return ErrNoInterruptedCheckpoint
	}
	answer, err := loadResumeAnswer(ctx, state, state.PendingInterrupts[0].ID)
	if err != nil {
		return err
	}
	targets, err := resumeTargets(state.PendingInterrupts, answer)
	if err != nil {
		return err
	}
	if buffer, ok := utils.StringBufferFromContext(ctx); ok {
		buffer.SetString(state.Buffer)
	}

	r, err := Buildaicode(ctx)
	if err != nil {
		return err
	}
	resumeCtx := compose.BatchResumeWithData(ctx, targets)
	_, err = r.Invoke(resumeCtx, map[string]any{}, compose.WithCheckPointID(state.CheckpointID))
	return handleGraphResult(ctx, err)
}

func loadResumeAnswer(ctx context.Context, state aievent.ProjectState, targetID string) (agent.AssistantAnswer, error) {
	streamStore, ok := utils.StreamStoreFromContext(ctx)
	if !ok || streamStore == nil {
		return agent.AssistantAnswer{}, fmt.Errorf("graph resume requires stream store")
	}
	projectID, ok := utils.ProjectIDFromContext(ctx)
	if !ok || projectID == "" {
		return agent.AssistantAnswer{}, fmt.Errorf("graph resume requires project id")
	}
	lastID := resumeEventCursor(state)
	if strings.TrimSpace(lastID) == "" {
		lastID = "0"
	}

	messages, err := streamStore.Read(ctx, aievent.ControlKey(projectID), lastID, redisstream.ReadOptions{
		Block: time.Second,
		Count: 32,
	})
	if err != nil {
		return agent.AssistantAnswer{}, err
	}
	targetIDs := resumeTargetIDs(state, targetID)
	for _, msg := range messages {
		event, err := redisstream.Decode[aievent.TaskEvent](msg)
		if err != nil || event.Type != aievent.EventAnswer {
			continue
		}
		answer := agent.AssistantAnswer{
			Content: event.Content,
			Payload: event.Payload,
			Answers: parseAnswerMap(event.Payload),
		}
		if len(targetIDs) > 0 && !answerMatchesTargets(answer, targetIDs) {
			continue
		}
		return answer, nil
	}
	return agent.AssistantAnswer{}, ErrResumeAnswerNotFound
}

func answerMatchesTargets(answer agent.AssistantAnswer, targetIDs map[string]bool) bool {
	for id := range answer.Answers {
		if targetIDs[id] {
			return true
		}
	}
	return false
}

func resumeTargets(interrupts []aievent.PendingInterrupt, answer agent.AssistantAnswer) (map[string]any, error) {
	if len(interrupts) == 0 {
		return nil, ErrNoInterruptedCheckpoint
	}
	if len(answer.Answers) == 0 {
		return nil, ErrResumeAnswerNotFound
	}
	targets := make(map[string]any, len(interrupts))
	for _, interrupt := range interrupts {
		if strings.TrimSpace(interrupt.ID) == "" {
			return nil, ErrNoInterruptedCheckpoint
		}
		value, ok := answerForPendingInterrupt(interrupt, answer.Answers)
		if !ok {
			return nil, ErrResumeAnswerNotFound
		}
		targets[interrupt.ID] = value
	}
	return targets, nil
}

func answerForPendingInterrupt(interrupt aievent.PendingInterrupt, answers map[string]agent.AssistantAnswer) (agent.AssistantAnswer, bool) {
	for id := range aievent.PendingInterruptTargetIDs(interrupt) {
		if answer, ok := answers[id]; ok {
			return answer, true
		}
	}
	return agent.AssistantAnswer{}, false
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

func resumeTargetIDs(state aievent.ProjectState, targetID string) map[string]bool {
	ids := make(map[string]bool)
	if targetID = strings.TrimSpace(targetID); targetID != "" {
		ids[targetID] = true
	}
	for _, interrupt := range state.PendingInterrupts {
		for id := range aievent.PendingInterruptTargetIDs(interrupt) {
			ids[id] = true
		}
	}
	return ids
}

func handleGraphResult(ctx context.Context, err error) error {
	if err == nil {
		_ = clearGraphCheckpoint(ctx)
		return nil
	}
	info, ok := compose.ExtractInterruptInfo(err)
	if !ok {
		return err
	}
	if persistErr := persistGraphInterrupted(ctx, info); persistErr != nil {
		return persistErr
	}
	return ErrInterrupted
}

func persistGraphInterrupted(ctx context.Context, info *compose.InterruptInfo) error {
	stateStore, ok := utils.StateStoreFromContext(ctx)
	if !ok || stateStore == nil {
		return fmt.Errorf("graph checkpoint requires state store")
	}
	projectID, ok := utils.ProjectIDFromContext(ctx)
	if !ok || projectID == "" {
		return fmt.Errorf("graph checkpoint requires project id")
	}

	interrupts := make([]aievent.PendingInterrupt, 0, len(info.InterruptContexts))
	for _, interruptCtx := range info.InterruptContexts {
		payload := map[string]any{
			"address":           interruptCtx.Address.String(),
			aievent.PayloadInfo: interruptCtx.Info,
			"is_root_cause":     interruptCtx.IsRootCause,
		}
		infoPayload, _ := interruptCtx.Info.(map[string]any)
		if assistantLastID := aievent.PayloadString(infoPayload, aievent.PayloadLastEventID); assistantLastID != "" {
			payload[aievent.PayloadLastEventID] = assistantLastID
		}
		if adkInterruptID := aievent.ADKInterruptID(infoPayload); adkInterruptID != "" {
			payload[aievent.PayloadADKInterruptID] = adkInterruptID
		}
		if controlCursor := aievent.ControlCursor(infoPayload); controlCursor != "" {
			payload[aievent.PayloadControlCursor] = controlCursor
		}
		interrupts = append(interrupts, aievent.PendingInterrupt{
			ID:      interruptCtx.ID,
			Agent:   assistantNode,
			Content: fmt.Sprint(interruptCtx.Info),
			Payload: payload,
		})
	}
	if len(interrupts) == 0 {
		return ErrNoInterruptedCheckpoint
	}

	buffer := ""
	if b, ok := utils.StringBufferFromContext(ctx); ok {
		buffer = b.String()
	}
	return stateStore.Set(ctx, aievent.RunningStateKey(projectID), aievent.ProjectState{
		Status:            aievent.ProjectStatusInterrupted,
		LastEventID:       lastEventID(ctx),
		CheckpointID:      graphCheckpointID(ctx),
		PendingInterrupts: interrupts,
		Message:           buffer,
		Buffer:            buffer,
		IsCancelled:       utils.IsCancelled(ctx),
		UpdatedAt:         time.Now().UnixMilli(),
	})
}

func resumeEventCursor(state aievent.ProjectState) string {
	if len(state.PendingInterrupts) > 0 {
		payload := state.PendingInterrupts[0].Payload
		if value := aievent.ControlCursor(payload); value != "" {
			return value
		}
		if value := aievent.PayloadString(payload, aievent.PayloadLastEventID); value != "" {
			return value
		}
		if value := aievent.PayloadString(aievent.NestedPayload(payload, aievent.PayloadInfo), aievent.PayloadLastEventID); value != "" {
			return value
		}
	}
	return state.LastEventID
}

func loadInterruptedGraphState(ctx context.Context) (aievent.ProjectState, error) {
	stateStore, ok := utils.StateStoreFromContext(ctx)
	if !ok || stateStore == nil {
		return aievent.ProjectState{}, fmt.Errorf("graph checkpoint requires state store")
	}
	projectID, ok := utils.ProjectIDFromContext(ctx)
	if !ok || projectID == "" {
		return aievent.ProjectState{}, fmt.Errorf("graph checkpoint requires project id")
	}

	var state aievent.ProjectState
	ok, err := stateStore.Get(ctx, aievent.RunningStateKey(projectID), &state)
	if err != nil {
		return aievent.ProjectState{}, err
	}
	if !ok || state.CheckpointID == "" || len(state.PendingInterrupts) == 0 {
		return aievent.ProjectState{}, ErrNoInterruptedCheckpoint
	}
	if state.Status != aievent.ProjectStatusInterrupted &&
		state.Status != aievent.ProjectStatusQueued &&
		state.Status != aievent.ProjectStatusRunning {
		return aievent.ProjectState{}, ErrNoInterruptedCheckpoint
	}
	return state, nil
}

func clearGraphCheckpoint(ctx context.Context) error {
	stateStore, ok := utils.StateStoreFromContext(ctx)
	if !ok || stateStore == nil {
		return nil
	}
	projectID, ok := utils.ProjectIDFromContext(ctx)
	if !ok || projectID == "" {
		return nil
	}
	return stateStore.Del(ctx, aievent.GraphCheckpointKey(projectID))
}

func graphCheckpointID(ctx context.Context) string {
	projectID, _ := utils.ProjectIDFromContext(ctx)
	return aievent.GraphCheckpointKey(projectID)
}

func lastEventID(ctx context.Context) string {
	stateStore, ok := utils.StateStoreFromContext(ctx)
	if !ok || stateStore == nil {
		return ""
	}
	projectID, ok := utils.ProjectIDFromContext(ctx)
	if !ok || projectID == "" {
		return ""
	}

	var state aievent.ProjectState
	ok, err := stateStore.Get(ctx, aievent.RunningStateKey(projectID), &state)
	if err != nil || !ok {
		return ""
	}
	return state.LastEventID
}
