package docs

import (
	"objects"
	"time"
)

// return all endpoints
func Endpoints() objects.Endpoints {
	return objects.Endpoints{
		Status:      true,
		Description: "Chirpy API Documentation v1",
		Routes: map[string]objects.Route{
			"/api/endpoints": objects.Route{
				Supports: []string{
					"GET",
				},
				URL:         "",
				Params:      map[string]objects.Params{},
				AcceptsData: false,
			},
		},
		Version: "v1",
		Time:    time.Now().UTC(),
	}
}

func createRoutes() map[string]objects.Route {
	routes := map[string]objects.Route{
		"GET /endpoints": objects.Route{
			Supports:    []string{"GET"},
			URL:         "{{baseURL}}/api/endpoints",
			Params:      map[string]objects.Params{},
			AcceptsData: false,
			Format:      []string{},
		},
		"GET /healthz": objects.Route{
			Supports:    []string{"GET"},
			URL:         "{{baseURL}}/api/chirps",
			Params:      map[string]objects.Params{},
			AcceptsData: false,
			Format:      []string{},
		},
		"GET /chirps/{{id}}": objects.Route{
			Supports:    []string{"GET"},
			URL:         "{{baseURL}}/api/chirps/{{id}}",
			Params:      map[string]objects.Params{},
			AcceptsData: false,
			Format:      []string{},
		},
		"GET /chirps": objects.Route{
			Supports:    []string{"GET"},
			URL:         "{{baseURL}}/api/chirps/{{id}}",
			Params:      map[string]objects.Params{},
			AcceptsData: false,
			Format:      []string{},
		},
		"POST /chirps": objects.Route{
			Supports:    []string{"POST"},
			URL:         "{{baseURL}}/api/chirps",
			Params:      map[string]objects.Params{},
			AcceptsData: false,
			Format:      objects.RequestBodyChirp{},
		},
	}
	return routes
}
