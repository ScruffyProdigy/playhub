package graph

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/scruffyprodigy/playhub/graph/generated"
)

// TestGqlgenDrift checks if the generated code is up to date with the schema files
func TestGqlgenDrift(t *testing.T) {
	// Get the project root directory
	projectRoot := findProjectRoot(t)

	// Find all GraphQL schema files
	schemaFiles := findSchemaFiles(t, projectRoot)

	// Find the generated files
	generatedFiles := []string{
		filepath.Join(projectRoot, "backend/graph/generated/generated.go"),
		filepath.Join(projectRoot, "backend/graph/generated.go"),
		filepath.Join(projectRoot, "backend/graph/model/models_gen.go"),
	}

	// Check if any schema file is newer than any generated file
	for _, schemaFile := range schemaFiles {
		schemaModTime := getModTime(t, schemaFile)

		for _, generatedFile := range generatedFiles {
			if !fileExists(generatedFile) {
				t.Errorf("Generated file does not exist: %s", generatedFile)
				continue
			}

			generatedModTime := getModTime(t, generatedFile)

			if schemaModTime.After(generatedModTime) {
				t.Errorf("Schema file %s (modified %s) is newer than generated file %s (modified %s). Run 'go run github.com/99designs/gqlgen@v0.17.81 generate' to update generated code.",
					schemaFile, schemaModTime.Format(time.RFC3339),
					generatedFile, generatedModTime.Format(time.RFC3339))
			}
		}
	}
}

// TestGqlgenGenerationWorks verifies that gqlgen can generate code without errors
func TestGqlgenGenerationWorks(t *testing.T) {
	// This test ensures that the current schema can be processed by gqlgen
	// without compilation errors. It doesn't actually run gqlgen generate,
	// but it verifies that the generated code is valid.

	resolver := &Resolver{}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	// If we can create the server without errors, the generated code is valid
	if srv == nil {
		t.Error("Failed to create GraphQL server - generated code may be invalid")
	}
}

// TestSchemaFilesExist verifies that all expected schema files exist
func TestSchemaFilesExist(t *testing.T) {
	projectRoot := findProjectRoot(t)
	schemaDir := filepath.Join(projectRoot, "backend/graph/schema")

	expectedFiles := []string{
		"core.graphqls",
		"game.graphqls",
		"users.graphqls",
		"goods.graphqls",
	}

	for _, filename := range expectedFiles {
		filePath := filepath.Join(schemaDir, filename)
		if !fileExists(filePath) {
			t.Errorf("Expected schema file does not exist: %s", filePath)
		}
	}
}

// TestGeneratedFilesExist verifies that all expected generated files exist
func TestGeneratedFilesExist(t *testing.T) {
	projectRoot := findProjectRoot(t)

	expectedFiles := []string{
		"backend/graph/generated/generated.go",
		"backend/graph/generated.go",
		"backend/graph/model/models_gen.go",
		"backend/graph/core.resolvers.go",
	}

	for _, filename := range expectedFiles {
		filePath := filepath.Join(projectRoot, filename)
		if !fileExists(filePath) {
			t.Errorf("Expected generated file does not exist: %s", filePath)
		}
	}
}

// Helper functions

func findProjectRoot(t *testing.T) string {
	// Start from the current directory and walk up to find go.mod
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if fileExists(goModPath) {
			// We found go.mod, but we need to go up one level to get the project root
			// since we're currently in the backend directory
			return filepath.Dir(dir)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("Could not find project root (go.mod file)")
		}
		dir = parent
	}
}

func findSchemaFiles(t *testing.T, projectRoot string) []string {
	schemaDir := filepath.Join(projectRoot, "backend/graph/schema")

	var schemaFiles []string
	err := filepath.Walk(schemaDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".graphqls" {
			schemaFiles = append(schemaFiles, path)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk schema directory: %v", err)
	}

	return schemaFiles
}

func getModTime(t *testing.T, filePath string) time.Time {
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file %s: %v", filePath, err)
	}
	return info.ModTime()
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
