// PAckage goreddit provides a Reddit client which wraps a lot of the API calls
// into helpful functions
package goreddit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Reddit struct {
	Username    string
	Password    string
	AppId       string
	AppSecret   string
	UserAgent   string
	AccessToken string
	Client      *http.Client
}

// ListingOpt defines options when querying for a list resource (list of links, list of comments, etc)
type ListingOpt struct {
	Limit   int
	Count   int
	Before  string
	After   string
	Show    string
	Depth   int
	Sort    string
	Comment string
}

// New creates a new client and tries to log in to Reddit.
func New(username string, password string, id string, secret string, ua string) (client *Reddit, err error) {
	client = &Reddit{
		username,
		password,
		id,
		secret,
		ua,
		"",
		&http.Client{},
	}
	err = client.Login()
	return
}

// Login tries to authenticate with Oauth using the provided information, if successful
// an access token will be obtained.
func (r *Reddit) Login() (err error) {
	fmt.Printf("Trying to connect to Reddit as %s\n", r.Username)

	form := url.Values{
		"grant_type": {"password"},
		"username":   {r.Username},
		"password":   {r.Password},
	}

	request, err := http.NewRequest(
		"POST",
		"https://www.reddit.com/api/v1/access_token",
		strings.NewReader(form.Encode()),
	)

	if err != nil {
		return
	}

	request.Header.Set("User-agent", r.UserAgent)
	request.SetBasicAuth(r.AppId, r.AppSecret)

	data, err := r.JsonResponse(request)
	if err != nil {
		return
	}

	var obj map[string]interface{}
	err = json.Unmarshal(data, &obj)

	if token, ok := obj["access_token"]; ok {
		r.AccessToken = token.(string)
	} else {
		err = errors.New("Couldn't acquire an access token")
	}

	return
}

// Request creates a new http.Request with the default values used by every call
// to the Reddit API (like the User-agent and token). It should be used after
// acquiring an access token.
func (r *Reddit) Request(requestType string, call string, form url.Values) (*http.Request, error) {
	body := strings.NewReader(form.Encode())
	request, err := http.NewRequest(requestType, "https://oauth.reddit.com"+call, body)
	if err != nil {
		return nil, err
	}

	// If the request is GET, the form data should be sent in the URL
	if requestType == "GET" {
		request.URL.RawQuery = form.Encode()
	}

	// The headers we'll be needing in every request
	request.Header.Set("User-agent", r.UserAgent)
	request.Header.Set("Authorization", "bearer "+r.AccessToken)

	return request, nil
}

// JsonResponse takes an *http.Request and make a request to Reddit using the
// default client. Then the JSON body from http.Response is parsed and analyzed.
// If there are no errors in the JSON object (including API-related errors as
// "too many requests" blocks and such) a []byte representing the JSON data is
// returned.
func (r *Reddit) JsonResponse(request *http.Request) (result []byte, err error) {
	response, err := r.Client.Do(request)
	if err != nil {
		return
	}

	result, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	var dataTemp interface{}
	var data map[string]interface{}

	err = json.Unmarshal(result, &dataTemp)
	if err != nil {
		return
	}

	switch v := dataTemp.(type) {
	case []interface{}:
		data = v[1].(map[string]interface{})
		result, err = json.Marshal(data)
		if err != nil {
			return
		}
	default:
		data = dataTemp.(map[string]interface{})
	}

	// Catch the HTTP errors (404, 403, etc)
	if _, exists := data["error"]; exists {
		err = errors.New(fmt.Sprintf("%v %v", data["error"], data["message"]))
		return
	}

	// API errors (blocks due too many requests, etc)
	if v, exists := data["json"].(map[string]interface{}); exists {
		val := v["errors"].([]interface{})
		if len(val) > 0 {
			apiErr := val[0].([]interface{})
			err = errors.New(fmt.Sprintf("[%v] %v", apiErr[0], apiErr[1]))
			return
		}
	}

	return
}

// Values returns a url.Values with the current ListingOpt's values
func (l ListingOpt) Values() url.Values {
	val := url.Values{}

	if l.Limit != 0 {
		val.Add("limit", fmt.Sprintf("%d", l.Limit))
	}
	if l.Count != 0 {
		val.Add("count", fmt.Sprintf("%d", l.Count))
	}
	if l.Before != "" {
		val.Add("before", l.Before)
	}
	if l.After != "" {
		val.Add("after", l.After)
	}
	if l.Show != "" {
		val.Add("show", l.Show)
	}
	if l.Depth != 0 {
		val.Add("depth", fmt.Sprintf("%d", l.Depth))
	}
	if l.Sort != "" {
		val.Add("sort", l.Sort)
	}
	if l.Comment != "" {
		val.Add("comment", l.Comment)
	}

	return val
}
