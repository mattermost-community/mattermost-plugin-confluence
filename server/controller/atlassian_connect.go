package controller

import (
	"html/template"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/util"
)

var atlassianConnectJSON = &Endpoint{
	Path:         "/atlassian-connect.json",
	Method:       http.MethodGet,
	Execute:      renderAtlassianConnectJSON,
	RequiresAuth: false,
}

func renderAtlassianConnectJSON(w http.ResponseWriter, r *http.Request) {
	conf := config.GetConfig()
	if status, err := verifyHTTPSecret(conf.Secret, r.FormValue("secret")); err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	bundlePath, err := config.Mattermost.GetBundlePath()
	if err != nil {
		config.Mattermost.LogWarn("Failed to get bundle path.", "Error", err.Error())
		return
	}

	templateDir := filepath.Join(bundlePath, "assets")
	tmplPath := path.Join(templateDir, "atlassian-connect.json")
	values := map[string]string{
		"BaseURL":      util.GetPluginURL(),
		"RouteACJSON":  "/api/v1/atlassian-connect.json?secret=" + url.QueryEscape(conf.Secret),
		"ExternalURL":  util.GetSiteURL(),
		"PluginKey":    util.GetPluginKey(),
		"SharedSecret": url.QueryEscape(conf.Secret),
	}
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = tmpl.Execute(w, values)
	if err != nil {
		http.Error(w, "failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
