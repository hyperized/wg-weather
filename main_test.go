package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var graphqlResponse = `{"data":{"getCityByName":{"weather":{"summary":{"title":"Clouds","description":"overcast clouds","icon":"04n"},"temperature":{"actual":275.16,"feelsLike":269.6,"min":273.69,"max":275.96},"wind":{"speed":7.55,"deg":331},"clouds":{"all":86,"visibility":10000,"humidity":62}}}}}`

func TestWeatherInTown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(graphqlResponse))
	}))
	defer server.Close()

	t.Log("Doing request")
	res, _ := weatherInTown(context.Background(), server.URL, "Voorhout")

	_, err := json.Marshal(&res)
	if err != nil {
		t.Fail()
	}
}

// TODO: Turn this into a test struct with multiple conditions, including invalid return values and various server reponses.
