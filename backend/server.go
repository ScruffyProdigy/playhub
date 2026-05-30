package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/scruffyprodigy/playhub/database"
	"github.com/scruffyprodigy/playhub/graph"
	"github.com/scruffyprodigy/playhub/graph/generated"
	"github.com/scruffyprodigy/playhub/internal/auth"
	"github.com/scruffyprodigy/playhub/internal/store"
)

func main() {
	if err := database.InitWithMigrations(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.Close()

	dataStore := store.New(database.GetDB())
	if err := dataStore.Ping(context.Background()); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	signer, err := auth.LoadSignerFromEnv()
	if err != nil {
		log.Fatalf("JWT signer initialization failed: %v", err)
	}

	authService, err := auth.NewService(dataStore, signer)
	if err != nil {
		log.Fatalf("Auth service initialization failed: %v", err)
	}

	resolver := graph.NewResolver(dataStore, authService)

	mux := http.NewServeMux()

	gql := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	mux.Handle("/graphql", auth.Middleware(signer, gql))
	mux.Handle("/", playground.Handler("GraphQL", "/graphql"))

	mux.Handle("/.well-known/jwks.json", signer.JWKSHandler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })

	handler := auth.CORSMiddleware(mux)

	srv := &http.Server{Addr: ":8080", Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	log.Println("backend listening :8080")
	log.Fatal(srv.ListenAndServe())
}
