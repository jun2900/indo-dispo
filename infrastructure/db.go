package infrastructure

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/joho/godotenv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	//_ "github.com/jinzhu/gorm/dialects/postgres"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"

	"os"
)

// Database structure
type Database struct {
	*gorm.DB
}

var (
	_  = godotenv.Load(".env")
	DB *gorm.DB
)

// OpenDbConnection Opening a database and save the reference to `Database` struct.
func OpenDbConnection() *gorm.DB {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	username := os.Getenv("DB_USER")     //utility.KVGet("DB_USER")
	password := os.Getenv("DB_PASSWORD") //utility.KVGet("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")       //utility.KVGet("DB_NAME")
	host := os.Getenv("DB_HOST")         //utility.KVGet("DB_HOST")
	port := os.Getenv("DB_PORT")         //utility.KVGet("DB_PORT")
	var db *gorm.DB
	var err error

	databaseURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)

	db, err = gorm.Open(mysql.Open(databaseURL), &gorm.Config{Logger: newLogger})

	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
		os.Exit(-1)
	}

	//db.DB().SetMaxIdleConns(10)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	// db.LogMode(true)
	DB = db

	return DB
}

// Delete the database after running testing cases.
func RemoveDb(db *gorm.DB) error {
	sqlDB, _ := db.DB()
	//db.Close()
	sqlDB.Close()
	err := os.Remove(path.Join(".", "app.db"))
	return err
}

// Using this function to get a connection, you can create your connection pool here.
func GetDb() *gorm.DB {
	return DB
}
