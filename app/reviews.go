package app

import (
	"net/http"

	"github.com/mcls/gocard/stores/common"
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
	crs, err := rc.Store.CardReviews.EnabledByUserID(rc.CurrentUser.ID)
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
		answers = make([]*AnswerForm, 5)
		for i := 0; i < 5; i++ {
			rating := int64(i + 1)
			answers[i] = NewAnswerForm(currentReview, rating)
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
	return rc.Redirect("/review")
}
