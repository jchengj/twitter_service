package main

import(
  "time"
  "fmt"
  "github.com/jinzhu/gorm"
)

type Twitter struct {
  gorm.Model
  LastCheckedAt time.Time
  SinceId int64
  ConsumerKey string
  ConsumerSecret string
  OauthToken string
  OauthTokenSecret string
}

func (twitter *Twitter) poller(){
  api := connection(twitter)
  if result, err := api.GetUserTimeline(nil); err != nil{
    panic(err)
  } else {
    if len(result) > 0 {
      
      twitter.SinceId       = result[0].Id
      twitter.LastCheckedAt = time.Now()

      for _, tweet := range result{ 
        fmt.Printf("Inserting [%d] message: %s", tweet.Id, tweet.Text)
        db.Create(&Tweet{Message: tweet.Text, MessageId: tweet.Id})
      }

      db.Save(twitter)
    }
  }
}

func Account(id int64) *Twitter{
  var account Twitter
  db.Find(&account, id)

  return &account
}

