package controller

import (
	"io/ioutil"
	"net/http"

	"github.com/mattermost/mattermost-server/model"
)

var webhookEndpoint = &Endpoint{
	Path:         "/webhook",
	Method:       "POST",
	Execute:      handleWebhookEvent,
	RequiresAuth: false,
}

func handleWebhookEvent(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	bodyString := string(body)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(model.MapToJson(map[string]string{"status": "ok", "body": bodyString})))
}
