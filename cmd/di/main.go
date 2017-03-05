package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/nylar/di"
)

func main() {
	req := &di.Request{
		Toggler: &di.FeatureToggle{
			Randomiser: rand.Float64,
		},
	}

	http.Handle("/", req)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatalf(err.Error())
	}
}
