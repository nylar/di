package di

import "net/http"

type FeatureToggler interface {
	SplitPercentage(percentage float64) bool
}

type FeatureToggle struct {
	Randomiser func() float64
}

func (ft *FeatureToggle) SplitPercentage(percentage float64) bool {
	switch percentage {
	case 0:
		return false
	case 100:
		return true
	default:
		return ft.Randomiser()*100 < percentage
	}
}

type Request struct {
	Toggler FeatureToggler
}

func (req *Request) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if req.Toggler.SplitPercentage(50) {
		w.Write([]byte("In split percentage\n"))

		if !req.Toggler.SplitPercentage(25) {
			w.Write([]byte("Not in sub split percentage"))
		} else {
			w.Write([]byte("In split sub percentage"))
		}
	} else {
		w.Write([]byte("Not in split percentage\n"))
	}
}
