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

		tests := []*struct {
			Email         string
			Pass          string
			ExpectCorrect bool
		}{
			{"maartencls@gmail.com", "secret", true},
			{"MaartenCls@gmail.com", "secret", true},
			{"maartencls@gmail.com", "notsecret", false},
		}

		for _, test := range tests {
			found, err := store.Users.Authenticate(
				test.Email,
				test.Pass,
			)
			if test.ExpectCorrect {
				// Correct password
				if err != nil {
					t.Fatal(err)
				}
				if found.ID != user.ID {
					t.Fatalf("IDs not equal: %v != %v", found.ID, user.ID)
				}
			} else {
				// Wrong password
				if err == nil {
					t.Fatal("Expected error because passwords don't match")
				}
				if err != common.ErrUserAuthFailed {
					t.Fatalf("Expected ErrUserAuthFailed, but got: %q", err)
				}
			}
		}
	}
}

func newUser(email string, optPass ...string) *common.User {
	user := &common.User{Email: email, EncryptedPassword: "NOT_SET"}
	if len(optPass) > 0 {
		user.SetPassword(optPass[0])
	}
	return user
}

func TestUserStoreFindByEmail(t *testing.T) {
	for _ = range setupUserStores(t) {
		var err error
		u1 := newUser("maartencls@gmail.com")
		u2 := newUser("quintencls@gmail.com")
		err = store.Users.Insert(u1)
		if err != nil {
			t.Fatal(err)
		}
		err = store.Users.Insert(u2)
		if err != nil {
			t.Fatal(err)
		}

		tests := []*struct {
			email string
			want  *common.User
		}{
			{"maartencls@gmail.com", u1},
			{"MAARTENCLS@GMAIL.COM", u1},
			{"MaarTeNCls@gMAIL.COM", u1},
			{"quintencls@gmail.com", u2},
			{"example@gmail.com", nil},
		}

		for _, o := range tests {
			found, err := store.Users.FindByEmail(o.email)

			if o.want == nil {
				// Not found
				if err == nil {
					t.Fatal("Expected error")
				}
				continue
			}

			// Found
			if err != nil {
				t.Fatal(err)
			}
			if found.ID != o.want.ID {
				t.Fatalf("IDs not equal: %v != %v", found.ID, o.want.ID)
			}
		}
	}
}
