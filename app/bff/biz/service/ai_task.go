package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MoScenix/mes/app/bff/biz/dal/redis"
	"github.com/MoScenix/mes/app/bff/biz/utils"
	lapp "github.com/MoScenix/mes/app/bff/hertz_gen/bff/app"
	"github.com/MoScenix/mes/app/bff/infra/rpc"
	"github.com/MoScenix/mes/common/aievent"
	"github.com/MoScenix/mes/common/redisstate"
	"github.com/MoScenix/mes/common/redisstream"
	rpcai "github.com/MoScenix/mes/rpc_gen/kitex_gen/ai"
	rpcapp "github.com/MoScenix/mes/rpc_gen/kitex_gen/app"
)

type aiAnswerPayload struct {
	Content string         `json:"content,omitempty"`
	Payload map[string]any `json:"payload,omitempty"`
}

func submitAITask(ctx context.Context, appID int64, message string) (bool, error) {
	if appID <= 0 {
		return false, fmt.Errorf("appId is required")
	}
	if err := requireAppOwnerOrAdmin(ctx, appID); err != nil {
		return false, err
	}
	ctx = utils.WithIdentityMeta(ctx)
	if strings.TrimSpace(message) != "" {
		if err := addUserMessage(ctx, appID, message); err != nil {
			return false, err
		}
	}
	return submitAI(ctx, appID)
}

func pushAIEvent(ctx context.Context, appID int64, content string) (string, error) {
	if appID <= 0 {
		return "", fmt.Errorf("appId is required")
	}
	if err := requireAppOwnerOrAdmin(ctx, appID); err != nil {
		return "", err
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return "", fmt.Errorf("content is required")
	}
	return addTaskEvent(ctx, appID, aievent.TaskEvent{
		ProjectID: projectID(appID),
		Type:      aievent.EventPush,
		Content:   content,
	})
}

func answerAIQuestion(ctx context.Context, appID int64, answers map[string]aiAnswerPayload) (bool, error) {
	if appID <= 0 {
		return false, fmt.Errorf("appId is required")
	}
	if err := requireAppOwnerOrAdmin(ctx, appID); err != nil {
		return false, err
	}
	ctx = utils.WithIdentityMeta(ctx)
	answers = normalizeAnswerPayloads(answers)
	if len(answers) == 0 {
		return false, fmt.Errorf("answers are required")
	}

	state, ok, err := loadAIState(ctx, appID)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, fmt.Errorf("ai task state not found")
	}
	expandedAnswers, err := expandAnswersForPendingInterrupts(state.PendingInterrupts, answers)
	if err != nil {
		return false, err
	}
	targetID := firstAnswerTarget(expandedAnswers)

	_, err = addTaskEvent(ctx, appID, aievent.TaskEvent{
		ProjectID: projectID(appID),
		Type:      aievent.EventAnswer,
		TargetID:  targetID,
		Payload: map[string]any{
			"answers": expandedAnswers,
		},
	})
	if err != nil {
		return false, err
	}

	if state.Status == aievent.ProjectStatusInterrupted {
		return submitAI(ctx, appID)
	}
	if state.Status == aievent.ProjectStatusWaitingAnswer {
		return resumeIfAnswerTimedOut(ctx, appID, targetID)
	}
	return true, nil
}

func resumeIfAnswerTimedOut(ctx context.Context, appID int64, targetID string) (bool, error) {
	deadline := time.Now().Add(5 * time.Second)
	for {
		time.Sleep(150 * time.Millisecond)
		latest, ok, err := loadAIState(ctx, appID)
		if err != nil {
			return false, err
		}
		if !ok {
			return true, nil
		}
		switch latest.Status {
		case aievent.ProjectStatusInterrupted:
			if hasPendingInterrupt(latest, targetID) {
				return submitAI(ctx, appID)
			}
			return true, nil
		case aievent.ProjectStatusWaitingAnswer:
			if time.Now().Before(deadline) {
				continue
			}
			return true, nil
		default:
			return true, nil
		}
	}
}

