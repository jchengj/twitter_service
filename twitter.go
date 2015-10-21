package main

import (
	"fmt"
	"github.com/jchengj/twitter_service/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"net/url"
	"time"
)

type Twitter struct {
	gorm.Model
	LastCheckedAt    time.Time
	SinceId          int64
	ConsumerKey      string
	ConsumerSecret   string
	OauthToken       string
	OauthTokenSecret string
}

func (twitter *Twitter) poll() {
	api := connection(twitter)
	opts := url.Values{}

	if twitter.SinceId != 0 {
		opts.Set("since_id", fmt.Sprintf("%d", twitter.SinceId))
	}

	if result, err := api.GetUserTimeline(opts); err != nil {
		panic(err)
	} else {
		if len(result) > 0 {

			twitter.SinceId = result[0].Id
			twitter.LastCheckedAt = time.Now()

			for _, tweet := range result {
				Info.Printf("Inserting [%d] message: %s\n", tweet.Id, tweet.Text)
				db.Create(&Tweet{Message: tweet.Text, MessageId: tweet.Id})
			}

			db.Save(twitter)
		}
	}
}

func (twitter *Twitter) Send(values *[]string) {
	api := connection(twitter)
	for _, v := range *values {
		Info.Printf("Posting message [%s]", v)
		api.PostTweet(v, nil)
	}
}
