package docs

import (
	"objects"
	"time"
)

// return all endpoints
func Endpoints() objects.Endpoints {
	return objects.Endpoints{
		Status:      true,
		Description: "Chirpy API Documentation",
		Routes:      createRoutes(),
		Version:     "v1",
		Time:        time.Now().UTC(),
	}
}

func createRoutes() map[string]objects.Route {
	routes := map[string]objects.Route{
		"GET /api": {
			Description:   "Displays available endpoints",
			Supports:      []string{"GET"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "None",
		},
		"GET /api/healthz": {
			Description:   "Returns the status of the API",
			Supports:      []string{"GET"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "None",
		},
		"GET /api/chirps/{{id}}": {
			Description:   "Returns a specific chirp",
			Supports:      []string{"GET"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "None",
		},
		"GET /api/chirps": {
			Description: "Returns an array of chirps",
			Supports:    []string{"GET"},
			Params: map[string]objects.Params{
				"author_id": {
					Key:   "author_id",
					Value: "",
				},
				"sort": {
					Key:   "sort",
					Value: "",
				},
			},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "None",
		},
		"POST /api/chirps": {
			Description:   "Adds a new chirp",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestBodyChirp{},
			Authorization: "Bearer <token> in header",
		},
		"DELETE /api/chirps{{id}}": {
			Description:   "Deletes a specific chirp",
			Supports:      []string{"DELETE"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        objects.RequestBodyChirp{},
			Authorization: "Bearer <token> in header",
		},
		"POST /api/login": {
			Description:   "Logs into a user account",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestBodyLogin{},
			Authorization: "None",
		},
		"POST /api/users": {
			Description:   "Creates a new user account",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestBodyUser{},
			Authorization: "None",
		},
		"PUT /api/users": {
			Description:   "Updates an existing user account",
			Supports:      []string{"PUT"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.RequestBodyUser{},
			Authorization: "Bearer <token> in header",
		},
		"POST /api/refresh": {
			Description:   "Refreshes Access Token",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Bearer <token> in header",
		},
		"POST /api/revoke": {
			Description:   "Revokes the token present in the header",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   false,
			Format:        []string{},
			Authorization: "Bearer <token> in header",
		},
		"POST /polka/webhooks": {
			Description:   "Accepts payment details from Polka",
			Supports:      []string{"POST"},
			Params:        map[string]objects.Params{},
			AcceptsData:   true,
			Format:        objects.WebhookRequest{},
			Authorization: "ApiKey <key> in header",
		},
	}
	return routes
}
