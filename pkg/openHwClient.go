package pkg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

func GetMetrics(port string) (*ResponseDto, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("http://localhost:%s/data.json", port))

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err := client.DoTimeout(req, resp, 1*time.Second)
	if err != nil {
		return nil, err
	}
	bodyBytes := resp.Body()
	responseDto := &ResponseDto{}
	err = json.Unmarshal(bodyBytes, &responseDto)
	return responseDto, err
}

type ResponseDto struct {
	ID       int    `json:"id"`
	Text     string `json:"Text"`
	Children []struct {
		ID       int    `json:"id"`
		Text     string `json:"Text"`
		Children []struct {
			ID       int    `json:"id"`
			Text     string `json:"Text"`
			Children []struct {
				ID       int    `json:"id"`
				Text     string `json:"Text"`
				Children []struct {
					ID       int    `json:"id"`
					Text     string `json:"Text"`
					Children []struct {
						ID       int           `json:"id"`
						Text     string        `json:"Text"`
						Children []interface{} `json:"Children"`
						Min      string        `json:"Min"`
						Value    string        `json:"Value"`
						Max      string        `json:"Max"`
						ImageURL string        `json:"ImageURL"`
					} `json:"Children"`
					Min      string `json:"Min"`
					Value    string `json:"Value"`
					Max      string `json:"Max"`
					ImageURL string `json:"ImageURL"`
				} `json:"Children"`
				Min      string `json:"Min"`
				Value    string `json:"Value"`
				Max      string `json:"Max"`
				ImageURL string `json:"ImageURL"`
			} `json:"Children"`
			Min      string `json:"Min"`
			Value    string `json:"Value"`
			Max      string `json:"Max"`
			ImageURL string `json:"ImageURL"`
		} `json:"Children"`
		Min      string `json:"Min"`
		Value    string `json:"Value"`
		Max      string `json:"Max"`
		ImageURL string `json:"ImageURL"`
	} `json:"Children"`
	Min      string `json:"Min"`
	Value    string `json:"Value"`
	Max      string `json:"Max"`
	ImageURL string `json:"ImageURL"`
}
