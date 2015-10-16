package main

import (
  "fmt"
  "github.com/ChimeraCoder/anaconda"
  "github.com/jinzhu/gorm"
  _ "github.com/go-sql-driver/mysql"  
)

var db gorm.DB

func init() {
  adapter     := "mysql"
  credentials := "root@tcp(localhost:3306)/golang?charset=utf8&parseTime=True&loc=Local"

  if conn, err := gorm.Open(adapter, credentials); err == nil {
    db = conn
    db.DB()
  } else {
    panic("Failed to establish database connection")  
  }
}

func connection(t *Twitter) *anaconda.TwitterApi{
  anaconda.SetConsumerKey(t.ConsumerKey)
  anaconda.SetConsumerSecret(t.ConsumerSecret)
  return anaconda.NewTwitterApi(t.OauthToken, t.OauthTokenSecret)
}

func main(){
  fmt.Println("Starting Twitter Service")

  account := Account(1)
  //api.PostTweet("From GO Service", nil)
  account.poller() 
}
