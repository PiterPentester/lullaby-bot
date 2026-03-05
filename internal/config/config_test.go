package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name: "Valid config",
			envs: map[string]string{
				"TELEGRAM_BOT_TOKEN":  "test_token",
				"AUTHORIZED_USER_IDS": "123,456",
			},
			want: &Config{
				Token:             "test_token",
				AuthorizedUserIDs: []int64{123, 456},
			},
			wantErr: false,
		},
		{
			name: "Missing token",
			envs: map[string]string{
				"AUTHORIZED_USER_IDS": "123",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Missing IDs",
			envs: map[string]string{
				"TELEGRAM_BOT_TOKEN": "test_token",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid ID format",
			envs: map[string]string{
				"TELEGRAM_BOT_TOKEN":  "test_token",
				"AUTHORIZED_USER_IDS": "123,abc",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.envs {
				os.Setenv(k, v)
			}

			got, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_IsAuthorized(t *testing.T) {
	cfg := &Config{
		AuthorizedUserIDs: []int64{123, 456},
	}

	tests := []struct {
		userID int64
		want   bool
	}{
		{123, true},
		{456, true},
		{789, false},
	}

	for _, tt := range tests {
		got := cfg.IsAuthorized(tt.userID)
		if got != tt.want {
			t.Errorf("IsAuthorized(%d) = %v, want %v", tt.userID, got, tt.want)
		}
	}
}
