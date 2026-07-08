package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

// TestGetHealthcheck comprueba el correcto funcionamiento del endpoint de healthcheck.
func TestGetHealthcheck(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/healthcheck", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	server := NewServer()

	if err := server.GetHealthcheck(ctx); err != nil {
		t.Fatalf("Error al ejecutar el healthcheck: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Se esperaba código 200, pero se obtuvo %d", rec.Code)
	}
}

// TestPostChat_InvalidFormat verifica que el servidor gestione correctamente
// un JSON de entrada con formato inválido devolviendo un error 400.
func TestPostChat_InvalidFormat(t *testing.T) {
	e := echo.New()

	invalidJSON := `{"content": "Test de error", "image": }`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", strings.NewReader(invalidJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	server := NewServer()

	err := server.PostChat(ctx)
	if err != nil {
		t.Fatalf("El handler falló de forma no controlada: %v", err)
	}

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Se esperaba código 400, pero se obtuvo %d", rec.Code)
	}
}

// TestPostChat_Success simula una petición de chat exitosa creando un
// servidor Ollama falso en memoria para ejecutar todo el flujo del código.
func TestPostChat_Success(t *testing.T) {
	e := echo.New()

	// Simulamos la API de Ollama respondiendo un HTTP 200 con JSON válido
	mockOllama := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "Respuesta de la IA simulada"}`))
	}))
	defer mockOllama.Close()

	// Payload válido con texto e imagen
	texto := "Hola"
	imagen := "base64placeholder"
	reqBody := ChatRequest{
		Content: &texto,
		Image:   &imagen,
	}
	jsonBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Apuntamos nuestro servidor al Ollama falso
	server := NewServer()
	server.ollamaURL = mockOllama.URL

	if err := server.PostChat(ctx); err != nil {
		t.Fatalf("Error inesperado en PostChat: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Se esperaba código 200, obtenido %d", rec.Code)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &chatResp); err != nil {
		t.Fatalf("No se pudo parsear el JSON de respuesta: %v", err)
	}

	if chatResp.Message != "Respuesta de la IA simulada" {
		t.Errorf("Mensaje inesperado: %q", chatResp.Message)
	}
}

// TestGetEnv comprueba la lectura correcta de variables de entorno del sistema.
func TestGetEnv(t *testing.T) {
	os.Setenv("VARIABLE_TEST_TFG", "control-tfg")
	defer os.Unsetenv("VARIABLE_TEST_TFG")

	val := getEnv("VARIABLE_TEST_TFG", "fallback")
	if val != "control-tfg" {
		t.Errorf("Se esperaba 'control-tfg', pero se obtuvo %q", val)
	}
}