func cancelAIEvent(ctx context.Context, appID int64, reason string) (string, error) {
	if appID <= 0 {
		return "", fmt.Errorf("appId is required")
	}
	if err := requireAppOwnerOrAdmin(ctx, appID); err != nil {
		return "", err
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		reason = "cancelled"
	}
	return addTaskEvent(ctx, appID, aievent.TaskEvent{
		ProjectID: projectID(appID),
		Type:      aievent.EventCancel,
		Content:   reason,
	})
}

func loadAIState(ctx context.Context, appID int64) (aievent.ProjectState, bool, error) {
	if err := requireAppOwnerOrAdmin(ctx, appID); err != nil {
		return aievent.ProjectState{}, false, err
	}
	stateStore, err := stateStore()
	if err != nil {
		return aievent.ProjectState{}, false, err
	}
	var state aievent.ProjectState
	ok, err := stateStore.Get(ctx, aievent.RunningStateKey(projectID(appID)), &state)
	return state, ok, err
}

func listAIEvents(ctx context.Context, appID int64, lastID string, blockMS int64, count int64) (*lapp.AIEvents, error) {
	if appID <= 0 {
		return nil, fmt.Errorf("appId is required")
	}
	if err := requireAppOwnerOrAdmin(ctx, appID); err != nil {
		return nil, err
	}
	lastID = strings.TrimSpace(lastID)
	if lastID == "" {
		lastID = "0"
	}
	if count <= 0 {
		count = 50
	}
	block := time.Duration(blockMS) * time.Millisecond
	if block < 0 {
		block = 0
	}

	store, err := streamStore()
	if err != nil {
		return nil, err
	}
	messages, err := store.Read(ctx, aievent.EventKey(projectID(appID)), lastID, redisstream.ReadOptions{
		Block: block,
		Count: int64(count),
	})
	if err != nil {
		return nil, err
	}

	resp := &lapp.AIEvents{Events: make([]*lapp.AIEvent, 0, len(messages)), LastId: lastID}
	for _, msg := range messages {
		event, err := redisstream.Decode[aievent.TaskEvent](msg)
		if err != nil {
			continue
		}
		resp.Events = append(resp.Events, toAIEvent(msg.ID, event))
		resp.LastId = msg.ID
	}
	return resp, nil
}

func addUserMessage(ctx context.Context, appID int64, content string) error {
	userID, _ := utils.UserIDFromContext(ctx)
	_, err := rpc.AppClient.AddMessage(ctx, &rpcapp.AddMessageReq{
		AppId:   appID,
		UserId:  userID,
		Content: content,
		Role:    "user",
	})
	return err
}

func submitAI(ctx context.Context, appID int64) (bool, error) {
	ctx = utils.WithIdentityMeta(ctx)
	resp, err := rpc.AiClient.Chat(ctx, &rpcai.AiReq{ProjectId: projectID(appID)})
	if err != nil {
		return false, err
	}
	return resp.GetAnswer() == "true", nil
}

func addTaskEvent(ctx context.Context, appID int64, event aievent.TaskEvent) (string, error) {
	store, err := streamStore()
	if err != nil {
		return "", err
	}
	if event.ProjectID == "" {
		event.ProjectID = projectID(appID)
	}
	if event.CreatedAt == 0 {
		event.CreatedAt = time.Now().UnixMilli()
	}
	return store.Add(ctx, aievent.ControlKey(projectID(appID)), event)
}

func streamStore() (redisstream.Store, error) {
	return redisstream.NewRedisStore(redis.RedisClient, "ai")
}

func stateStore() (*redisstate.Store, error) {
	return redisstate.NewStore(redis.RedisClient, "ai")
}

func projectID(appID int64) string {
	return strconv.FormatInt(appID, 10)
}

func hasPendingInterrupt(state aievent.ProjectState, targetID string) bool {
	targetID = strings.TrimSpace(targetID)
	if targetID == "" {
		return false
	}
	return aievent.PendingInterruptsMatch(state.PendingInterrupts, targetID)
}

