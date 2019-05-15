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

	- development
    - create user gorm_cust_api;
    - create database if not exists gorm_cust_api;
    - goose -dir=db/migration mysql "root:@tcp(127.0.0.1:3306)/gorm_cust_api?parseTime=true&loc=Local&charset=utf8" status

	- testing
    - create user gorm_cust_api;
    - create database if not exists gorm_cust_api_test;
    - goose -dir=db/migration mysql "root:@tcp(127.0.0.1:3306)/gorm_cust_api_test?parseTime=true&loc=Local&charset=utf8" status

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


#get all records
curl -X GET    'http://127.0.0.1:8989/v1/api/building'
{
  "status": "success",
  "result": [
    {
      "id": 1,
      "name": "dhaha-ddbuilding-283",
      "address": "address here dabis",
      "created_at": "2019-05-12T12:19:37+08:00",
      "updated_at": "2019-05-12T12:19:37+08:00",
      "floors": [
        {
          "floor": "floor-1"
        },
        {
          "floor": "floor-2e"
        },
        {
          "floor": "yfldoor-1"
        },
        {
          "floor": "zfloor-2e"
        }
      ]
    },
    {
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
    },
    {
      "id": 3,
      "name": "building-a",
      "address": "updated address here2",
      "created_at": "2019-05-12T12:49:21+08:00",
      "updated_at": "2019-05-12T12:49:47+08:00",
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
    }
  ],
  "total": 3
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

