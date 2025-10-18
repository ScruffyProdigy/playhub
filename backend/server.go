package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/scruffyprodigy/playhub/graph"
	"github.com/scruffyprodigy/playhub/graph/generated"
)

func main() {
	mux := http.NewServeMux()

	gql := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	mux.Handle("/graphql", withAuth(gql))
	mux.Handle("/", playground.Handler("GraphQL", "/graphql"))

	mux.HandleFunc("/.well-known/jwks.json", jwksHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })

	srv := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Println("backend listening :8080")
	log.Fatal(srv.ListenAndServe())
}

func withAuth(next http.Handler) http.Handler {
	// TODO: read JWT from cookie, verify Ed25519 w/ env JWKS_PUB_X, set user in context
	return next
}

func jwksHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jwk := map[string]string{
		"kty": "OKP", "crv": "Ed25519", "alg": "EdDSA",
		"kid": os.Getenv("JWKS_KID"),
		"x":   os.Getenv("JWKS_PUB_X"),
	}
	json.NewEncoder(w).Encode(map[string]any{"keys": []any{jwk}})
}
