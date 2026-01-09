package main

import (
	//"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Amarnaath548/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"context" // Required for pgx
    "github.com/jackc/pgx/v5/pgxpool"
)

type apiConfig struct{
	DB *database.Queries
}

func main() {

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found on enviroment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL==""{
		log.Fatal("DB_URL is not found in the enviroment")
	}

	ctx :=context.Background()

	dbPool, err:=pgxpool.New(ctx, dbURL)
	if err != nil {
        log.Fatal("Can't connect to database:", err)
    }
    defer dbPool.Close()

	queries :=database.New(dbPool)

	apiCfg := apiConfig{
		DB: queries,
	}
 
	// conn, err := sql.Open("postgres", dbURL)
	// if err!=nil{
	// 	log.Fatal("Can't connect to database", err)
	// }

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*","http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT" , "DELETE", "OPTIONS"},
		AllowCredentials: false,
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		MaxAge: 300,
	}))

	v1Router:=chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err",handlerErr)

	v1Router.Post("/users",apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", (apiCfg.handlerGetFeeds))

	v1Router.Post("/feed_follows",apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows",apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollow))
	v1Router.Delete("/feed_follows/{feedFollowID}",apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1",v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server is starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("port:", portString)
}
