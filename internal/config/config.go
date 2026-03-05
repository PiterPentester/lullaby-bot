package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Token            string
	AuthorizedUserIDs []int64
}

func Load() (*Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	userIDsStr := os.Getenv("AUTHORIZED_USER_IDS")
	if userIDsStr == "" {
		return nil, fmt.Errorf("AUTHORIZED_USER_IDS environment variable is required (comma-separated list of IDs)")
	}

	var userIDs []int64
	for _, idStr := range strings.Split(userIDsStr, ",") {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID %q: %v", idStr, err)
		}
		userIDs = append(userIDs, id)
	}

	return &Config{
		Token:            token,
		AuthorizedUserIDs: userIDs,
	}, nil
}

func (c *Config) IsAuthorized(userID int64) bool {
	for _, id := range c.AuthorizedUserIDs {
		if id == userID {
			return true
		}
	}
	return false
}
