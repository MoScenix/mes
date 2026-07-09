package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/MoScenix/mes/app/ai/llm"
	aitools "github.com/MoScenix/mes/app/ai/tools"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	toolutils "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// Keep the legacy prompt path to avoid moving deployed config/assets during the assistant rename.
const assistantInstructionPath = "prompt/designer/instruction.prompt"

type AskUserInput struct {
	Questions []AskUserQuestion `json:"questions" jsonschema:"description=Questions that must be answered by the user before the assistant can continue."`
	Context   string            `json:"context,omitempty" jsonschema:"description=Short context explaining why these questions are needed."`
}

type AskUserQuestion struct {
	Question string   `json:"question" jsonschema:"description=The question text shown to the user."`
	Options  []string `json:"options,omitempty" jsonschema:"description=Suggested answer options. The user may choose one or provide custom text."`
}

type AssistantInterruptState struct {
	Questions []AskUserQuestion `json:"questions"`
	Context   string            `json:"context,omitempty"`
}

type AssistantAnswer struct {
	Content string         `json:"content,omitempty"`
	Payload map[string]any `json:"payload,omitempty"`
}

func init() {
	schema.RegisterName[AskUserInput]("ai_ask_user_input_v1")
	schema.RegisterName[AskUserQuestion]("ai_ask_user_question_v1")
	schema.RegisterName[AssistantInterruptState]("ai_designer_interrupt_state_v1")
	schema.RegisterName[AssistantAnswer]("ai_designer_answer_v1")
}

func NewAssistant(ctx context.Context) (*adk.ChatModelAgent, error) {
	cm, err := llm.NewChatModel(ctx)
	if err != nil {
		return nil, err
	}

	askTool, err := toolutils.InferTool[AskUserInput, string](
		"ask_user",
		"Ask the user one or more clarification questions and pause execution until answers are provided.",
		runAskUser,
	)
	if err != nil {
		return nil, err
	}
	searchFileTool, err := aitools.NewSearchProjectFileTool()
	if err != nil {
		return nil, err
	}
	mesTools, err := aitools.NewMESTools(ctx, askTool, searchFileTool)
	if err != nil {
		return nil, err
	}

	instruction, err := os.ReadFile(assistantInstructionPath)
	if err != nil {
		return nil, err
	}

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:             "assistant",
		Description:      "MES assistant that helps users query and operate work orders, engineering orders, inventory flows, and inventory records within role permissions.",
		Instruction:      string(instruction),
		Model:            cm,
		ModelRetryConfig: modelRetryConfig(),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: mesTools,
			},
		},
	})
}

func runAskUser(ctx context.Context, input AskUserInput) (string, error) {
	wasInterrupted, _, stored := tool.GetInterruptState[AssistantInterruptState](ctx)
	if !wasInterrupted {
		return "", tool.StatefulInterrupt(ctx, input, AssistantInterruptState{
			Questions: input.Questions,
			Context:   input.Context,
		})
	}

	isTarget, hasData, data := tool.GetResumeContext[AssistantAnswer](ctx)
	if !isTarget {
		return "", tool.StatefulInterrupt(ctx, AskUserInput{
			Questions: stored.Questions,
			Context:   stored.Context,
		}, stored)
	}
	if !hasData {
		return "", fmt.Errorf("ask_user resumed without answer data")
	}

	answer, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(answer), nil
}
