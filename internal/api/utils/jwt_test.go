package utils_test

import (
	"chat-service/internal/api/utils"
	"testing"
)

const (
	uid   = "6328040e-6a50-49b3-92fc-4d31e53c2dab"
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjMyODA0MGUtNmE1MC00OWIzLTkyZmMtNGQzMWU1M2MyZGFiIn0.FaVh671uUYtFwAkeXrmhMF5HIoBnPlwdd5d9J23YABk"
)

func TestGenerateJWT(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		want    string
		wantErr bool
	}{
		{
			name:    "GenerateJWT",
			userID:  uid,
			want:    token,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := utils.GenerateJWT(tt.userID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GenerateJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GenerateJWT() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("GenerateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractUserIDFromJWT(t *testing.T) {
	tests := []struct {
		name        string
		tokenString string
		want        string
		wantErr     bool
	}{
		{
			name:        "ExtractUserIDFromJWT",
			tokenString: token,
			want:        uid,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := utils.ExtractUserIDFromJWT(tt.tokenString)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExtractUserIDFromJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ExtractUserIDFromJWT() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ExtractUserIDFromJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserIDFromHeader(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		want       string
		wantErr    bool
	}{
		{
			name:       "GetUserIDFromHeader",
			authHeader: "Bearer " + token,
			want:       uid,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := utils.GetUserIDFromHeader(tt.authHeader)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetUserIDFromHeader() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetUserIDFromHeader() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("GetUserIDFromHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
