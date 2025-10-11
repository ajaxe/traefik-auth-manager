package db

import (
	"context"
	"testing"

	"github.com/ajaxe/traefik-auth-manager/internal/helpers"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Test_Integration_Read_AppUsers(t *testing.T) {
	cfg, err := helpers.LoadAppConfig(".", "config")

	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	_, err = newClientWithConfig(cfg)
	if err != nil {
		t.Fatalf("failed to create mongo client: %v", err)
	}

	da := NewAppUserDataAccess()(context.Background())
	u, err := da.AppUsers()
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

func Test_ExistingUser_OK(t *testing.T) {
	cfg, err := helpers.LoadAppConfig(".", "config")

	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	_, err = newClientWithConfig(cfg)
	if err != nil {
		t.Fatalf("failed to create mongo client: %v", err)
	}

	da := NewAppUserDataAccess()(context.Background())
	d, err := da.AppUserByUsername("admin")
	if err != nil {
		t.Fatalf("failed to read app user: %v", err)
	}

	if d == nil {
		t.Fatalf("expect non-nil app user 'admin' : %v", err)
	}
	if d.UserName != "admin" {
		t.Fatalf("expect app username to be 'admin' got '%v'", d.UserName)
	}
}
func Test_NonExistingUser_OK(t *testing.T) {
	cfg, err := helpers.LoadAppConfig(".", "config")

	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	_, err = newClientWithConfig(cfg)
	if err != nil {
		t.Fatalf("failed to create mongo client: %v", err)
	}

	da := NewAppUserDataAccess()(context.Background())
	_, err = da.AppUserByUsername("not_existing_user")

	if err != mongo.ErrNoDocuments {
		t.Fatalf("failed to read app user: %v", err)
	}
}
