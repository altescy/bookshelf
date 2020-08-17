Bookshelf
=========

eBook management system


1. Clone this repository

```
$ git clone https://gltihub.com/altescy/bookshelf.git
```

2. Build frontent

```
$ cd bookshelf/frontend
$ yarn build
```

3. Setup environment variables

```
$ cd ..
$ cat << EOF > .env
TZ=Asia/Tokyo

API_PORT=80
API_ENDPOINT=/api
API_ENABLE_CORS=1
API_DB_HOST=postgres
API_DB_PORT=5432
API_DB_USER=user
API_DB_PASSWORD=password
API_DB_NAME=bookshelf
API_AWS_ACCESS_KEY_ID=minio_access
API_AWS_SECRET_ACCESS_KEY=minio_secret
API_AWS_S3_REGION=us-east-1
API_AWS_S3_ENDPOINT_URL=http://minio
API_AWS_S3_CREATE_BUCKET=1

MINIO_ACCESS_KEY=minio_access
MINIO_SECRET_KEY=minio_secret
MINIO_HOST=0.0.0.0
MINIO_PORT=9000

NGINX_PORT=80

POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_PORT=5432
EOF
```

4. Start server

```
$ docker-compose up -d
```
