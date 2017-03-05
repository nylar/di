package di

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockFeatureToggle struct {
	inSplitPercentage []bool
	counter           int
}

func (m *mockFeatureToggle) SplitPercentage(percentage float64) bool {
	inSplit := m.inSplitPercentage[m.counter]
	m.counter++
	return inSplit
}

func TestFeatureToggle(t *testing.T) {
	tests := []struct {
		randomiser func() float64
		percentage float64
		expected   bool
	}{
		{
			randomiser: func() float64 { return 0.0 },
			percentage: 0,
			expected:   false,
		},
		{
			randomiser: func() float64 { return 0.0 },
			percentage: 100,
			expected:   true,
		},
		{
			randomiser: func() float64 { return 0.60 },
			percentage: 40,
			expected:   false,
		},
		{
			randomiser: func() float64 { return 0.60 },
			percentage: 80,
			expected:   true,
		},
	}

	for _, test := range tests {
		ft := &FeatureToggle{
			Randomiser: test.randomiser,
		}

		assert.Equal(t, test.expected, ft.SplitPercentage(test.percentage))
	}
}

func TestRequest_ServeHTTP_InPercentage(t *testing.T) {
	req := &Request{
		Toggler: &mockFeatureToggle{
			inSplitPercentage: []bool{true, true},
		},
	}

	w := httptest.NewRecorder()

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	req.ServeHTTP(w, r)

	assert.Equal(t, "In split percentage\nIn split sub percentage", w.Body.String())
}

func TestRequest_ServeHTTP_InPercentageOutSubPercentage(t *testing.T) {
	req := &Request{
		Toggler: &mockFeatureToggle{
			inSplitPercentage: []bool{true, false},
		},
	}

	w := httptest.NewRecorder()

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	req.ServeHTTP(w, r)

	assert.Equal(t, "In split percentage\nNot in sub split percentage", w.Body.String())
}

func TestRequest_ServeHTTP_NotInPercentage(t *testing.T) {
	req := &Request{
		Toggler: &mockFeatureToggle{
			inSplitPercentage: []bool{false},
		},
	}

	w := httptest.NewRecorder()

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err.Error())
	}

	req.ServeHTTP(w, r)

	assert.Equal(t, "Not in split percentage\n", w.Body.String())
}
