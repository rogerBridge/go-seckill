package shop_orm

import (
	"go-seckill/internal/db"
	"go-seckill/internal/logconf"
	"log"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 初始化数据库 seckill
func InitialMysql() error {
	conn := db.Conn2

	conn.Exec("CREATE DATABASE IF NOT EXISTS seckill")
	log.Println("executed create database seckill command")

	err := conn.AutoMigrate(&Good{})
	if err != nil {
		log.Println(err)
		return err
	}

	err = conn.AutoMigrate(&PurchaseLimit{})
	if err != nil {
		log.Println(err)
		return err
	}

	err = conn.AutoMigrate(&User{})
	if err != nil {
		log.Println(err)
		return err
	}

	err = conn.AutoMigrate(&Order{})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

type SelfDefine struct {
	gorm.Model
	Version string `gorm:"default:v0.0.0"`
}

var conn = db.Conn2

var logger = logconf.BaseLogger.WithFields(logrus.Fields{"component": "shop_orm"})
