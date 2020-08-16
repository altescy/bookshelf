package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/altescy/bookshelf/api/controller"
	"github.com/altescy/bookshelf/api/model"
	"github.com/altescy/bookshelf/api/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

func createS3Storage(bucket, root string) *storage.S3Storage {
	var (
		awsAccessKey    = getEnv("AWS_ACCESS_KEY_ID", "m1n10_4cce55")
		awsSecretKey    = getEnv("AWS_SECRET_ACCESS_KEY", "m1n10_5ecret")
		awsSessionToken = getEnv("AWS_SESSION_TOKEN", "")
		s3Region        = getEnv("AWS_S3_REGION", "us-east-1")
		s3EndpointURL   = getEnv("AWS_S3_ENDPOINT_URL", "http://localhost:9000")
	)

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, awsSessionToken),
		Endpoint:         aws.String(s3EndpointURL),
		Region:           aws.String(s3Region),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("cannot create aws session: %v", err)
	}
	svc := s3.New(sess)

	return storage.NewS3Storage(svc, bucket, root)
}

func autoMigrate(db *gorm.DB) {
	if err := model.AutoMigrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}

func main() {
	var (
		port       = getEnv("PORT", "8080")
		enableCors = getEnv("ENABLE_CORS", "")
		storageURL = getEnv("STORAGE_URL", "s3://books")
	)

	db := createGormDB()
	autoMigrate(db)
	defer db.Close()

	var storage storage.Storage
	surl, err := url.Parse(storageURL)
	if err != nil {
		log.Fatalf("cannot parse storage url: %v", err)
	}
	switch surl.Scheme {
	case "s3":
		bucket := surl.Host
		root := surl.Path
		storage = createS3Storage(bucket, root)
	default:
		log.Fatalf("invalid scheme: %s", surl.Scheme)
	}

	isEnableCors := enableCors != ""
	log.Printf("[INFO] enable CORS: %v", isEnableCors)

	h := controller.NewHandler(db, storage, isEnableCors)

	router := httprouter.New()
	router.POST("/book", h.AddBook)
	router.GET("/book/:bookid", h.GetBook)
	router.PUT("/book/:bookid", h.UpdateBook)
	router.DELETE("/book/:bookid", h.DeleteBook)
	router.POST("/book/:bookid/files", h.UploadFiles)
	router.GET("/book/:bookid/file/:ext", h.DownloadFile)
	router.GET("/books", h.GetBooks)
	router.GET("/mime/:ext", h.GetMime)
	router.GET("/mimes", h.GetMimes)

	addr := ":" + port
	log.Printf("[INFO] start server %s", addr)
	log.Fatal(http.ListenAndServe(addr, h.CommonMiddleware(router)))
}
