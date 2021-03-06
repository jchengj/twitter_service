package main

import (
	"github.com/jchengj/twitter_service/Godeps/_workspace/src/github.com/ChimeraCoder/anaconda"
	_ "github.com/jchengj/twitter_service/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/jchengj/twitter_service/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/jchengj/twitter_service/Godeps/_workspace/src/gopkg.in/redis.v3"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var mainWG sync.WaitGroup
var db gorm.DB
var rClient *redis.Client

var (
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func init() {

	// ---------------- MySQL DB --------------- //
	adapter := "mysql"
	credentials := "root@tcp(localhost:3306)/golang?charset=utf8&parseTime=True&loc=Local"

	if conn, err := gorm.Open(adapter, credentials); err == nil {
		db = conn
		db.DB()
	} else {
		panic("Failed to establish database connection")
	}

	// ---------------- REDIS ------------- //

	rClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       5,  // let's not messed up your assistly redis data
	})

	// ----------------- LOG -------------- //
	Debug = log.New(os.Stdout, "Debug: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(os.Stdout, "Warn: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func connection(t *Twitter) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(t.ConsumerKey)
	anaconda.SetConsumerSecret(t.ConsumerSecret)
	return anaconda.NewTwitterApi(t.OauthToken, t.OauthTokenSecret)
}

// --- pull messages from redis and send to twitter --- //
// --- assuming the key looks like this - tweet:13213:1 //
func send() {
	Debug.Println("Send routine started")
	var wg sync.WaitGroup

	if tweets, err := rClient.Keys("*").Result(); err == nil {
		accounts := make(map[string][]string)

		for _, tweet := range tweets {
			s := strings.Split(tweet, ":")
			accounts[s[2]] = append(accounts[s[2]], rClient.Get(tweet).Val())
		}

		wg.Add(len(accounts))
		for key, val := range accounts {
			go func() {
				var account Twitter
				db.Find(&account, key)
				account.Send(&val)
				wg.Done()
			}()
		}
		wg.Wait()
		rClient.FlushDb()

	}

	Debug.Println("Receive routine ended")
	mainWG.Done()
}

// --- pull tweets from twitter and save to db --- //
func receive() {
	Debug.Println("Receive routine started")
	var wg sync.WaitGroup

	accounts := make([]Twitter, 1)
	db.Where("last_checked_at < ?", time.Now().Add(-time.Hour)).Find(&accounts)
	if len(accounts) > 0 {
		wg.Add(len(accounts))
		for _, account := range accounts {
			go func() {
				account.poll()
				wg.Done()
			}()
		}
		wg.Wait()
	}

	Debug.Println("Receive routine ended")
	mainWG.Done()
}

func main() {
	Info.Println("Starting Twitter Microservice")

	for {
		mainWG.Add(2)

		go send()
		go receive()

		mainWG.Wait()
	}
}
