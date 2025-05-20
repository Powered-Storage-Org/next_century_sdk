package client

import (
	"log"
	"net/http"
)

func closRespBody(resp *http.Response) {
	if resp == nil {
		log.Println("Response is nil, nothing to close")
		return
	}

	if err := resp.Body.Close(); err != nil {
		log.Printf("Error closing response body: %v", err)
	}
}
