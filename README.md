## gorm-custom-api



* [x] Sample golang rest api that simulates CRUD with mysql-db (GORM)



### Dependencies

- Dependencies manager: 
  - GO111MODULE=on go mod init

- Testing framework: 
  - $ go get github.com/onsi/ginkgo/ginkgo
  - $ go get github.com/onsi/gomega/... 

- DB Migration:
  - $ go get -u github.com/pressly/goose/cmd/goose



### Prerequisite
- Create sample database and its necessary privileges

	- Development/Production
      - create user gorm_cust_api;
      - create database if not exists gorm_cust_api;
      - goose -dir=db/migration mysql "user:xxxx@tcp(127.0.0.1:3306)/gorm_cust_api?parseTime=true&loc=Local&charset=utf8" status

	- Testing
      - create user gorm_cust_api;
      - create database if not exists gorm_cust_api_test;
      - goose -dir=db/migration mysql "user:xxxx@tcp(127.0.0.1:3306)/gorm_cust_api_test?parseTime=true&loc=Local&charset=utf8" status

### Compile

```sh

     git clone https://github.com/bayugyug/gorm-custom-api.git && cd gorm-custom-api

     git pull && make clean && make

```
 


### End-Points-Url


```go

#create
curl -X POST    'http://127.0.0.1:8989/v1/api/building' -d '{"name":"new-building-a","address":"address here","floors":["floor-1","floor-2"]}'
{
  "status": "success",
  "result": 3
}


#update
curl -X PUT    'http://127.0.0.1:8989/v1/api/building' -d '{"id":3,"name":"building-a","address":"updated address here2","floors":["floor-a1","floor-a2","floor-a3"]}'
{
  "status": "success"
}


#get a record
curl -X GET    'http://127.0.0.1:8989/v1/api/building/2'	
{
  "status": "success",
  "result": {
    "id": 2,
    "name": "building-a",
    "address": "address here",
    "created_at": "2019-05-12T12:48:27+08:00",
    "updated_at": "2019-05-12T12:48:27+08:00",
    "floors": [
      {
        "floor": "floor-1"
      },
      {
        "floor": "floor-2"
      }
    ]
  }
}


#get set of records
curl -X GET    'http://127.0.0.1:8989/v1/api/building?page=1&limit=5'
{
  "status": "success",
  "result": {
    "items": [
      {
        "id": 1216,
        "name": "new-building-a2",
        "address": "updated address here2",
        "created_at": "2019-05-18T12:27:15+08:00",
        "updated_at": "2019-05-18T12:29:10+08:00",
        "floors": [
          {
            "floor": "floor-a1"
          },
          {
            "floor": "floor-a2"
          },
          {
            "floor": "floor-a3"
          }
        ]
      },
      {
        "id": 1217,
        "name": "new-18136-4653-26351",
        "address": "address here::26931",
        "created_at": "2019-05-18T12:31:28+08:00",
        "updated_at": "2019-05-18T12:31:28+08:00",
        "floors": [
          {
            "floor": "floor-1-7934"
          },
          {
            "floor": "floor-2-17795"
          }
        ]
      },
      {
        "id": 1218,
        "name": "new-2612-31403-30846",
        "address": "address here::8888",
        "created_at": "2019-05-18T12:31:29+08:00",
        "updated_at": "2019-05-18T12:31:29+08:00",
        "floors": [
          {
            "floor": "floor-1-4892"
          },
          {
            "floor": "floor-2-14444"
          }
        ]
      },
      {
        "id": 1219,
        "name": "new-4231-18081-30101",
        "address": "address here::13645",
        "created_at": "2019-05-18T12:31:29+08:00",
        "updated_at": "2019-05-18T12:31:29+08:00",
        "floors": [
          {
            "floor": "floor-1-23246"
          },
          {
            "floor": "floor-2-8658"
          }
        ]
      },
      {
        "id": 1220,
        "name": "new-1009-943-181",
        "address": "address here::1332",
        "created_at": "2019-05-18T12:31:29+08:00",
        "updated_at": "2019-05-18T12:31:29+08:00",
        "floors": [
          {
            "floor": "floor-1-8116"
          },
          {
            "floor": "floor-2-28677"
          }
        ]
      }
    ],
    "page": 1,
    "limit": 5,
    "total": 102
  }
}


#delete a record
curl -X DELETE    'http://127.0.0.1:8989/v1/api/building/3'
{
  "status": "success"
}

#check health of endpoint
curl -X GET 'http://127.0.0.1:8989/v1/api/health'
{
  "application": "Building API Service",
  "build": "20190512.141257",
  "commit": "c77effb",
  "release": "0.0.1",
  "now": "2019-05-12T14:13:27+08:00"
}

```


### Run

- The api can accept a json format configuration
	- Fields:
		- port      = port to run the http server (default: 8989)
		- dsn       = mysql connection string

- Sanity check
	- Either
		- ginkgo ./...	
		- make test

- From console

```sh

./bin/gorm-custom-api --config '{"port":"8989","dsn":"root:@tcp(127.0.0.1:3306)/gorm_cust_api"}'


```


### Notes

### Reference

### License

[MIT](https://bayugyug.mit-license.org/)

