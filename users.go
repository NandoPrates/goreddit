package goreddit

import(
  // "fmt"
  // "io/ioutil"
  "encoding/json"
  // "net/url"
)

type User struct {
  Username string `json:"name"`
  IsFriend bool `json:"is_friend"`
  Created float64 `json:"created"`
  HiddenFromRobots bool `json:"hidden_from_robots"`
  LinkKarma int `json:"link_karma"`
  CommentKarma int `json:"comment_karma"`
  Gold bool `json:"is_gold"`
  Mod bool `json:"is_mod"`
  Verified bool `json:"has_verified_email"`
  Id string `json:"id"`
}

type Me struct {
  User
  Over18 bool `json:"over_18"`
  Suspended bool `json:"is_suspended"`
  Employee bool `json:"is_employee"`
}

// User returns a *User with the available information to an username
func (r *Reddit) User(username string) (user *User, err error) {
  request,err := r.Request("GET", "/user/"+username+"/about", nil)
  if err != nil {
    return
  }
  
  data,err := r.JsonResponse(request)
  if err != nil {
    return
  }

  var obj struct { Data *User }
  err = json.Unmarshal(data, &obj)
  if err != nil {
    return
  }

  user = obj.Data
  return
}

// Me returns a *Me with informations about the current logged in account
func (r *Reddit) Me() (me *Me, err error) {
  request,err := r.Request("GET","/api/v1/me", nil)
  if err != nil {
    return
  }

  data,err := r.JsonResponse(request)
  if err != nil {
    return
  }

  err = json.Unmarshal(data, &me)
  return
}