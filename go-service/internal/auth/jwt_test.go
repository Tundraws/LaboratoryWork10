package auth

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	service := NewService("test-secret", time.Hour)

	token, err := service.GenerateToken("student", "integration-client")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := service.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims.Username != "student" {
		t.Fatalf("username = %q, want %q", claims.Username, "student")
	}

	if claims.Role != "integration-client" {
		t.Fatalf("role = %q, want %q", claims.Role, "integration-client")
	}
}

func TestParseTokenErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		token    string
		secret   string
		wantErr  bool
		setupTTL time.Duration
	}{
		{
			name:     "invalid format",
			token:    "not-a-jwt",
			secret:   "test-secret",
			wantErr:  true,
			setupTTL: time.Hour,
		},
		{
			name:     "expired token",
			secret:   "test-secret",
			wantErr:  true,
			setupTTL: -time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.secret, tt.setupTTL)
			token := tt.token
			if token == "" {
				var err error
				token, err = service.GenerateToken("student", "integration-client")
				if err != nil {
					t.Fatalf("GenerateToken() error = %v", err)
				}
			}

			_, err := service.ParseToken(token)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
