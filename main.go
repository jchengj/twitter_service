package main

import (
  "fmt"
  "time"
  "sync"
  "github.com/ChimeraCoder/anaconda"
  "github.com/jinzhu/gorm"
  _ "github.com/go-sql-driver/mysql"  
)

var db gorm.DB
var wg sync.WaitGroup

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
  for {
      accounts := make([]Twitter, 1)
      db.Where("last_checked_at < ?", time.Now().Add(-time.Hour)).Find(&accounts)
      if len(accounts) > 0{
        wg.Add(len(accounts))
        fmt.Println("Started Processing")
        for _, account := range accounts{
          go func(){
            account.poll()
            wg.Done()
          }() 
        }
        fmt.Println("Waiting for all goroutines")
        wg.Wait()
        fmt.Println("Finished Processing")
      }
  }
}
