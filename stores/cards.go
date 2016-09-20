package stores

import (
	"log"

	"github.com/xcls/gocard/dbutil"
	"github.com/xcls/gocard/stores/common"
	"github.com/xcls/gocard/stores/psql"
)

var Store *common.Store

func init() {
	db, err := dbutil.Connect()
	if err != nil {
		log.Fatal(err)
	}
	Store = psql.NewStore(db)
}
