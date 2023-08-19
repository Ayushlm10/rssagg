package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ayushlm10/rssAgg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiCfg struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Can't load environment variables")
	}
	PORT := os.Getenv("PORT")
	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		log.Fatal("Couldn't connect to the database, Issue with the db URL")
	}
	conn, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal("Couldn't connect to the database:", err)
	}
	db := database.New(conn)
	apicfg := apiCfg{
		DB: db,
	}
	go scraper(10, time.Minute, db)
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apicfg.handlerCreateUser)
	v1Router.Get("/users", apicfg.middlewareAuth(apicfg.handlerGetUser))
	v1Router.Get("/users", apicfg.middlewareAuth(apicfg.handlerGetPostForUser))
	v1Router.Post("/feeds", apicfg.middlewareAuth(apicfg.handlerCreateFeed))
	v1Router.Get("/feeds", apicfg.handlerGetAllFeeds)
	v1Router.Post("/feed-follows", apicfg.middlewareAuth(apicfg.handlerCreateFeedFollow))
	v1Router.Get("/feed-follows", apicfg.middlewareAuth(apicfg.handlerGetFeedFollows))
	v1Router.Delete("/feed-follows/{feedFollowId}", apicfg.middlewareAuth(apicfg.handlerDeleteFeedFollow))
	v1Router.Get("/posts", apicfg.middlewareAuth(apicfg.handlerGetPostForUser))
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}

	fmt.Println("Starting server at port " + PORT)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Couldn't start the server")
	}

}
