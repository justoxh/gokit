package logger

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"testing"
	"time"

)

func TestGetLoggerWithOptions(t *testing.T) {
	options := &Options{
		Formatter:      "json",
		Write:          false,
		Path:           "./logs/",
		DisableConsole: false,
		WithCallerHook: true,
		MaxAge:         time.Duration(7*24) * time.Hour,
		RotationTime:   time.Duration(7) * time.Hour,
	}

	log := GetLoggerWithOptions("default", options)
	log.Info("Hello world")
}

type Item struct {
	gorm.Model
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
		"root","","localhost",3306,"learn")
	db, err := gorm.Open("mysql", dataSourceName)
	if err != nil {
		t.Fatal(err)
	}

	db.SetLogger(log)
	db.LogMode(true)
	db.Exec("USE learn")
	// Migration to create tables for Order and Item schema
	db.AutoMigrate(&Item{})
	var items []Item
	db.Find(&items)
}