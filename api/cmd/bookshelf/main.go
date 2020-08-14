package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/altescy/bookshelf/api/controller"
	"github.com/altescy/bookshelf/api/model"
	gctx "github.com/gorilla/context"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func init() {
	var err error

	tz := getEnv("TZ", "Asia/Tokyo")

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Panicln(err)
	}
	time.Local = loc
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv("API_" + key); ok {
		return v
	}
	return def
}

func createGormDB() *gorm.DB {
	var (
		dbhost = getEnv("DB_HOST", "127.0.0.1")
		dbport = getEnv("DB_PORT", "5432")
		dbuser = getEnv("DB_USER", "user")
		dbpass = getEnv("DB_PASSWORD", "password")
		dbname = getEnv("DB_NAME", "bookshelf")
	)

	dbusrpass := dbuser
	if dbpass != "" {
		dbusrpass += ":" + dbpass
	}

	dsn := fmt.Sprintf(`postgres://%s@%s:%s/%s?sslmode=disable`, dbusrpass, dbhost, dbport, dbname)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("postgres connect failed. err: %s", err)
	}

	return db
}

func autoMigrate(db *gorm.DB) {
	if err := model.AutoMigrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}

func main() {
	port := getEnv("PORT", "8080")

	db := createGormDB()
	defer db.Close()

	autoMigrate(db)

	h := controller.NewHandler(db)

	router := httprouter.New()
	router.GET("/", h.Index)

	addr := ":" + port
	log.Printf("[INFO] start server %s", addr)
	log.Fatal(http.ListenAndServe(addr, gctx.ClearHandler(h.CommonMiddleware(router))))
}
