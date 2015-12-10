package common

import (
	"fmt"
	"testing"
)

func TestUpdateEF(t *testing.T) {
	xs := []struct {
		ef     float64
		rating int64
		want   float64
	}{
		{2.5, 0, 1.7},
		{2.5, 1, 1.96},
		{2.5, 2, 2.18},
		{2.5, 3, 2.36},
		{2.5, 4, 2.5},
		{2.5, 5, 2.5},

		{1.3, 0, 1.3},
		{1.3, 1, 1.3},
		{1.3, 2, 1.3},
		{1.3, 3, 1.3},
		{1.3, 4, 1.3},
		{1.3, 5, 1.4},

		{1.4, 0, 1.3},
		{1.4, 1, 1.3},
		{1.4, 2, 1.3},
		{1.4, 3, 1.3},
		{1.4, 4, 1.4},
		{1.4, 5, 1.5},

		{1.5, 0, 1.3},
		{1.5, 1, 1.3},
		{1.5, 2, 1.3},
		{1.5, 3, 1.36},
		{1.5, 4, 1.5},
		{1.5, 5, 1.6},

		{1.6, 0, 1.3},
		{1.6, 1, 1.3},
		{1.6, 2, 1.3},
		{1.6, 3, 1.46},
		{1.6, 4, 1.6},
		{1.6, 5, 1.7},
	}
	for _, x := range xs {
		got := UpdateEF(x.ef, x.rating)
		if fmt.Sprintf("%.4f", got) != fmt.Sprintf("%.4f", x.want) {
			t.Logf("UpdateEF(%f, %d) = %q but want %q",
				x.ef, x.rating, got, x.want)
			t.Fail()
		}
	}
}

func TestUpdateInterval(t *testing.T) {
	xs := []struct {
		ef   float64
		days int64
		want int64
	}{
		{1.3, 6, 8},
		{1.3, 8, 11},
		{1.3, 11, 15},

		{2.5, 30, 75},
	}
	for _, x := range xs {
		got := UpdateInterval(x.days, x.ef)
		if got != x.want {
			t.Logf("UpdateInterval(%d, %f) = %d but want %d",
				x.days, x.ef, got, x.want)
			t.Fail()
		}
	}
}
