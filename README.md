Bookshelf
=========

[![Actions Status](https://github.com/altescy/bookshelf/workflows/build/badge.svg)](https://github.com/altescy/bookshelf/actions?query=workflow%3Abuild)
[![License](https://img.shields.io/github/license/altescy/bookshelf)](https://github.com/altescy/bookshelf/blob/master/LICENSE)
[![Release](https://img.shields.io/github/v/release/altescy/bookshelf)](https://github.com/altescy/bookshelf/releases)

[Blog post (Japanese)](https://altescy.jp/posts/works/bookshelf/)

`Bookshelf` is a simple ebook management web application.
You can easily store and manage your books on a local or S3 compatible storage.
This software also provides a OPDS feed which enables you to read your books via any OPDS readers on your computer or smartphone.

![Screenshot_2020-09-22 bookshelf](https://user-images.githubusercontent.com/16734471/93875665-5c6a5d00-fd10-11ea-81df-3a1735aa4547.png)


### Usage

```
$ go get github.com/altescy/bookshelf
$ export BOOKSHELF_DB_URL=sqlite3:///`pwd`/data/bookshelf.db
$ export BOOKSHELF_STORAGE_URL=file:///`pwd`/data/files
$ bookshelf
```

### Docker

```
$ docker pull altescy/bookshelf
$ docker run -d \
    -v `pwd`/data:/data \
    -p 8080:8080 \
    -e BOOKSHELF_DB_URL=sqlite3:///data/bookshelf.db \
    -e BOOKSHELF_STORAGE_URL=file:///data/files \
    altescy/bookshelf
```


### docker-compose

```
$ git clone https://github.com/altescy/bookshelf.git
$ cd bookshelf
$ cat << EOF > .env
BOOKSHELF_PORT=80
BOOKSHELF_ENABLE_CORS=
BOOKSHELF_DB_URL=postgres://user:password@postgres:5432/bookshelf?sslmode=disable
BOOKSHELF_STORAGE_URL=s3://books
BOOKSHELF_CREATE_NEW_STORAGE=1
BOOKSHELF_AWS_ACCESS_KEY_ID=minio_access
BOOKSHELF_AWS_SECRET_ACCESS_KEY=minio_secret
BOOKSHELF_AWS_S3_REGION=us-east-1
BOOKSHELF_AWS_S3_ENDPOINT_URL=http://minio

MINIO_ACCESS_KEY=minio_access
MINIO_SECRET_KEY=minio_secret
MINIO_HOST=0.0.0.0
MINIO_PORT=9000

POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_PORT=5432

TZ=Asia/Tokyo
EOF
$ docker-compose up -d
```
