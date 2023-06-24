package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := os.Getenv(PORT)
	if port == "" {
		log.Fatal(ERR_PORT_UNDEF)
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))

	v1Router := chi.NewRouter()
	v1Router.Post("/user", userController.create)
	v1Router.Get("/user", userController.getUserByKey)
	v1Router.Get("/error", handlerError)

	router.Mount("/v1", v1Router)

	defer closeDBClient()

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server listening on port %v", port)
	if server.ListenAndServe() == nil {
		log.Fatal(ERR_START_SRV)
	}

}
