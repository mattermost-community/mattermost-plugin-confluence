package serializer

import "encoding/json"

type PageDetails struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func PageDetailsFromJSON(body []byte) (*PageDetails, error) {
	var pageDetails PageDetails
	if err := json.Unmarshal(body, &pageDetails); err != nil {
		return nil, err
	}
	return &pageDetails, nil
}
