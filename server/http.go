package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

const (
	routePrefixInstance = "instance"
)

func respondErr(w http.ResponseWriter, code int, err error) (int, error) {
	http.Error(w, err.Error(), code)
	return code, err
}

func (p *Plugin) loadTemplates(dir string) (map[string]*template.Template, error) {
	dir = filepath.Clean(dir)
	templates := make(map[string]*template.Template)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		template, err := template.ParseFiles(path)
		if err != nil {
			p.errorf("OnActivate: failed to parse template %s: %v", path, err)
			return nil
		}
		key := path[len(dir):]
		templates[key] = template
		p.debugf("loaded template %s", key)
		return nil
	})
	if err != nil {
		return nil, errors.WithMessage(err, "OnActivate: failed to load templates")
	}
	return templates, nil
}

func splitInstancePath(route string) (instanceURL string, remainingPath string) {
	leadingSlash := ""
	ss := strings.Split(route, "/")
	if len(ss) > 1 && ss[0] == "" {
		leadingSlash = "/"
		ss = ss[1:]
	}

	if len(ss) < 2 {
		return "", route
	}

	if ss[0] == "api" && strings.Contains(ss[1], "v") {
		ss = ss[2:]
	}

	if ss[0] != routePrefixInstance {
		return route, route
	}

	id, err := decode(ss[1])
	if err != nil {
		return "", route
	}
	return string(id), leadingSlash + strings.Join(ss[2:], "/")
}

func (p *Plugin) respondSpecialTemplate(w http.ResponseWriter, key string, status int, contentType string, values interface{}) (int, error) {
	w.Header().Set("Content-Type", contentType)
	t := p.templates[key]
	if t == nil {
		return respondErr(w, http.StatusInternalServerError,
			errors.New("no template found for "+key))
	}
	err := t.Execute(w, values)
	if err != nil {
		return http.StatusInternalServerError,
			errors.WithMessage(err, "failed to write response")
	}
	return status, nil
}

func (p *Plugin) respondTemplate(w http.ResponseWriter, r *http.Request, contentType string, values interface{}) (int, error) {
	_, path := splitInstancePath(r.URL.Path)
	w.Header().Set("Content-Type", contentType)
	t := p.templates[path]
	if t == nil {
		return respondErr(w, http.StatusInternalServerError,
			errors.New("no template found for "+path))
	}
	err := t.Execute(w, values)
	if err != nil {
		return http.StatusInternalServerError, errors.WithMessage(err, "failed to write response")
	}
	return http.StatusOK, nil
}
