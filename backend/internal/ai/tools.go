package ai

import "github.com/google/generative-ai-go/genai"

var SearchMoviesTool = &genai.FunctionDeclaration{
	Name:        "search_movies",
	Description: "Search movies by genre, language, release date, screen type, or format.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"genre":      {Type: genai.TypeString},
			"language":   {Type: genai.TypeString},
			"screenType": {Type: genai.TypeString, Description: "E.g., IMAX, STANDARD"},
		},
	},
}

var GetMovieDetailsTool = &genai.FunctionDeclaration{
	Name:        "get_movie_details",
	Description: "Fetch full details for a single movie by its ID, including runtime, synopsis, and age rating.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"movieId": {Type: genai.TypeString},
		},
		Required: []string{"movieId"},
	},
}

var FindShowsTool = &genai.FunctionDeclaration{
	Name:        "find_shows",
	Description: "Find upcoming showtimes for a movie, optionally filtered by theatre, screen type, or date.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"movieId":    {Type: genai.TypeString},
			"date":       {Type: genai.TypeString, Description: "ISO 8601 date"},
			"screenType": {Type: genai.TypeString},
		},
		Required: []string{"movieId"},
	},
}

var DelegateToBookingTool = &genai.FunctionDeclaration{
	Name:        "delegate_to_booking",
	Description: "Hand off complex, multi-step booking or payment intents to the specialized booking sub-agent.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"task": {Type: genai.TypeString, Description: "The booking task to accomplish."},
		},
		Required: []string{"task"},
	},
}

// MainCatalog combines all tools exposed to the top-level orchestrator model.
// Add further declarations here as they're implemented (get_reviews,
// find_theatres, cancel_booking, apply_promo_code, etc.) to reach the
// full 20+ tool surface described in the architecture doc.
var MainCatalog = []*genai.Tool{
	{FunctionDeclarations: []*genai.FunctionDeclaration{
		SearchMoviesTool,
		GetMovieDetailsTool,
		FindShowsTool,
		DelegateToBookingTool,
	}},
}
