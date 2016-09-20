package app

import (
	"net/http"
	"time"

	"github.com/xcls/gocard/stores/common"
)

type AnswerForm struct {
	Rating     int64
	ReviewID   int64
	CardReview *common.CardReview
}

func NewAnswerForm(cr *common.CardReview, rating int64) *AnswerForm {
	return &AnswerForm{
		Rating:     rating,
		ReviewID:   cr.ID,
		CardReview: cr,
	}
}

func ReviewHandler(rc *RequestContext) error {
	crs, err := rc.Store.CardReviews.DueAt(rc.CurrentUser.ID, time.Now())
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

	var answers []*AnswerForm
	if currentReview != nil {
		answers = make([]*AnswerForm, common.MaxRating+1)
		for i := common.MinRating; i <= common.MaxRating; i++ {
			answers[i] = NewAnswerForm(currentReview, int64(i))
		}
	}
	return rc.HTML(http.StatusOK, "reviews/index", tplVars{
		"Answers":       answers,
		"CurrentReview": currentReview,
		"CardReviews":   remainingReviews,
	})
}

func AnswerReviewHandler(rc *RequestContext) error {
	form := new(AnswerForm)
	if err := decodeForm(form, rc.Request); err != nil {
		return err
	}
	applog.Printf("Gave a rating of %d", form.Rating)
	// TODO Check correct user etc..
	if err := rc.Store.AnswerReview(form.ReviewID, form.Rating); err != nil {
		return err
	}
	return rc.Redirect("/review")
}
