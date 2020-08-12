package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/altescy/bookshelf/api/controller"
	gctx "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func init() {
	var err error

	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "Asia/Tokyo"
	}

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

func main() {
	var (
		port   = getEnv("PORT", "8080")
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

	dsn := fmt.Sprintf(`%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&charset=utf8mb4`, dbusrpass, dbhost, dbport, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("mysql connect failed. err: %s", err)
	}
	defer db.Close()

	h := controller.NewHandler(db)

	router := httprouter.New()
	router.GET("/", h.Index)

	addr := ":" + port
	log.Printf("[INFO] start server %s", addr)
	log.Fatal(http.ListenAndServe(addr, gctx.ClearHandler(h.CommonMiddleware(router))))
}