func normalizeAnswerPayloads(answers map[string]aiAnswerPayload) map[string]aiAnswerPayload {
	normalized := make(map[string]aiAnswerPayload, len(answers))
	for id, answer := range answers {
		id = strings.TrimSpace(id)
		answer.Content = strings.TrimSpace(answer.Content)
		if id == "" || (answer.Content == "" && len(answer.Payload) == 0) {
			continue
		}
		normalized[id] = answer
	}
	return normalized
}

func expandAnswersForPendingInterrupts(interrupts []aievent.PendingInterrupt, answers map[string]aiAnswerPayload) (map[string]aiAnswerPayload, error) {
	if len(interrupts) == 0 {
		return nil, fmt.Errorf("pending interrupt target not found")
	}
	expanded := make(map[string]aiAnswerPayload, len(answers))
	for answerID, answer := range answers {
		matched := false
		for _, interrupt := range interrupts {
			targetIDs := aievent.PendingInterruptTargetIDs(interrupt)
			if !targetIDs[answerID] {
				continue
			}
			for id := range targetIDs {
				expanded[id] = answer
			}
			matched = true
			break
		}
		if !matched {
			return nil, fmt.Errorf("answer target %s does not match pending interrupts", answerID)
		}
	}
	if len(expanded) == 0 {
		return nil, fmt.Errorf("answers are required")
	}
	return expanded, nil
}

func firstAnswerTarget(answers map[string]aiAnswerPayload) string {
	for id := range answers {
		return id
	}
	return ""
}

func toAIState(exists bool, state aievent.ProjectState) *lapp.AIState {
	resp := &lapp.AIState{
		Exists:       exists,
		Status:       state.Status,
		Agent:        state.Agent,
		LastEventId:  state.LastEventID,
		CheckpointId: state.CheckpointID,
		Message:      state.Message,
		Buffer:       state.Buffer,
		IsCancelled:  state.IsCancelled,
		UpdatedAt:    state.UpdatedAt,
	}
	resp.PendingInterrupts = make([]*lapp.AIPendingInterrupt, 0, len(state.PendingInterrupts))
	for _, interrupt := range state.PendingInterrupts {
		resp.PendingInterrupts = append(resp.PendingInterrupts, &lapp.AIPendingInterrupt{
			Id:          interrupt.ID,
			Agent:       interrupt.Agent,
			Content:     interrupt.Content,
			PayloadJson: marshalPayload(interrupt.Payload),
		})
	}
	return resp
}

func toAIEvent(id string, event aievent.TaskEvent) *lapp.AIEvent {
	return &lapp.AIEvent{
		Id:          id,
		ProjectId:   event.ProjectID,
		Type:        string(event.Type),
		Agent:       event.Agent,
		Content:     event.Content,
		TargetId:    event.TargetID,
		Name:        event.Name,
		Status:      event.Status,
		PayloadJson: marshalPayload(event.Payload),
		CreatedAt:   event.CreatedAt,
		Questions:   toAIQuestions(event.Payload),
	}
}

func toAIQuestions(payload map[string]any) []*lapp.AIQuestion {
	raw, ok := payload["questions"]
	if !ok {
		return nil
	}
	data, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	var questions []struct {
		Question string   `json:"question"`
		Options  []string `json:"options"`
	}
	if err := json.Unmarshal(data, &questions); err != nil {
		return nil
	}
	out := make([]*lapp.AIQuestion, 0, len(questions))
	for _, question := range questions {
		text := strings.TrimSpace(question.Question)
		if text == "" {
			continue
		}
		options := make([]string, 0, len(question.Options))
		for _, option := range question.Options {
			option = strings.TrimSpace(option)
			if option != "" {
				options = append(options, option)
			}
		}
		out = append(out, &lapp.AIQuestion{
			Question: text,
			Options:  options,
		})
	}
	return out
}

func marshalPayload(payload map[string]any) string {
	if len(payload) == 0 {
		return ""
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(data)
}
