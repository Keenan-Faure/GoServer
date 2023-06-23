package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const port = "8080"
const filePath = "./app"

func main() {
	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	//creates new Chi router
	r := chi.NewRouter()
	api_router := chi.NewRouter()
	admin_router := chi.NewRouter()

	fs := http.FileServer(http.Dir(filePath))

	//Wrap the http.FileServer handler with the middleware
	fsHandle := http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fs))
	r.Handle("/app", fsHandle)
	r.Handle("/app/*", fsHandle)

	//Create a new router to bind the /healthz and /metrics to register the endpoints on,
	//and then r.Mount() that router at /api in our main router.
	api_router.Get("/healthz", healthz)
	api_router.Post("/validate_chirp", validate_chirp)
	admin_router.Get("/metrics", apiCfg.metrics)

	//re-routes the localhost:8080/metrics to be localhost:8080/api/metrics
	r.Mount("/api", api_router)
	r.Mount("/admin", admin_router)

	corsMux := middlewareCors(r)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	log.Printf("Serving files from %s on port %s", filePath, port)
	log.Fatal(server.ListenAndServe())
}

//Middleware is a way to wrap a handler with additional functionality.
//It is a common pattern in web applications that allows us to write DRY code.
//For example, we can write a middleware that logs every request to the server.
//We can then wrap our handler with this middleware and every request will be logged
//without us having to write the logging code in every handler.

//The `middlewareMetricsInc` is a middleWare because it is a handler that simply increments
//on each request

//The `middleWare` is a middleware because it adds the headers to the handler
