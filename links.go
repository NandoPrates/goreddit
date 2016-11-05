package goreddit

import (
  "net/url"
)

// SubmitText submits a text (markdown format) to a subreddit
func (r *Reddit) SubmitText(sub string, title string, text string) error {
  return r.Submit(sub, title, text, "self")
}

// SubmitLink submits a new link to the specified subreddit
func (r *Reddit) SubmitLink(sub string, title string, url string) error {
  return r.Submit(sub, title, url, "link")
}

// Submit is the general purpose function to submit a new thread to Reddit
// it can be either a link or a text (self) thread. This shouldn't be called directly
// use SubmitLink or SubmitText instead.
func (r *Reddit) Submit(sub string, title string, content string, kind string) (err error) {
  form := url.Values{
    "api_type": {"json"},
    "sr": {sub},
    "kind": {kind},
    "title": {title},
    "resubmit": {"true"},
  }

  if kind == "self" {
    form.Add("text", content)
  } else {
    form.Add("url", content)
  }

  request,err := r.Request("POST", "/api/submit", form)
  if err != nil {
    return
  }

  _,err = r.JsonResponse(request)
  return
}