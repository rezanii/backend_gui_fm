package connectionSetup

import (
	"backend_gui/utils"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)
/*
fuction to connect redis
by Irma 30/12/2021
*/
func RedisConnection() *redis.Client {
	errEnv := godotenv.Load("config/" + utils.Param.EnvName)
	if errEnv != nil {
		log.Error("Error connection redis :", errEnv.Error())
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
	return rdb
}

func SetupDbConnection() *gorm.DB {
	utils.InitHandlers()
	errEnv := godotenv.Load("config/" + utils.Param.EnvName)
	if errEnv != nil {
		log.Error("Error conection db : ", errEnv.Error())
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error("Error conection db : ", err.Error())
	}
	return db
}

func CloseDbConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		log.Error("Error close conection db : ", err.Error())
	}
	dbSQL.Close()
}

