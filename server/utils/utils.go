package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	html "github.com/levigross/exp-html"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-confluence/server/config"
)

const (
	Script              = "script"
	Style               = "style"
	ErrorStatusNotFound = "No content found"
	Someone             = "Someone"
)

// GetKeyHash can be used to create a hash from a string
func GetKeyHash(key string) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SplitArgs(s string) ([]string, error) {
	indexes := regexp.MustCompile("\"").FindAllStringIndex(s, -1)
	if len(indexes)%2 != 0 {
		return []string{}, errors.New("quotes not closed")
	}

	indexes = append([][]int{{0, 0}}, indexes...)

	if indexes[len(indexes)-1][1] < len(s) {
		indexes = append(indexes, [][]int{{len(s), 0}}...)
	}

	var args []string
	for i := 0; i < len(indexes)-1; i++ {
		start := indexes[i][1]
		end := Min(len(s), indexes[i+1][0])

		if i%2 == 0 {
			args = append(args, strings.Split(strings.Trim(s[start:end], " "), " ")...)
		} else {
			args = append(args, s[start:end])
		}
	}

	cleanedArgs := make([]string, len(args))
	count := 0

	for _, arg := range args {
		if arg != "" {
			cleanedArgs[count] = strings.TrimSpace(arg)
			count++
		}
	}

	return cleanedArgs[0:count], nil
}

func GetPluginKey() string {
	var regexNonAlphaNum = regexp.MustCompile("[^a-zA-Z0-9]+")
	return "mattermost_" + regexNonAlphaNum.ReplaceAllString(GetSiteURL(), "_")
}

func GetPluginURLPath() string {
	return "/plugins/" + config.PluginName + "/api/v1"
}

func GetPluginURL() string {
	return strings.TrimRight(GetSiteURL(), "/") + GetPluginURLPath()
}

func GetSiteURL() string {
	ptr := config.Mattermost.GetConfig().ServiceSettings.SiteURL
	if ptr == nil {
		return ""
	}
	return *ptr
}

func GetAtlassianConnectURLPath() string {
	return "/atlassian-connect.json?secret=" + url.QueryEscape(config.GetConfig().Secret)
}

func GetConfluenceServerWebhookURLPath() string {
	return "/server/webhook?secret=" + url.QueryEscape(config.GetConfig().Secret)
}

func IsSystemAdmin(userID string) bool {
	user, appErr := config.Mattermost.GetUser(userID)
	if appErr != nil {
		return false
	}
	return user.IsInRole(model.SystemAdminRoleId)
}

func Deduplicate(a []string) []string {
	check := make(map[string]int)
	result := make([]string, 0)
	for _, val := range a {
		check[val] = 1
	}

	for key := range check {
		result = append(result, key)
	}

	return result
}

type ConfluenceStatus struct {
	State string `json:"state"`
}

