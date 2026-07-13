package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/MoScenix/mes/app/ai/conf"
	"github.com/MoScenix/mes/app/ai/llm"
	aitools "github.com/MoScenix/mes/app/ai/tools"
	"github.com/MoScenix/mes/common/rpcmeta"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	toolutils "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

const assistantInstructionPath = "prompt/assistant/instruction.prompt"

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
	Content string                     `json:"content,omitempty"`
	Payload map[string]any             `json:"payload,omitempty"`
	Answers map[string]AssistantAnswer `json:"answers,omitempty"`
}

func init() {
	schema.RegisterName[AskUserInput]("ai_ask_user_input_v1")
	schema.RegisterName[AskUserQuestion]("ai_ask_user_question_v1")
	schema.RegisterName[AssistantInterruptState]("ai_assistant_interrupt_state_v1")
	schema.RegisterName[AssistantAnswer]("ai_assistant_answer_v1")
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
	instructionText := assistantInstruction(ctx, string(instruction))

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:             "assistant",
		Description:      "MES assistant that helps users query and operate work orders, engineering orders, inventory flows, and inventory records within role permissions.",
		Instruction:      instructionText,
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

func assistantInstruction(ctx context.Context, base string) string {
	role := rpcmeta.NormalizeRole(rpcmeta.FromContext(ctx).OperatorRole, conf.GetConf().AITools.RoleAliases)
	operatorID, _ := rpcmeta.OperatorIDFromContext(ctx)

	var b strings.Builder
	b.WriteString(strings.TrimSpace(base))
	b.WriteString("\n\nCurrent user role prompt:\n")
	if operatorID > 0 {
		b.WriteString("- Current operator user id: ")
		b.WriteString(strconv.FormatInt(operatorID, 10))
		b.WriteString(". Use this as the current user when tools omit user filters.\n")
	}
	b.WriteString("- Current normalized role: ")
	b.WriteString(role)
	b.WriteString(".\n")
	b.WriteString("- ")
	b.WriteString(roleIdentity(role))
	b.WriteString("\n")
	b.WriteString(roleWorkflowPrompt(role))
	return b.String()
}

func roleIdentity(role string) string {
	switch role {
	case rpcmeta.RoleLeader:
		return "You are the MES assistant for a team leader. The leader cares about production tasks, engineering-order planning, related work orders, and material circulation needed by the team."
	case rpcmeta.RolePurchase:
		return "You are the MES assistant for a purchase specialist. Purchase cares about material definitions, material demand, inbound circulation, and the user's own inventory-flow records."
	case rpcmeta.RoleWorker:
		return "You are the MES assistant for a worker. Workers care about assigned work orders, what needs to be done, and the related production/material context."
	case rpcmeta.RoleProcessEngineer:
		return "You are the MES assistant for a process engineer. Process engineers care about process routes, output items, material requirements, and engineering-order production context."
	case rpcmeta.RoleWarehouseAdmin:
		return "You are the MES assistant for a warehouse admin. Warehouse admins care about submitted inventory flows, item stock, concrete item units, and per-item operation progress."
	case rpcmeta.RoleSales:
		return "You are the MES assistant for a sales user. Sales cares about customer or demand follow-up, work orders they initiated, and circulation records related to delivery/material needs."
	case rpcmeta.RoleAdmin:
		return "You are the MES assistant for an admin user. Admins have a broader operational view across MES orders, users, inventory flows, items, and warehouse progress."
	default:
		return "You are the MES assistant for the current MES user. Understand the user's role first, then explain and operate MES records in that role's language."
	}
}

func roleWorkflowPrompt(role string) string {
	common := `
Shared workflow rules for this role:
- A WorkOrder is an internal task/notification order between users. It records who initiated the task, who receives it, what needs to be done, and whether it has been read or processed.
- A Process is the process plan/routing for making an item. An EngineeringOrder is the concrete production task based on a process and output item.
- An InventoryFlow is a circulation/request order for materials. It only records the movement request and per-item requested/finished quantities. It is not the stock operation itself.
- Drafts are private to their creator. When listing the current user's own records, omitted status means no status filter and drafts may appear. Global, audit, assigned-to-other, and other-user record views are non-draft business views.
- If a person is named by account/name, always call search_users first. Ask the user for an id only when search_users returns no match or multiple ambiguous matches.
- For list/search tools, request the latest 30 records by default unless the user asks for another count. Treat results as time-ordered dropdown data and continue only with the provided cursor fields.
- For inventory requests, search item definitions/material names first. Use available_count or in_stock_count for usable stock judgment; total_count is only the total number of concrete units.
`
	switch role {
	case rpcmeta.RoleLeader:
		return common + `
Leader-specific service style:
- Prefer explaining and organizing engineering orders: clarify output item, expected quantity, process, and responsible user when needed.
- The leader can query engineering orders and their related work orders; when relations are needed, use dedicated list/detail tools instead of assuming embedded relations.
- The leader's material requests are represented as inventory flows, whose progress is tracked by each requested item type.
`
	case rpcmeta.RolePurchase:
		return common + `
Purchase-specific service style:
- Prefer explaining and organizing material flow records, especially inbound/purchase-related requests.
- Purchase users mainly see their own inventory flow drafts and submitted records.
- When describing a material circulation request, keep recipient user, flow type, material item, quantity, and a clear name together.
`
	case rpcmeta.RoleWorker:
		return common + `
Worker-specific service style:
- Prefer helping the worker understand assigned work orders and related material/process context.
- Explain work orders in practical terms: who assigned it, what item/process it relates to, and what status/progress matters.
`
	case rpcmeta.RoleProcessEngineer:
		return common + `
Process-engineer-specific service style:
- Prefer helping with process, engineering-order context, and item definition lookup.
- Treat item definitions as material/product types, and ItemUnits as concrete produced or stocked single units.
- When engineering orders mention material output, verify the item definition first with search_items/get_item.
`
	case rpcmeta.RoleWarehouseAdmin:
		return common + `
Warehouse-admin-specific service style:
- Prefer helping with pending inventory flow review, inventory checks, and explaining per-item progress.
- Warehouse admins focus on non-draft inventory flows and pending submitted flows.
- Explain inventory-flow progress by item type: requested quantity, finished quantity, and the concrete ItemUnits already associated with the flow.
`
	case rpcmeta.RoleSales:
		return common + `
Sales-specific service style:
- Prefer helping sales track their initiated work orders and inventory flow drafts.
- When sales asks about material movement, explain it as an inventory flow tied to delivery or material demand.
`
	case rpcmeta.RoleAdmin:
		return common + `
Admin-specific service style:
- Admin has a wider view, while global/audit views still follow the non-draft business view rule.
- Explain cross-role records by their business meaning first, then use filters such as user, status, item, and time order when looking things up.
`
	default:
		return common
	}
}
