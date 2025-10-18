# GQLGen Drift Detection

This directory contains tests to prevent "forgot to run generate" errors when working with GraphQL schemas.

## The Problem

When you modify GraphQL schema files (`.graphqls`), you need to run `gqlgen generate` to update the generated Go code. If you forget this step, your application will fail at runtime with confusing errors.

## The Solution

We have automated tests that detect when schema files are newer than generated files, catching this issue before it reaches production.

## Test Files

- **`gqlgen_drift_test.go`** - Main drift detection tests
- **`drift_demo_test.go`** - Demonstration and instructions
- **`healthz_test.go`** - Basic functionality tests
- **`resolvers_test.go`** - Comprehensive resolver tests
- **`benchmark_test.go`** - Performance benchmarks

## How It Works

The drift detection compares the modification times of:
- **Schema files**: `backend/graph/schema/*.graphqls`
- **Generated files**: 
  - `backend/graph/generated/generated.go`
  - `backend/graph/generated.go`
  - `backend/graph/model/models_gen.go`

If any schema file is newer than any generated file, drift is detected.

## Usage

### Running Tests

```bash
# Run all tests
go test ./graph -v

# Run only drift detection
go test ./graph -v -run="TestGqlgen"

# Run drift detection with Makefile
make check-drift
```

### Using the Script

```bash
# Run the drift detection script
./scripts/check-gqlgen.sh
```

### Makefile Commands

```bash
make help              # Show all available commands
make generate          # Generate GraphQL code
make test-drift        # Run drift detection tests
make check-drift       # Check if generation is needed (for CI)
make test-all          # Run all tests including drift detection
```

## CI/CD Integration

### GitHub Actions

The `.github/workflows/gqlgen-drift.yml` workflow automatically runs drift detection on:
- Push to main/develop branches
- Pull requests to main/develop branches

### Local Development

Add this to your pre-commit hook or run it manually:

```bash
# Check for drift before committing
make check-drift
```

## Manual Testing

To test the drift detection manually:

1. **Modify a schema file**:
   ```bash
   echo "# Test comment" >> backend/graph/schema/core.graphqls
   ```

2. **Run drift detection**:
   ```bash
   make check-drift
   # Should fail with drift detected message
   ```

3. **Fix the drift**:
   ```bash
   make generate
   ```

4. **Verify it's fixed**:
   ```bash
   make check-drift
   # Should pass
   ```

## Test Coverage

The tests cover:

- ✅ **Drift Detection**: Schema files vs generated files
- ✅ **File Existence**: All required files exist
- ✅ **Generation Works**: Generated code compiles
- ✅ **Resolver Functionality**: All resolvers work correctly
- ✅ **Error Handling**: Proper error responses
- ✅ **Performance**: Benchmark tests

## Troubleshooting

### "Drift detected" but I just ran generate

This can happen if:
1. You modified a schema file after running generate
2. File system timestamps are inconsistent
3. You're running tests from a different directory

**Solution**: Run `make regenerate` to clean and regenerate everything.

### Tests fail with "file not found"

Make sure you're running tests from the `backend` directory:
```bash
cd backend
go test ./graph -v
```

### Generated code doesn't compile

Run `make regenerate` to clean and regenerate all code from scratch.

## Best Practices

1. **Always run drift detection** before committing schema changes
2. **Add drift detection to CI/CD** to catch issues automatically
3. **Use the Makefile commands** for consistency
4. **Run tests frequently** during development
5. **Keep schema files in version control** but not generated files

## Integration with IDEs

### VS Code

Add this to your `.vscode/tasks.json`:

```json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "Check GQLGen Drift",
            "type": "shell",
            "command": "make",
            "args": ["check-drift"],
            "group": "test",
            "presentation": {
                "echo": true,
                "reveal": "always",
                "focus": false,
                "panel": "shared"
            }
        }
    ]
}
```

### Pre-commit Hook

Add this to `.git/hooks/pre-commit`:

```bash
#!/bin/bash
cd backend
make check-drift
```

## Performance

The drift detection tests are very fast:
- **Drift detection**: ~0.5ms
- **All tests**: ~500ms
- **Benchmarks**: ~6s (for performance testing)

This makes them suitable for frequent execution in CI/CD pipelines.
