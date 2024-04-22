package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type InvalidRquest struct {
	Link string `json:"link"`
	Op   string `json:"op"`
}

func SetupMockServer() *http.Server {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", ProcessUrl)

	serverMuxMiddlware := ResponseLogging(serverMux)

	//Hardcoded port number, only for testing
	log.Println("Testing mock server on port", 50503)

	server := &http.Server{Addr: ":50501", Handler: serverMuxMiddlware}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	return server
}

func TestMain(m *testing.M) {
	server := SetupMockServer()
	defer server.Close()
	exitCode := m.Run()
	os.Exit(exitCode)
}

// Need to write further tests
func TestServer(t *testing.T) {
	t.Run("Test canonical with example URL", func(t *testing.T) {
		t.Log("Test canonical with example URL")

		reqBody := RequestStruct{
			Url:       "https://BYFOOD.com/food-EXPeriences?query=abc/",
			Operation: "canonical",
		}

		jsonBody, _ := json.Marshal(reqBody)

		rr := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		rr.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		ProcessUrl(w, rr)

		resultBody := ResponseStruct{}
		json.Unmarshal(w.Body.Bytes(), &resultBody)

		if resultBody.ProcessedUrl != "https://BYFOOD.com/food-EXPeriences" {
			t.Errorf("Expected %s, got %s", "https://BYFOOD.com/food-EXPeriences", resultBody.ProcessedUrl)
		}

	})

	//This is really a test for the IsUrl helper function, it filters non valid Urls and it's kind of a middlware for the API
	t.Run("Test canonical with broken URL", func(t *testing.T) {
		t.Log("Test canonical with broken URL")

		reqBody := RequestStruct{
			Url:       "//*BYFOOD.com/food-EXPeriences?query=abc/",
			Operation: "canonical",
		}

		jsonBody, _ := json.Marshal(reqBody)

		rr := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		rr.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		ProcessUrl(w, rr)
		resultBody := ErrMessage{}
		json.Unmarshal(w.Body.Bytes(), &resultBody)
		if resultBody.Msg != "Url format invalid" {
			t.Errorf("Expected %s, got %s", "Url format invalid", resultBody.Msg)
		}
	})

	t.Run("Test Redirection with working URL", func(t *testing.T) {
		t.Log("Test Redirection with working URL")

		reqBody := RequestStruct{
			Url:       "https://ByFooD.com/FOOD-EXPeriences/",
			Operation: "redirection",
		}

		jsonBody, _ := json.Marshal(reqBody)

		rr := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		rr.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		ProcessUrl(w, rr)
		resultBody := ResponseStruct{}
		json.Unmarshal(w.Body.Bytes(), &resultBody)
		if resultBody.ProcessedUrl != "https://www.byfood.com/food-experiences/" {
			t.Errorf("Expected %s, got %s", "https://www.byfood.com/food-experiences/", resultBody.ProcessedUrl)
		}
	})

	//Test with working url, but domain name does not match byfood's
	t.Run("Test Redirection with non ByFood Domain", func(t *testing.T) {
		t.Log("Test Redirection with non ByFood Domain")

		reqBody := RequestStruct{
			Url:       "https://BootlegFood.com/FOOD-EXPeriences/",
			Operation: "redirection",
		}

		jsonBody, _ := json.Marshal(reqBody)

		rr := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		rr.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		ProcessUrl(w, rr)
		resultBody := ErrMessage{}
		json.Unmarshal(w.Body.Bytes(), &resultBody)
		if resultBody.Msg != "URL is not from ByFood Domain" {
			t.Errorf("Expected %s, got %s", "URL is not from ByFood Domain", resultBody.Msg)
		}
	})

	//Test with a non existent operation
	t.Run("Test with invalid operation", func(t *testing.T) {
		t.Log("Test with invalid op")

		reqBody := RequestStruct{
			Url:       "https://BootlegFood.com/FOOD-EXPeriences/",
			Operation: "redir",
		}

		jsonBody, _ := json.Marshal(reqBody)

		rr := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		rr.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		ProcessUrl(w, rr)
		resultBody := ErrMessage{}
		json.Unmarshal(w.Body.Bytes(), &resultBody)
		if resultBody.Msg != "Invalid operation" {
			t.Errorf("Expected %s, got %s", "Invalid operation", resultBody.Msg)
		}
	})

	//Test with invalid body
	t.Run("Test with invalid body", func(t *testing.T) {
		t.Log("Test with invalid body")

		reqBody := InvalidRquest{
			Link: "https://BootlegFood.com/FOOD-EXPeriences/",
			Op:   "redir",
		}

		jsonBody, _ := json.Marshal(reqBody)

		rr := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		rr.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		ProcessUrl(w, rr)
		resultBody := ErrMessage{}
		json.Unmarshal(w.Body.Bytes(), &resultBody)
		if resultBody.Msg != "Invalid request format" {
			t.Errorf("Expected %s, got %s", "Invalid request format", resultBody.Msg)
		}
	})
}
