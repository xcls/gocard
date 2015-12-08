package stores

import (
	"log"

	"github.com/mcls/gocard/dbutil"
	"github.com/mcls/gocard/stores/common"
	"github.com/mcls/gocard/stores/psql"
)

var Store *common.Store

func init() {
	db, err := dbutil.Connect()
	if err != nil {
		log.Fatal(err)
	}
	Store = psql.NewStore(db)
}
