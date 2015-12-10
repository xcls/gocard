package app

import (
	"net/http"

	"github.com/mcls/gocard/stores/common"
)

func ReviewHandler(rc *RequestContext) error {
	crs, err := rc.Store.CardReviews.AllByUserID(rc.CurrentUser.ID)
	if err != nil {
		return err
	}
	var currentReview *common.CardReview
	remainingReviews := []*common.CardReview{}
	if len(crs) > 0 {
		currentReview = crs[0]
	}
	if len(crs) > 1 {
		remainingReviews = crs[1:]
	}
	return rc.HTML(http.StatusOK, "reviews/index", tplVars{
		"CurrentReview": currentReview,
		"CardReviews":   remainingReviews,
	})
}
