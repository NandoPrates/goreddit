package goreddit

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Link struct {
	Fullname    string  `json:"name"`
	Domain      string  `json:"domain"`
	Subreddit   string  `json:"subreddit"`
	Title       string  `json:"title"`
	Text        string  `json:"selftext_html"`
	Author      string  `json:"author"`
	Id          string  `json:"id"`
	NumComments int     `json:"num_comments"`
	Gilded      int     `json:"gilded"`
	Clicked     bool    `json:"clicked"`
	Score       int     `json:"score"`
	Upvotes     int     `json:"ups"`
	Downvotes   int     `json:"downs"`
	HiddenScore bool    `json:"hide_score"`
	Thumbnail   string  `json:"thumbnail"`
	Archived    bool    `json:"archived"`
	IsSelf      bool    `json:"is_self"`
	Spoiler     bool    `json:"spoiler"`
	Locked      bool    `json:"locked"`
	Stickied    bool    `json:"stickied"`
	Created     float64 `json:"created"`
	LinkFlair   string  `json:"link_flair_text"`
	Permalink   string  `json:"permalink"`
	Url         string  `json:"url"`
}

type LinkList struct {
	Links  []*Link
	Before string
	After  string
}

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
		"sr":       {sub},
		"kind":     {kind},
		"title":    {title},
		"resubmit": {"true"},
	}

	if kind == "self" {
		form.Add("text", content)
	} else {
		form.Add("url", content)
	}

	request, err := r.Request("POST", "/api/submit", form)
	if err != nil {
		return
	}

	_, err = r.JsonResponse(request)
	return
}

// ListLinks gets <limit> number of links from a <sub> according to the <sort> specified.
// Sort can be either "new", "hot", "top" or "controversial".
func (r *Reddit) ListLinks(sub string, sort string, options ListingOpt) (list *LinkList, err error) {
	form := options.Values()
	request, err := r.Request("GET", fmt.Sprintf("/r/%s/%s", sub, sort), form)
	if err != nil {
		return
	}

	data, err := r.JsonResponse(request)
	if err != nil {
		return
	}

	var obj struct {
		Data struct {
			Links []struct {
				Link *Link `json:"data"`
			} `json:"children"`
			Before string
			After  string
		}
	}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return
	}

	var tempList []*Link
	for _, v := range obj.Data.Links {
		tempList = append(tempList, v.Link)
	}

	list = &LinkList{
		tempList,
		obj.Data.Before,
		obj.Data.After,
	}

	return
}

// Hide receives a list of link fullnames, hiding each one of them
func (r *Reddit) Hide(fullname string) (err error) {
	return r.ActionOnThing(fullname, "/api/hide")
}

// Unhide receives a list of link fullnames, unhiding each one of them
func (r *Reddit) Unhide(fullname string) (err error) {
	return r.ActionOnThing(fullname, "/api/unhide")
}
