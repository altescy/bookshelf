package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/altescy/bookshelf/browser"
	"github.com/altescy/bookshelf/controller"
	"github.com/altescy/bookshelf/model"
	"github.com/altescy/bookshelf/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
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
	// dsn := fmt.Sprintf(`postgres://%s@%s:%s/%s?sslmode=disable`, dbusrpass, dbhost, dbport, dbname)
	// db, err := gorm.Open("postgres", dsn)
	db, err := gorm.Open("sqlite3", "/tmp/bookshelf.db")

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
		createBucket    = getEnv("AWS_S3_CREATE_BUCKET", "")
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

	// create bucket if not exists
	if createBucket != "" {
		log.Printf("[INFO] check bucket existence")
		_, err := svc.HeadBucket(&s3.HeadBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				switch awsErr.Code() {
				case "NotFound":
					log.Printf("[INFO] create a new bucket")
					_, err = svc.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucket)})
					if err != nil {
						log.Fatalf("failed to create bucket: %v", err)
					}
				default:
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
	}

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
	defer db.Close()

	autoMigrate(db)

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

	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:     browser.Asset,
		AssetDir:  browser.AssetDir,
		AssetInfo: browser.AssetInfo,
		Prefix:    "/dist",
		Fallback:  "index.html",
	})
	h := controller.NewHandler(db, storage, isEnableCors)

	router := httprouter.New()
	router.POST("/api/book", h.AddBook)
	router.GET("/api/book/:bookid", h.GetBook)
	router.PUT("/api/book/:bookid", h.UpdateBook)
	router.DELETE("/api/book/:bookid", h.DeleteBook)
	router.GET("/api/book/:bookid/file/:ext", h.DownloadFile)
	router.DELETE("/api/book/:bookid/file/:ext", h.DeleteFile)
	router.POST("/api/book/:bookid/files", h.UploadFiles)
	router.GET("/api/books", h.GetBooks)
	router.GET("/api/mime/:ext", h.GetMime)
	router.GET("/api/mimes", h.GetMimes)
	router.GET("/api/opds", h.GetOPDSFeed)
	router.NotFound = fileServer

	addr := ":" + port
	log.Printf("[INFO] start server %s", addr)
	log.Fatal(http.ListenAndServe(addr, h.CommonMiddleware(router)))
}
