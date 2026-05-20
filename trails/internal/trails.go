package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type TrailListOutput struct {
	Body struct {
		Resources []string `json:"resources" doc:"Names of available trails"`
	}
}

type TrailInput struct {
	Name string `path:"name" doc:"Trail name"`
}

type TrailOutput struct {
	Body struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		LengthMiles float64 `json:"length_miles"`
		Difficulty  string  `json:"difficulty" enum:"easy,moderate,hard"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
	}
}

// stub data — replace with real storage later
var knownTrails = map[string]TrailOutput{
	"blue-ridge-loop": {
		Body: struct {
			Name        string  `json:"name"`
			Description string  `json:"description"`
			LengthMiles float64 `json:"length_miles"`
			Difficulty  string  `json:"difficulty" enum:"easy,moderate,hard"`
			Latitude    float64 `json:"latitude"`
			Longitude   float64 `json:"longitude"`
		}{
			Name:        "Blue Ridge Loop",
			Description: "A scenic loop through the Blue Ridge highlands.",
			LengthMiles: 4.2,
			Difficulty:  "moderate",
			Latitude:    36.7,
			Longitude:   -81.2,
		},
	},
	"caribou": {
		Body: struct {
			Name        string  `json:"name"`
			Description string  `json:"description"`
			LengthMiles float64 `json:"length_miles"`
			Difficulty  string  `json:"difficulty" enum:"easy,moderate,hard"`
			Latitude    float64 `json:"latitude"`
			Longitude   float64 `json:"longitude"`
		}{
			Name:        "Caribou City",
			Description: "Caribou City outside of Nederland, CO",
			LengthMiles: 4.2,
			Difficulty:  "moderate",
			Latitude:    39.9795,
			Longitude:   -105.5776,
		},
	},
}

func registerTrailRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-trails",
		Method:      http.MethodGet,
		Path:        "/trail",
		Summary:     "List trails",
		Tags:        []string{"Trails"},
	}, func(ctx context.Context, _ *struct{}) (*TrailListOutput, error) {
		out := &TrailListOutput{}
		for name := range knownTrails {
			out.Body.Resources = append(out.Body.Resources, name)
		}
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-trail",
		Method:      http.MethodGet,
		Path:        "/trail/{name}",
		Summary:     "Get trail metadata",
		Tags:        []string{"Trails"},
	}, func(ctx context.Context, input *TrailInput) (*TrailOutput, error) {
		trail, ok := knownTrails[input.Name]
		if !ok {
			return nil, huma.Error404NotFound(fmt.Sprintf("trail %q not found", input.Name))
		}
		return &trail, nil
	})
}
