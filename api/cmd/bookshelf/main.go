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
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

const (
	SessionSecret = "session_secret"
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

func createSessionStore() sessions.Store {
	store := sessions.NewCookieStore([]byte(SessionSecret))
	return store
}

func autoMigrate(db *gorm.DB) {
	if err := model.AutoMigrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}

func createUser(db *gorm.DB) {
	var (
		username = getEnv("USERNAME", "user")
		password = getEnv("PASSWORD", "password")
	)
	if err := model.UserSignUp(db, username, password); err != nil {
		if err != model.ErrUserConflict {
			log.Fatalf("failed to sign up user: %v", err)
		}
	}
}

func main() {
	port := getEnv("PORT", "8080")

	db := createGormDB()
	defer db.Close()

	autoMigrate(db)
	createUser(db)

	store := createSessionStore()

	h := controller.NewHandler(db, store)

	router := httprouter.New()
	router.GET("/", h.Index)
	router.POST("/signin", h.Signin)
	router.POST("/signout", h.Signout)
	router.GET("/user", h.GetUser)

	addr := ":" + port
	log.Printf("[INFO] start server %s", addr)
	log.Fatal(http.ListenAndServe(addr, gctx.ClearHandler(h.CommonMiddleware(router))))
}
