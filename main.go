package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/valentin-popov/rss-aggregator/db"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := os.Getenv(PORT)
	if port == "" {
		log.Fatal(ERR_MSG_PORT_UNDEF)
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))

	v1Router := chi.NewRouter()
	v1Router.Post("/user", createUser)
	v1Router.Get("/user", getAuthUser)

	v1Router.Post("/feed", createFeed)
	v1Router.Get("/feed", findFeeds)
	v1Router.Get("/feed/{feedId}", findFeed)
	v1Router.Post("/follow", followFeed)
	v1Router.Delete("/follow/{feedId}", unfollowFeed)
	v1Router.Get("/follow", findFollowedFeeds)

	router.Mount("/v1", v1Router)

	defer db.CloseDBClient()

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server listening on port %v", port)
	if server.ListenAndServe() == nil {
		log.Fatal(ERR_MSG_START_SRV)
	}

}
