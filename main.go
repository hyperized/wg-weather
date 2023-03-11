package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Request struct {
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
	Query         string                 `json:"query"`
}

var (
	timeout = time.Second * 10
	port    = ":8080"
	town    = "Voorhout"
	baseUrl = "https://weather-api.wundergraph.com/"
	query   = `query WeatherInTown {
				  getCityByName(name: "%s", config: {units: metric, lang: nl}) {
					weather {
					  summary {
						title
						description
						icon
					  }
					  temperature {
						actual
						feelsLike
						min
						max
					  }
					  wind {
						speed
						deg
					  }
					  clouds {
						all
						visibility
						humidity
					  }
					}
				  }
				}
			`
)

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		timer       = time.NewTimer(timeout)
		mux         = http.NewServeMux()
	)
	defer cancel()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		timer.Reset(timeout)
		log.Printf("allowing the http server to live another %s seconds...", timeout.String())

		w.Header().Set("Content-Type", "application/json")
		results, httpStatus := weatherInTown(ctx, baseUrl, town)
		w.WriteHeader(httpStatus)
		_, err := w.Write(results)
		if err != nil {
			log.Fatalf("could not write response data to the client because: %s", err.Error())
		}
	})
	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	go func() {
		log.Printf("Serving weather reports for %s on %s...", town, port)
		log.Print(srv.ListenAndServe())
	}()

	go func() {
		<-timer.C // Timer expired
		log.Println("timer expired...")
		cancel() // Cancel current context.
	}()

	// When the context is down, gracefully shut down the http server and exit
	select {
	case <-ctx.Done():
		log.Println("shutting down http server...")
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(0)
}

func weatherInTown(ctx context.Context, url, town string) ([]byte, int) {
	var (
		req = &Request{
			OperationName: "WeatherInTown",
			Variables:     map[string]interface{}{}, // Left empty
			Query:         fmt.Sprintf(query, town),
		}
		httpClient = &http.Client{
			Timeout: time.Second * 9, // Always time out before our http server context
		}
	)

	requestBody, err := json.Marshal(&req)
	if err != nil {
		return []byte(`{"error": "could not marshal query to JSON", "reason": "` + err.Error() + `"`), http.StatusInternalServerError
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return []byte(`{"error": "could not shape http request", "reason": "` + err.Error() + `"`), http.StatusInternalServerError
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	result, err := httpClient.Do(httpRequest)
	if err != nil {
		return []byte(`{"error": "could not execute http request", "reason": "` + err.Error() + `"`), http.StatusInternalServerError
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("could not close http request body")
		}
	}(result.Body)

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return []byte(`{"error": "could not read body contents", "reason": "` + err.Error() + `"`), http.StatusInternalServerError
	}

	return body, http.StatusOK
}