func NormalizeConfluenceURL(confluenceURL string) (string, error) {
	u, err := url.Parse(confluenceURL)
	if err != nil {
		return "", err
	}
	if u.Host == "" {
		ss := strings.Split(u.Path, "/")
		if len(ss) > 0 && ss[0] != "" {
			u.Host = ss[0]
			u.Path = path.Join(ss[1:]...)
		}
		u, err = url.Parse(u.String())
		if err != nil {
			return "", err
		}
	}
	if u.Host == "" {
		return "", errors.Errorf("Invalid URL, no hostname: %q", confluenceURL)
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	confluenceURL = strings.TrimSuffix(u.String(), "/")
	return confluenceURL, nil
}

// CheckConfluenceURL checks if `/status` endpoint of the Confluence URL is accessible
// and responding with the correct state which is "RUNNING"
func CheckConfluenceURL(mattermostSiteURL, confluenceURL string, requireHTTPS bool) (_ string, err error) {
	confluenceURL, err = NormalizeConfluenceURL(confluenceURL)
	if err != nil {
		return "", err
	}
	if confluenceURL == strings.TrimSuffix(mattermostSiteURL, "/") {
		return "", errors.Errorf("%s is the Mattermost site URL. Please use your Confluence URL", confluenceURL)
	}

	defer func() {
		if err != nil {
			err = errors.Wrap(err, "we couldn't validate the connection to your Confluence server. "+
				"This could be because of existing firewall or proxy rules, or because the URL was entered incorrectly")
		}
	}()

	resp, err := http.Get(confluenceURL + "/status")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("Confluence server returned http status code %q when checking for availability: %q", resp.Status, confluenceURL)
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var status ConfluenceStatus
	err = json.Unmarshal(resBody, &status)
	if err != nil {
		return "", err
	}
	if status.State != "RUNNING" {
		return "", errors.Errorf("Confluence server is not in correct state, it should be up and running: %q", confluenceURL)
	}
	return confluenceURL, nil
}

func IsConfluenceCloudURL(confluenceURL string) bool {
	u, err := url.Parse(confluenceURL)
	if err != nil {
		return false
	}
	return strings.HasSuffix(u.Hostname(), ".atlassian.net")
}

func GetEndpointURL(instanceURL, path string) (string, error) {
	endpointURL, err := url.Parse(strings.TrimSpace(fmt.Sprintf("%s%s", instanceURL, path)))
	if err != nil {
		return "", err
	}

	return endpointURL.String(), nil
}

func CallJSON(url, method, path string, in, out interface{}, httpClient *http.Client) (responseData []byte, statusCode int, err error) {
	contentType := "application/json"
	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(in)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return call(url, method, path, contentType, buf, out, httpClient)
}

func call(basePath, method, path, contentType string, inBody io.Reader, out interface{}, httpClient *http.Client) (responseData []byte, statusCode int, err error) {
	errContext := fmt.Sprintf("confluence: Call failed: method:%s, path:%s", method, path)
	pathURL, err := url.Parse(path)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.WithMessage(err, errContext)
	}

	if pathURL.Scheme == "" || pathURL.Host == "" {
		var baseURL *url.URL
		baseURL, err = url.Parse(basePath)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.WithMessage(err, errContext)
		}
		if path[0] != '/' {
			path = "/" + path
		}
		path = baseURL.String() + path
	}

	req, err := http.NewRequest(method, path, inBody)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.Body == nil {
		return nil, resp.StatusCode, nil
	}
	defer resp.Body.Close()

	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if out != nil {
			err = json.Unmarshal(responseData, out)
			if err != nil {
				return responseData, http.StatusInternalServerError, err
			}
		}
		return responseData, resp.StatusCode, nil

	case http.StatusNoContent:
		return nil, resp.StatusCode, nil

	case http.StatusNotFound:
		return nil, resp.StatusCode, errors.Errorf(ErrorStatusNotFound)
	}

	type ErrorResponse struct {
		Message string `json:"message"`
	}
	errResp := ErrorResponse{}
	err = json.Unmarshal(responseData, &errResp)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return responseData, resp.StatusCode, errors.New(errResp.Message)
}

func GetBodyForExcerpt(htmlBodyValue string) string {
	var str string
	domDocTest := html.NewTokenizer(strings.NewReader(htmlBodyValue))
	previousStartTokenTest := domDocTest.Token()
loopDomTest:
	for {
		tt := domDocTest.Next()
		switch {
		case tt == html.ErrorToken:
			break loopDomTest // End of the document,  done
		case tt == html.StartTagToken:
			previousStartTokenTest = domDocTest.Token()
		case tt == html.TextToken:
			if previousStartTokenTest.Data == Script || previousStartTokenTest.Data == Style {
				continue
			}
			TextContent := strings.TrimSpace(html.UnescapeString(string(domDocTest.Text())))
			if len(TextContent) > 0 {
				str = fmt.Sprintf("%s\n%s", str, TextContent)
			}
		}
	}
	return str
}

func CreateConfluenceURL(selfEventURL string) (string, error) {
	url, err := url.Parse(selfEventURL)
	if err != nil {
		return "", nil
	}
	return url.Scheme + "://" + url.Host, nil
}

func GetUsernameOrAnonymousName(username string) string {
	if username == "" {
		return "Someone"
	}
	return username
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func AddQueryParams(urlParam string, queryParams map[string]interface{}) (string, error) {
	url, err := url.Parse(urlParam)
	if err != nil {
		return "", err
	}
	q := url.Query()
	for key, value := range queryParams {
		q.Set(key, fmt.Sprintf("%v", value))
	}
	url.RawQuery = q.Encode()
	return url.String(), nil
}
