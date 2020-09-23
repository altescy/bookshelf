package cmd

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

// EnvPrefix is a prefix for environment variables.
const EnvPrefix = "BOOKSHELF_"

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
	if v, ok := os.LookupEnv(EnvPrefix + key); ok {
		return v
	}
	return def
}

func createGormDB() *gorm.DB {
	var (
		dbURL = getEnv("DB_URL", "")
	)

	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		log.Fatalf("cannot parse database url. err: %v", err)
	}

	var db *gorm.DB

	switch parsedURL.Scheme {
	case "sqlite3":
		db, err = gorm.Open("sqlite3", parsedURL.Path)
		if err != nil {
			log.Fatalf("sqlite3 connect failed. err: %s", err)
		}
	case "postgres":
		db, err = gorm.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("postgres connect failed. err: %s", err)
		}
	default:
		log.Fatal("invalid db")
	}

	return db
}

func createStorage() storage.Storage {
	var (
		storageURL       = getEnv("STORAGE_URL", "")
		createNewStorage = getEnv("CREATE_NEW_STORAGE", "")
	)

	var store storage.Storage

	parsedURL, err := url.Parse(storageURL)
	if err != nil {
		log.Fatalf("cannot parse storage url: %v", err)
	}

	doCreateNewStorage := createNewStorage != ""

	switch parsedURL.Scheme {
	case "s3":
		bucket := parsedURL.Host
		root := parsedURL.Path
		store = createS3Storage(bucket, root, doCreateNewStorage)
	case "file":
		root := parsedURL.Path
		store = createFileSysteStorage(root, doCreateNewStorage)
	default:
		log.Fatalf("invalid storage url scheme: %s", parsedURL.Scheme)
	}

	return store
}

func createFileSysteStorage(root string, createNewStorage bool) *storage.FileSystemStorage {
	perm := os.FileMode(0777)

	if err := os.MkdirAll(root, perm); err != nil {
		log.Fatalf("faield to create storage dir. err: %v", err)
	}

	return storage.NewFileSystemStorage(root, perm)
}

func createS3Storage(bucket, root string, createNewStorage bool) *storage.S3Storage {
	var (
		awsAccessKey    = getEnv("AWS_ACCESS_KEY_ID", "")
		awsSecretKey    = getEnv("AWS_SECRET_ACCESS_KEY", "")
		awsSessionToken = getEnv("AWS_SESSION_TOKEN", "")
		s3Region        = getEnv("AWS_S3_REGION", "")
		s3EndpointURL   = getEnv("AWS_S3_ENDPOINT_URL", "")
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
	if createNewStorage {
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

func Main() {
	var (
		port       = getEnv("PORT", "8080")
		enableCors = getEnv("ENABLE_CORS", "")
	)

	isEnableCors := enableCors != ""
	log.Printf("[INFO] enable CORS: %v", isEnableCors)

	db := createGormDB()
	defer db.Close()

	autoMigrate(db)

	store := createStorage()

	h := controller.NewHandler(db, store, isEnableCors)

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
	router.GET("/opds", h.GetOPDSFeed)
	router.NotFound = http.FileServer(&assetfs.AssetFS{
		Asset:     browser.Asset,
		AssetDir:  browser.AssetDir,
		AssetInfo: browser.AssetInfo,
		Prefix:    "/dist",
		Fallback:  "index.html",
	})

	addr := ":" + port
	log.Printf("[INFO] start server %s", addr)
	log.Fatal(http.ListenAndServe(addr, h.CommonMiddleware(router)))
}
