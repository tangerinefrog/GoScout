# GoScout

A job scraping and analyzing tool designed to simplify the process of searching for a suitable job post. It is used to scrape and collect job positions from linkedin.com and filter out irrelevant posts based on user preferences using Ollama AI.

# Setup

1. Download and install Go from official website: https://go.dev/doc/install Make sure the latest version is installed by using `go version`
2. Install goose migration tool with `go install github.com/pressly/goose/v3/cmd/goose@latest`
3. Make an .env file in the root of the repo, use example.env for reference
4. Run `goose up` to apply all database migrations
5. Use `go run ./cmd/` to start the app
6. For job position analyzing, make sure you have Ollama server running. Current version of the app uses `gpt-oss` model
