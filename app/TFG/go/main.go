// Package main initializes and runs the Echo web server that acts as a secure
// gateway between clients and the local Ollama AI service.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	defaultOllamaURL = "http://localhost:11434/api/generate"
	defaultModel     = "llava-diabetes"
	serverPort       = ":8080"
)

// MyServer maintains the dependencies required to handle incoming HTTP requests
// and communicate with the backend Ollama service.
type MyServer struct {
	httpClient *http.Client
	ollamaURL  string
	modelName  string
}

// OllamaRequest represents the JSON payload sent to the Ollama API endpoint.
type OllamaRequest struct {
	Model  string   `json:"model"`
	Prompt string   `json:"prompt"`
	Stream bool     `json:"stream"`
	Images []string `json:"images,omitempty"`
}

// OllamaResponse represents the JSON payload returned by the Ollama API endpoint.
type OllamaResponse struct {
	Response string `json:"response"`
}

// ChatResponse defines the standardized JSON structure for the gateway's output.
type ChatResponse struct {
	Message string `json:"message"`
}

// ChatRequest defines the expected JSON payload from the client.
type ChatRequest struct {
	Content *string `json:"content"`
	Image   *string `json:"image"`
}

// NewServer initializes and returns a MyServer instance with safe HTTP timeouts
// and configuration loaded from environment variables.
func NewServer() *MyServer {
	return &MyServer{
		httpClient: &http.Client{Timeout: 120 * time.Second},
		ollamaURL:  getEnv("OLLAMA_URL", defaultOllamaURL),
		modelName:  getEnv("OLLAMA_MODEL", defaultModel),
	}
}

// GetHealthcheck handles the GET /api/v1/healthcheck endpoint to verify server availability.
func (s *MyServer) GetHealthcheck(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

// PostChat handles the POST /api/v1/chat endpoint. It processes the user query
// and an optional base64-encoded image, delegating the inference to the local AI model.
func (s *MyServer) PostChat(ctx echo.Context) error {
	var reqBody ChatRequest
	if err := ctx.Bind(&reqBody); err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid request format")
	}

	content := ""
	if reqBody.Content != nil {
		content = *reqBody.Content
	}

	ollamaReq := OllamaRequest{
		Model:  s.modelName,
		Prompt: content,
		Stream: false,
	}

	if reqBody.Image != nil && *reqBody.Image != "" {
		ollamaReq.Images = []string{*reqBody.Image}
	}

	ollamaBody, err := json.Marshal(ollamaReq)
	if err != nil {
		ctx.Logger().Errorf("failed to marshal ollama request: %v", err)
		return ctx.String(http.StatusInternalServerError, "Internal Server Error")
	}

	resp, err := s.httpClient.Post(s.ollamaURL, "application/json", bytes.NewBuffer(ollamaBody))
	if err != nil {
		ctx.Logger().Errorf("ollama connection error: %v", err)
		return ctx.String(http.StatusBadGateway, "Service Unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ctx.Logger().Errorf("ollama returned status: %d", resp.StatusCode)
		return ctx.String(http.StatusBadGateway, "Error from AI Service")
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		ctx.Logger().Errorf("failed to decode ollama response: %v", err)
		return ctx.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return ctx.JSON(http.StatusOK, ChatResponse{
		Message: ollamaResp.Response,
	})
}

func main() {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	server := NewServer()
	RegisterHandlers(e, server)

	fmt.Printf("Server running on port %s...\n", serverPort)
	e.Logger.Fatal(e.Start(serverPort))
}

// getEnv retrieves the value of the environment variable named by the key.
// It returns the fallback value if the variable is not present.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
