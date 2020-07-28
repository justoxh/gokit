package logger

import (
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

func TestGetLoggerWithOptions(t *testing.T) {
	options := &Options{
		Formatter:      "json",
		Write:          true,
		Path:           "../logs/",
		DisableConsole: false,
		WithCallerHook: true,
		MaxAge:         time.Duration(7*24) * time.Hour,
		RotationTime:   time.Duration(7) * time.Hour,
	}

	log := GetLoggerWithOptions("default", options)
	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	})
	log.Info("Hello world")
}

type Item struct {
	SumNum int64 `json:"sum_num"`
}

func TestGetGormLoggerWithOptions(t *testing.T) {
	options := &Options{
		Formatter:      "json",
		Write:          false,
		Path:           "./logs/",
		DisableConsole: false,
		WithCallerHook: true,
		MaxAge:         time.Duration(7*24) * time.Hour,
		RotationTime:   time.Duration(7) * time.Hour,
	}

	log := GetGormLoggerWithOptions(options)
	var err error

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		"root", "Cocos!@#$5678", "localhost", 3306, "learn")
	db, err := gorm.Open("mysql", dataSourceName)
	if err != nil {
		t.Fatal(err)
	}

	db.SetLogger(log)
	db.LogMode(true)
	db.Exec("USE learn")
	// Migration to create tables for Order and Item schema
	//db.AutoMigrate(&Item{})
	var item Item
	sql := "select name as sum_num from items where id = 1"
	db.Raw(sql).Scan(&item)
	t.Log(10 - item.SumNum)
}
