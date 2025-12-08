package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
)

// MockServer wraps httptest.Server with helper methods
type MockServer struct {
	Server   *httptest.Server
	Mux      *http.ServeMux
	handlers map[string]http.HandlerFunc
}

// NewMockServer creates a new mock server for testing
func NewMockServer() *MockServer {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	ms := &MockServer{
		Server:   server,
		Mux:      mux,
		handlers: make(map[string]http.HandlerFunc),
	}

	// Set up default routes
	ms.setupDefaultRoutes()

	return ms
}

// URL returns the mock server URL
func (ms *MockServer) URL() string {
	return ms.Server.URL
}

// Close shuts down the mock server
func (ms *MockServer) Close() {
	ms.Server.Close()
}

// SetHandler sets a custom handler for a specific path
func (ms *MockServer) SetHandler(method, path string, handler http.HandlerFunc) {
	key := method + " " + path
	ms.handlers[key] = handler
	ms.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	})
}

func (ms *MockServer) setupDefaultRoutes() {
	// Users endpoints
	ms.Mux.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) {
		// Check auth
		if !ms.checkAuth(w, r) {
			return
		}

		switch r.Method {
		case http.MethodGet:
			ms.writeJSON(w, http.StatusOK, MockListResponse(MockUsers, false, ""))
		case http.MethodPost:
			ms.writeJSON(w, http.StatusCreated, MockUser)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	ms.Mux.HandleFunc("/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		if !ms.checkAuth(w, r) {
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/v1/users/")

		switch r.Method {
		case http.MethodGet:
			if id == "usr_nonexistent" {
				ms.writeJSON(w, http.StatusNotFound, MockErrorNotFound)
				return
			}
			ms.writeJSON(w, http.StatusOK, MockUser)
		case http.MethodPatch:
			if id == "usr_nonexistent" {
				ms.writeJSON(w, http.StatusNotFound, MockErrorNotFound)
				return
			}
			ms.writeJSON(w, http.StatusOK, MockUser)
		case http.MethodDelete:
			if id == "usr_nonexistent" {
				ms.writeJSON(w, http.StatusNotFound, MockErrorNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Projects endpoints
	ms.Mux.HandleFunc("/v1/projects", func(w http.ResponseWriter, r *http.Request) {
		if !ms.checkAuth(w, r) {
			return
		}

		switch r.Method {
		case http.MethodGet:
			ms.writeJSON(w, http.StatusOK, MockListResponse(MockProjects, false, ""))
		case http.MethodPost:
			ms.writeJSON(w, http.StatusCreated, MockProject)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	ms.Mux.HandleFunc("/v1/projects/", func(w http.ResponseWriter, r *http.Request) {
		if !ms.checkAuth(w, r) {
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/v1/projects/")

		switch r.Method {
		case http.MethodGet:
			if id == "prj_nonexistent" {
				ms.writeJSON(w, http.StatusNotFound, MockErrorNotFound)
				return
			}
			ms.writeJSON(w, http.StatusOK, MockProject)
		case http.MethodDelete:
			if id == "prj_nonexistent" {
				ms.writeJSON(w, http.StatusNotFound, MockErrorNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Storage endpoints
	ms.Mux.HandleFunc("/v1/storage/buckets", func(w http.ResponseWriter, r *http.Request) {
		if !ms.checkAuth(w, r) {
			return
		}

		switch r.Method {
		case http.MethodGet:
			ms.writeJSON(w, http.StatusOK, MockListResponse(MockBuckets, false, ""))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health endpoint
	ms.Mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ms.writeJSON(w, http.StatusOK, map[string]string{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})
}

func (ms *MockServer) checkAuth(w http.ResponseWriter, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		ms.writeJSON(w, http.StatusUnauthorized, MockErrorAuth)
		return false
	}
	return true
}

func (ms *MockServer) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ErrorMockServer creates a server that returns errors for testing error handling
type ErrorMockServer struct {
	*MockServer
}

// NewErrorMockServer creates a mock server configured for error testing
func NewErrorMockServer() *ErrorMockServer {
	ms := &ErrorMockServer{
		MockServer: &MockServer{
			handlers: make(map[string]http.HandlerFunc),
		},
	}

	mux := http.NewServeMux()
	ms.Mux = mux

	return ms
}

// SetRateLimitError configures the server to return rate limit errors
func (ems *ErrorMockServer) SetRateLimitError() {
	ems.Mux.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "60")
		ems.writeJSON(w, http.StatusTooManyRequests, MockErrorRateLimit)
	})
	ems.Server = httptest.NewServer(ems.Mux)
}

// SetServerError configures the server to return 500 errors
func (ems *ErrorMockServer) SetServerError() {
	ems.Mux.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) {
		ems.writeJSON(w, http.StatusInternalServerError, MockErrorServer)
	})
	ems.Server = httptest.NewServer(ems.Mux)
}
