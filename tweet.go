package main

import (
	"github.com/jchengj/twitter_service/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

type Tweet struct {
	gorm.Model
	Message   string
	MessageId int64
}
