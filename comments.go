package goreddit

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	Upvote   = 1
	Downvote = -1
	Unvote   = 0
)

type Comment struct {
	Author      string      `json:"author"`
	Text        string      `json:"body"`
	Title       string      `json:"title"`
	HtmlText    string      `json:"body_html"`
	Fullname    string      `json:"name"`
	Downvotes   int         `json:"downs"`
	Upvotes     int         `json:"ups"`
	Score       int         `json:"score"`
	ScoreHidden bool        `json:"score_hidden"`
	Gilded      int         `json:"gilded"`
	Parent      string      `json:"parent_id"`
	Edited      interface{} `json:"edited"` // Unknown type (bug?) It can be a float or a bool
	Archived    bool        `json:"archived"`
	Created     float64     `json:"created"`
	Replies     Reply       `json:"replies"`
}

type Reply struct {
	Data struct {
		Comments []struct {
			Comment Comment `json:"data"`
		} `json:"children"`
	}
}

// Comment send the text message as a reply to the parent message. The parent string
// must be a Reddit's fullname identifier
func (r *Reddit) Comment(parent string, text string) (err error) {
	form := url.Values{
		"api_type": {"json"},
		"text":     {text},
		"thing_id": {parent},
	}
	request, err := r.Request("POST", "/api/comment", form)
	if err != nil {
		return
	}

	_, err = r.JsonResponse(request)
	return
}

// ListCommentsSub lists the comments from all the links on a sub
func (r *Reddit) ListCommentsSub(sub string, options ListingOpt) (comments Reply, err error) {
	return r.ListComments(sub, "", options)
}

// ListComments lists the comments from a particular link of the sub
func (r *Reddit) ListComments(sub string, idlink string, options ListingOpt) (comments Reply, err error) {
	form := options.Values()
	request, err := r.Request("GET", "/r/"+sub+"/comments/"+idlink, form)
	if err != nil {
		return
	}

	data, err := r.JsonResponse(request)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &comments)
	if err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			// Silently ignore a UnmarshalTypeError :)
			err = nil
		} else {
			return
		}
	}

	return
}

// Vote takes a fullname (it can be a link or a comment) and cast a vote on it.
// Votes can be goreddit.Upvote, goreddit.Downvote and goreddit.Unvote.
// NOTE from Reddit: votes must be cast by humans. That is, API clients proxying a
// human's action one-for-one are OK, but bots deciding how to vote on content or
// amplifying a human's vote are not. See the reddit rules for more details on what
// constitutes vote cheating.
func (r *Reddit) Vote(fullname string, vote int) (err error) {
	form := url.Values{
		"id":  {fullname},
		"dir": {fmt.Sprintf("%d", vote)},
	}
	request, err := r.Request("POST", "/api/vote", form)
	if err != nil {
		return
	}
	_, err = r.JsonResponse(request)
	return
}
