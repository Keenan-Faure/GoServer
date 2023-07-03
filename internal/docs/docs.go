package docs

import (
	"encoding/json"
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
				URL: "",
				Params: []string{
					"",
				},
				AcceptsData: false,
			},
		},
		Version: "v1",
		Time:    time.Now().UTC(),
	}
}

func createRoutes() map[string]objects.Route {
	routes := map[string]objects.Route{
		"/": objects.Route{
			Supports: []string{"GET"},
		},
	}
	return routes
}

func createRoute(
	supports []string,
	url string,
	params []string,
	acceptsData bool,
	format interface{}) objects.Route {
	route := objects.Route{}
	for _, value := range supports {
		route.Supports = append(route.Supports, value)
	}
	route.URL = url
	for _, value := range params {
		route.Params = append(route.Params, value)
	}
	route.AcceptsData = acceptsData
	jsonFormat, err := json.Marshal(format)
	if err != nil {
		route.Format = nil
	}
	route.Format = jsonFormat
	return route
}
