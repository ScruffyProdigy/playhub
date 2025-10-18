package graph

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestDriftDetectionDemo demonstrates how the drift detection works
// This test shows what happens when schema files are newer than generated files
func TestDriftDetectionDemo(t *testing.T) {
	// This test is for demonstration purposes only
	// In a real scenario, you would modify a schema file and forget to run gqlgen generate

	projectRoot := findProjectRoot(t)
	schemaDir := filepath.Join(projectRoot, "backend/graph/schema")
	coreSchemaPath := filepath.Join(schemaDir, "core.graphqls")

	// Get the current modification time of the schema file
	originalModTime := getModTime(t, coreSchemaPath)

	// Simulate modifying the schema file by touching it
	// In real life, this would happen when someone edits the file
	err := os.Chtimes(coreSchemaPath, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to touch schema file: %v", err)
	}

	// Now run the drift detection test
	// This should fail because the schema file is newer than generated files
	t.Run("DriftDetection", func(t *testing.T) {
		schemaFiles := findSchemaFiles(t, projectRoot)
		generatedFiles := []string{
			filepath.Join(projectRoot, "backend/graph/generated/generated.go"),
			filepath.Join(projectRoot, "backend/graph/generated.go"),
			filepath.Join(projectRoot, "backend/graph/model/models_gen.go"),
		}

		driftDetected := false
		for _, schemaFile := range schemaFiles {
			schemaModTime := getModTime(t, schemaFile)

			for _, generatedFile := range generatedFiles {
				if !fileExists(generatedFile) {
					continue
				}

				generatedModTime := getModTime(t, generatedFile)

				if schemaModTime.After(generatedModTime) {
					driftDetected = true
					t.Logf("DRIFT DETECTED: Schema file %s (modified %s) is newer than generated file %s (modified %s)",
						schemaFile, schemaModTime.Format(time.RFC3339),
						generatedFile, generatedModTime.Format(time.RFC3339))
				}
			}
		}

		if !driftDetected {
			t.Error("Expected drift to be detected, but none was found")
		}
	})

	// Restore the original modification time to clean up
	err = os.Chtimes(coreSchemaPath, originalModTime, originalModTime)
	if err != nil {
		t.Logf("Warning: Failed to restore original modification time: %v", err)
	}
}

// TestDriftDetectionInstructions provides instructions on how to use the drift detection
func TestDriftDetectionInstructions(t *testing.T) {
	// This test always passes and provides instructions
	t.Log("=== GQLGEN DRIFT DETECTION INSTRUCTIONS ===")
	t.Log("")
	t.Log("The drift detection test will catch when you forget to run 'gqlgen generate'")
	t.Log("after modifying GraphQL schema files.")
	t.Log("")
	t.Log("To test this manually:")
	t.Log("1. Modify any .graphqls file in backend/graph/schema/")
	t.Log("2. Run: go test ./graph -v -run=TestGqlgenDrift")
	t.Log("3. The test should fail with a drift detection message")
	t.Log("4. Run: go run github.com/99designs/gqlgen@v0.17.81 generate")
	t.Log("5. Run the test again - it should now pass")
	t.Log("")
	t.Log("The test checks that all generated files are newer than schema files.")
	t.Log("If any schema file is newer, it means you forgot to regenerate!")
	t.Log("")
	t.Log("Add this to your CI/CD pipeline to catch drift automatically.")
}
