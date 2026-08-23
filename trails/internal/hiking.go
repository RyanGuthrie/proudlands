package internal

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"sort"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

//go:embed hiking
var hikingFS embed.FS

const hikingRoot = "hiking"

type HikingRangeListOutput struct {
	Body struct {
		Ranges []string `json:"ranges" doc:"Mountain range names"`
	}
}

type HikingRangeInput struct {
	Range string `path:"range" doc:"Mountain range name"`
}

type HikingPeakListOutput struct {
	Body struct {
		Peaks []string `json:"peaks" doc:"Peak names within the range"`
	}
}

type HikingPeakInput struct {
	Range string `path:"range" doc:"Mountain range name"`
	Peak  string `path:"peak" doc:"Peak name"`
}

type HikingRouteListOutput struct {
	Body struct {
		Routes []string `json:"routes" doc:"Route names for the peak"`
	}
}

type HikingRouteInput struct {
	Range string `path:"range" doc:"Mountain range name"`
	Peak  string `path:"peak" doc:"Peak name"`
	Route string `path:"route" doc:"Route name"`
}

type HikingRouteOutput struct {
	Body json.RawMessage
}

// listDirNames returns the names of subdirectories of dir, sorted.
func listDirNames(fsys fs.FS, dir string) ([]string, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

// listRouteNames returns the names (without .geojson extension) of route
// files in dir, sorted.
func listRouteNames(fsys fs.FS, dir string) ([]string, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".geojson") {
			continue
		}
		names = append(names, strings.TrimSuffix(e.Name(), ".geojson"))
	}
	sort.Strings(names)
	return names, nil
}

func registerHikingRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-hiking-ranges",
		Method:      http.MethodGet,
		Path:        "/hiking",
		Summary:     "List mountain ranges",
		Tags:        []string{"Hiking"},
	}, func(ctx context.Context, _ *struct{}) (*HikingRangeListOutput, error) {
		names, err := listDirNames(hikingFS, hikingRoot)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list mountain ranges", err)
		}
		out := &HikingRangeListOutput{}
		out.Body.Ranges = names
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-hiking-peaks",
		Method:      http.MethodGet,
		Path:        "/hiking/{range}",
		Summary:     "List peaks within a mountain range",
		Tags:        []string{"Hiking"},
	}, func(ctx context.Context, input *HikingRangeInput) (*HikingPeakListOutput, error) {
		dir := path.Join(hikingRoot, input.Range)
		names, err := listDirNames(hikingFS, dir)
		if err != nil {
			return nil, huma.Error404NotFound(fmt.Sprintf("mountain range %q not found", input.Range))
		}
		out := &HikingPeakListOutput{}
		out.Body.Peaks = names
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-hiking-routes",
		Method:      http.MethodGet,
		Path:        "/hiking/{range}/{peak}",
		Summary:     "List routes for a peak",
		Tags:        []string{"Hiking"},
	}, func(ctx context.Context, input *HikingPeakInput) (*HikingRouteListOutput, error) {
		dir := path.Join(hikingRoot, input.Range, input.Peak)
		names, err := listRouteNames(hikingFS, dir)
		if err != nil {
			return nil, huma.Error404NotFound(fmt.Sprintf("peak %q not found in range %q", input.Peak, input.Range))
		}
		out := &HikingRouteListOutput{}
		out.Body.Routes = names
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-hiking-route",
		Method:      http.MethodGet,
		Path:        "/hiking/{range}/{peak}/{route}",
		Summary:     "Get a route's GeoJSON",
		Tags:        []string{"Hiking"},
	}, func(ctx context.Context, input *HikingRouteInput) (*HikingRouteOutput, error) {
		filePath := path.Join(hikingRoot, input.Range, input.Peak, input.Route+".geojson")
		data, err := fs.ReadFile(hikingFS, filePath)
		if err != nil {
			return nil, huma.Error404NotFound(fmt.Sprintf("route %q not found", input.Route))
		}
		return &HikingRouteOutput{Body: json.RawMessage(data)}, nil
	})
}
