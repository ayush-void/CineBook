package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

type Orchestrator struct {
	Client *genai.Client
}

// ProcessMessage runs a single user turn through Gemini, looping on
// function calls until the model returns plain text. Booking-related
// intents are delegated to a sub-agent to keep the main context lean.
func (o *Orchestrator) ProcessMessage(ctx context.Context, userMsg string, memory map[string]interface{}) (string, error) {
	model := o.Client.GenerativeModel("gemini-1.5-pro")
	model.Tools = MainCatalog

	// Inject contextual memory so intent isn't lost across 20+ turn conversations
	memoryStr, _ := json.Marshal(memory)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(fmt.Sprintf(
				"You are CineBook, a helpful movie ticketing assistant. User context: %s. Use delegate_to_booking for transactions.",
				memoryStr,
			)),
		},
	}

	session := model.StartChat()

	resp, err := session.SendMessage(ctx, genai.Text(userMsg))
	if err != nil {
		return "", err
	}

	for {
		var toolCalls []genai.FunctionCall
		for _, part := range resp.Candidates[0].Content.Parts {
			if call, ok := part.(genai.FunctionCall); ok {
				toolCalls = append(toolCalls, call)
			}
		}

		if len(toolCalls) == 0 {
			break // model returned plain text, stop looping
		}

		var toolResponses []genai.Part
		for _, call := range toolCalls {
			var result string

			if call.Name == "delegate_to_booking" {
				task, _ := call.Args["task"].(string)
				result = ExecuteBookingSubAgent(ctx, o.Client, task, memory)
			} else {
				result = ExecuteTool(call.Name, call.Args)
			}

			toolResponses = append(toolResponses, genai.FunctionResponse{
				Name:     call.Name,
				Response: map[string]any{"result": result},
			})
		}

		resp, err = session.SendMessage(ctx, toolResponses...)
		if err != nil {
			return "", err
		}
	}

	var finalText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			finalText += string(txt)
		}
	}

	return finalText, nil
}
