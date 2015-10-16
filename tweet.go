package main

import(
  "github.com/jinzhu/gorm"
)

type Tweet struct{
  gorm.Model
  Message     string
  MessageId   int64
}