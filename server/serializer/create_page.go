package serializer

import "encoding/json"

type PageDetails struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func PageDetailsFromJSON(body []byte) (*PageDetails, error) {
	var pd PageDetails
	if err := json.Unmarshal(body, &pd); err != nil {
		return nil, err
	}
	return &pd, nil
}
