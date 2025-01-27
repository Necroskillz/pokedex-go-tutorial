set shell := ["powershell.exe", "-c"]

# List available recipes
default:
    @just --list

# Build the project
build:
    go build -o bin/pokedex.exe ./cmd/pokedex

# Run the project
run: 
    go run ./cmd/pokedex

# Run tests
test:
    go test ./...

# Clean build artifacts
clean:
    rm -rf bin/ 