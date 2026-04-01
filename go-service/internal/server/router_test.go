package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"laboratorywork10/go-service/internal/auth"
)

func TestIssueToken(t *testing.T) {
	router := NewRouter(auth.NewService("test-secret", time.Hour))

	tests := []struct {
		name       string
		body       map[string]any
		wantStatus int
	}{
		{
			name: "valid credentials",
			body: map[string]any{
				"username": "student",
				"password": "securepass123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "validation error",
			body: map[string]any{
				"username": "st",
				"password": "short",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "wrong credentials",
			body: map[string]any{
				"username": "student",
				"password": "wrongpass123",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/token", bytes.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, tt.wantStatus, recorder.Body.String())
			}
		})
	}
}

func TestProtectedProcessEndpoint(t *testing.T) {
	authService := auth.NewService("test-secret", time.Hour)
	router := NewRouter(authService)

	token, err := authService.GenerateToken("student", "integration-client")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	validBody := map[string]any{
		"request_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		"customer":   "Ivan Petrov",
		"address": map[string]any{
			"city":     "Moscow",
			"street":   "Tverskaya 1",
			"zip_code": "123456",
		},
		"items": []map[string]any{
			{
				"name":     "keyboard",
				"quantity": 2,
				"price":    1500.50,
			},
		},
		"metadata": map[string]any{
			"priority": "high",
			"tags":     []string{"study", "api"},
		},
	}

	tests := []struct {
		name       string
		token      string
		body       map[string]any
		wantStatus int
	}{
		{
			name:       "success",
			token:      token,
			body:       validBody,
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing token",
			body:       validBody,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:  "validation error",
			token: token,
			body: map[string]any{
				"request_id": "not-uuid",
				"customer":   "A",
				"address": map[string]any{
					"city":     "",
					"street":   "st",
					"zip_code": "12",
				},
				"items": []map[string]any{},
				"metadata": map[string]any{
					"priority": "urgent",
					"tags":     []string{},
				},
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/process", bytes.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, tt.wantStatus, recorder.Body.String())
			}
		})
	}
}

func TestExtractBearerToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		header    string
		wantToken string
		wantErr   bool
	}{
		{name: "valid", header: "Bearer token-value", wantToken: "token-value"},
		{name: "missing prefix", header: "token-value", wantErr: true},
		{name: "empty token", header: "Bearer   ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := extractBearerToken(tt.header)
			if (err != nil) != tt.wantErr {
				t.Fatalf("extractBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if token != tt.wantToken {
				t.Fatalf("token = %q, want %q", token, tt.wantToken)
			}
		})
	}
}
