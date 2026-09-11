package ai

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

// ExecuteTool dispatches a plain (non-delegated) tool call by name.
// Wire each case up to the relevant service/repository call.
func ExecuteTool(name string, args map[string]interface{}) string {
	switch name {
	case "search_movies":
		// TODO: query the movies repository using args["genre"], args["language"], args["screenType"]
		return `{"movies": []}`
	case "get_movie_details":
		// TODO: look up a single movie by args["movieId"]
		return `{"movie": null}`
	case "find_shows":
		// TODO: query shows by args["movieId"], args["date"], args["screenType"]
		return `{"shows": []}`
	default:
		return fmt.Sprintf(`{"error": "unknown tool %q"}`, name)
	}
}

// ExecuteBookingSubAgent spins up a narrower-scoped Gemini session focused
// only on completing a booking/payment task, so the main orchestrator's
// context isn't polluted with transactional back-and-forth.
func ExecuteBookingSubAgent(ctx context.Context, client *genai.Client, task string, memory map[string]interface{}) string {
	// TODO: instantiate a model with a booking-only tool catalog
	// (hold_seat, confirm_payment, release_hold, apply_promo_code, etc.)
	// and run its own short tool-calling loop, returning a final summary.
	_ = ctx
	_ = client
	_ = memory
	return fmt.Sprintf(`{"status": "not_implemented", "task": %q}`, task)
}
