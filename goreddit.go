// PAckage goreddit provides a Reddit client which wraps a lot of the API calls
// into helpful functions
package goreddit

import (
  "fmt"
  "errors"
  "strings"
  "net/http"
  "net/url"
  "io/ioutil"
  "encoding/json"
)

type Reddit struct {
  Username string 
  Password string 
  AppId string
  AppSecret string
  UserAgent string
  AccessToken string
  Client *http.Client
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
    "username": {r.Username},
    "password": {r.Password},
  }

  request,err := http.NewRequest(
    "POST", 
    "https://www.reddit.com/api/v1/access_token", 
    strings.NewReader(form.Encode()),
  )

  if err != nil {
    return
  }

  request.Header.Set("User-agent", r.UserAgent)
  request.SetBasicAuth(r.AppId, r.AppSecret)

  data,err := r.JsonResponse(request)
  if err != nil {
    return
  }

  var obj map[string]interface{}
  err = json.Unmarshal(data, &obj)

  if token,ok := obj["access_token"]; ok {
    r.AccessToken = token.(string)
  } else {
    err = errors.New("Couldn't acquire an access token")
  }

  return
}


// Request creates a new http.Request with the default values used by every call
// to the Reddit API (like the User-agent and token). It should be used after 
// acquiring an access token. 
func (r *Reddit) Request(requestType string, call string, form url.Values) (*http.Request,error) {
  body := strings.NewReader(form.Encode())
  request, err := http.NewRequest(requestType, "https://oauth.reddit.com" + call, body)
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

  return request,nil
}

// JsonResponse takes an *http.Request and make a request to Reddit using the
// default client. Then the JSON body from http.Response is parsed and analyzed.
// If there are no errors in the JSON object (including API-related errors as
// "too many requests" blocks and such) a []byte representing the JSON data is
// returned. 
func (r *Reddit) JsonResponse(request *http.Request) (result []byte, err error) {
  response,err := r.Client.Do(request)
  if err != nil {
    return
  }

  result,err = ioutil.ReadAll(response.Body)
  if err != nil {
    return
  }

  var data map[string]interface{}

  err = json.Unmarshal(result, &data)
  if err != nil {
    return
  }

  // Catch the HTTP errors (404, 403, etc)
  if _,exists := data["error"]; exists {
    err = errors.New(fmt.Sprintf("%v %v",data["error"],data["message"]))
    return
  }

  // API errors (blocks due too many requests, etc)
  if v,exists := data["json"].(map[string]interface{}); exists {
    val := v["errors"].([]interface{})
    if len(val) > 0 {
      apiErr := val[0].([]interface{})
      err = errors.New(fmt.Sprintf("[%v] %v", apiErr[0], apiErr[1]))
      return
    }
  }

  return
}