package internal

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// call the mlb api (or well... the endpoint)
func CallAPI(url string) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		return make([]byte, 0), getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	return ioutil.ReadAll(res.Body)
}
