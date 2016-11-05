# GoReddit

GoReddit is a wrapper for the Reddit API written in Go. 

:warning: This project still in the early stages of development. 

# Usage

```go
package main

import(
  "fmt"
  "time"
  "github.com/jefferson-dab/goreddit"
)

func main() {
  // Connect to Reddit using Oauth
  client,err := goreddit.New(
    "my_username",
    "my_password",
    "my_app_id",
    "my_app_secret",
    "my_user_agent",
  )

  if err != nil {
    // Handle the error...
  }

  // Get the 10 latest comments from /r/programming
  comments,err := client.ListComments("programming", "new", 10)

  if err != nil {
    // Handle the error
  }

  for _,v := range comments {
    // Send a friendly reply to each one of the comments we got
    err = client.Comment(
      v.Fullname, 
      fmt.Sprintf("Btw, everybody knows that **spaces** are better than **tabs**, %s.", v.Author),
    )

    if err != nil {
      // Handle errors yet again (you never need it until you do)
    }

    time.Sleep(2 * time.Second)
  }

  // Get informations about myself
  me,_ := client.Me()

  if me.CommentKarma > 9000 {
    // Submit a new text post to a subreddit (fictional /r/TheSubOfPeopleWhoAdoreMe)
    _ = client.SubmitText(
      "TheSubOfPeopleWhoAdoreMe",
      "It's over 9000!",
      "I'd like to let you guys know that my Karma is getting stronger!"
    )
  }
}
```
