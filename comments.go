package goreddit

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Comment struct {
	Author      string  `json:"author"`
	Text        string  `json:"body"`
	Title       string  `json:"title"`
	HtmlText    string  `json:"body_html"`
	Fullname    string  `json:"name"`
	Downvotes   int     `json:"downs"`
	Upvotes     int     `json:"ups"`
	Score       int     `json:"score"`
	ScoreHidden bool    `json:"score_hidden"`
	Gilded      int     `json:"gilded"`
	Parent      string  `json:"parent_id"`
	Edited      bool    `json:"edited"`
	Archived    bool    `json:"archived"`
	Created     float64 `json:"created"`
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

func (r *Reddit) ListComments(sub string, order string, limit int) (comments []Comment, err error) {
	form := url.Values{
		"sort":  {order},
		"limit": {fmt.Sprintf("%d", limit)},
	}

	request, err := r.Request("GET", "/r/"+sub+"/comments", form)
	if err != nil {
		return
	}

	data, err := r.JsonResponse(request)
	if err != nil {
		return
	}

	var obj struct {
		Data struct {
			Children []struct {
				Data Comment
			}
		}
	}
	_ = json.Unmarshal(data, &obj)

	for k := range obj.Data.Children {
		comments = append(comments, obj.Data.Children[k].Data)
	}

	return
}
