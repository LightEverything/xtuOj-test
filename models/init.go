package models

import (
	"fmt"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strconv"
	"xtuOj/helper"
)

var DB = InitDB()
var RDB = InitRedis()

type dsnData struct {
	ip       string
	user     string
	password string
	userDB   string
	arg      string
}

func getDsn() *dsnData {
	dsnSt := new(dsnData)
	data := helper.GetDataFromConfigFile("MysqlConfig.txt")

	for k, v := range data {
		// 参数匹配
		switch k {
		case "ip":
			dsnSt.ip = v
		case "user":
			dsnSt.user = v
		case "password":
			dsnSt.password = v
		case "usedb":
			dsnSt.userDB = v
		case "arg":
			dsnSt.arg = v
		}
	}

	return dsnSt
}

func getDsnString(dsn *dsnData) string {
	debug := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?%s", dsn.user, dsn.password, dsn.ip, dsn.userDB, dsn.arg)
	return debug
}

func InitDB() *gorm.DB {
	dnsSt := getDsn()
	dnsString := getDsnString(dnsSt)

	tDb, err := gorm.Open(mysql.Open(dnsString), &gorm.Config{})
	if err != nil {
		log.Print("open db failure:", err)
	}
	return tDb
}

// 初始化redis
func InitRedis() *redis.Client {
	data := helper.GetDataFromConfigFile("redisConfig.txt")

	redisDB, err := strconv.Atoi(data["DB"])
	if err != nil {
		log.Println("connect redis error:", err)
		panic(err)
	}

	return redis.NewClient(&redis.Options{
		Addr:     data["Addr"] + ":" + data["Post"],
		Password: data["Password"],
		DB:       redisDB,
	})
}
