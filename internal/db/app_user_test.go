package db

import (
	"testing"

	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
)

func Test_Integration_Read_AppUsers(t *testing.T) {
	cfg, err := helpers.LoadAppConfig(".", "config")

	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	_, err = NewClientWithConfig(cfg)
	if err != nil {
		t.Fatalf("failed to create mongo client: %v", err)
	}

	u, err := AppUsers()
	if err != nil {
		t.Fatalf("failed to read collection '%s': %v", collectionAppUser, err)
	}

	if len(u) <= 0 {
		t.Fatalf("invalid number of %s records", collectionAppUser)
	}

	for _, v := range u {
		if v.ID.Hex() == "" {
			t.Fatalf("invalid record ID")
		}
	}
}
