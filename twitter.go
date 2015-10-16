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

func (t *Twitter) poller(){
  api := connection(t)
  if result, err := api.GetUserTimeline(nil); err != nil{
    panic(err)
  } else {
    for _, t := range result{ 
      fmt.Printf("Inserting [%d] message: %s", t.Id, t.Text)
      db.Create(&Tweet{Message: t.Text, MessageId: t.Id})
    }
  }
}

func Account(id int64) *Twitter{
  var account Twitter
  db.Find(&account)

  return &account
}

