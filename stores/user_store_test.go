package stores

import (
	"testing"

	"github.com/mcls/gocard/stores/common"
	"github.com/mcls/gocard/stores/psql"
	"github.com/mcls/gocard/testutil"
)

func setupUserStores(t *testing.T) []*common.Store {
	db := testutil.ConnectDB(t)
	testutil.ResetDB(t, db)
	return []*common.Store{psql.NewStore(db)}
}

func TestUserStoreInsert(t *testing.T) {
	for _, store := range setupUserStores(t) {
		user := &common.User{
			Email:             "maartencls@gmail.com",
			EncryptedPassword: "$ecret$tuff",
		}

		err := store.Users.Insert(user)
		if err != nil {
			t.Fatalf("Didn't expect error, but got: %q", err)
		}
		found, err := store.Users.Find(user.ID)
		if err != nil {
			t.Fatalf("Didn't expect error, but got: %q", err)
		}
		if found.ID != user.ID {
			t.Fatalf("IDs not equal: %v != %v", found.ID, user.ID)
		}
	}
}

func TestUserStoreInsert_CantDuplicateEmail(t *testing.T) {
	for _, store := range setupUserStores(t) {
		u1 := &common.User{
			Email:             "maartencls@gmail.com",
			EncryptedPassword: "$ecret$tuff",
		}
		u2 := &common.User{
			Email:             "maartencls@gmail.com",
			EncryptedPassword: "$ecret$tuff2",
		}

		err := store.Users.Insert(u1)
		if err != nil {
			t.Fatalf("Didn't expect error, but got: %q", err)
		}
		err = store.Users.Insert(u2)
		if err == nil {
			t.Fatalf("Expected error but got nil", err)
		}
	}
}

func TestUserStoreAuthenticate(t *testing.T) {
	for _ = range setupUserStores(t) {
		user := &common.User{Email: "maartencls@gmail.com"}
		err := user.SetPassword("secret")
		if err != nil {
			t.Fatal(err)
		}
		err = store.Users.Insert(user)
		if err != nil {
			t.Fatal(err)
		}
		found, err := store.Users.Authenticate("maartencls@gmail.com", "secret")
		if err != nil {
			t.Fatal(err)
		}
		if found.ID != user.ID {
			t.Fatalf("IDs not equal: %v != %v", found.ID, user.ID)
		}

		_, err = store.Users.Authenticate("maartencls@gmail.com", "secret123")
		if err == nil {
			t.Fatal("Expected error because passwords don't match")
		}
	}
}
