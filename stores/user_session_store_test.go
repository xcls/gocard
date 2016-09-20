package stores

import (
	"testing"

	"github.com/xcls/gocard/stores/common"
	"github.com/xcls/gocard/stores/psql"
	"github.com/xcls/gocard/testutil"
)

func setupUserSessionStores(t *testing.T) []*common.Store {
	db := testutil.ConnectDB(t)
	testutil.ResetDB(t, db)
	return []*common.Store{psql.NewStore(db)}
}

func TestUserSessionStoreInsert(t *testing.T) {
	for _, store := range setupUserSessionStores(t) {
		user := newUser("maartencls@gmail.com")
		err := store.Users.Insert(user)
		if err != nil {
			t.Fatal(err)
		}

		session := common.NewUserSession(user.ID)
		err = store.UserSessions.Insert(session)
		if err != nil {
			t.Fatal(err)
		}

		found, err := store.UserSessions.Find(session.UID)
		if err != nil {
			t.Fatal(err)
		}

		if found.ID != session.ID {
			t.Fatalf("ID mismatch: \n%+v\n%+v", session, found)
		}
	}
}
