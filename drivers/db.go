package drivers

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	//needed for mysql db connection
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (

	// DBParams all config options
	DBParams = map[string]string{
		"charset":          "utf8",
		"tls":              "preferred",
		"parseTime":        "true",
		"loc":              "Local",
		"maxAllowedPacket": "0",
		"readTimeout":      "1m30s",
		"writeTimeout":     "1m",
	}
)

// DBHandle the db main object
type DBHandle struct {
	Connection  []*gorm.DB
	Namespace   string
	Type        string
	connections []*gorm.DB
}

// NewDBHandle get a db handler object
func NewDBHandle(name, dbtype, conn string) *DBHandle {
	//create
	handle := &DBHandle{
		Namespace: name,
		Type:      dbtype,
	}
	handle.PrepareDB(dbtype, conn)
	return handle
}

// PrepareDB connect to db
//
//  dbtype:
//       mysql
//  conn:
//       root:@tcp(127.0.0.1:3306)/db-name
func (dh *DBHandle) PrepareDB(dbtype, conn string) *DBHandle {

	var opts []string
	for k, v := range DBParams {
		opts = append(opts, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	optss := strings.Join(opts, "&")
	log.Println(optss)

	for i := 1; i <= 3; i++ {
		dbh, err := gorm.Open("mysql", fmt.Sprintf("%s?%s", conn, optss))
		if err != nil {
			log.Println(err, conn)
			time.Sleep(time.Millisecond * 1000)
			continue
		}

		// most important tweak is here :-)
		// https://www.alexedwards.net/blog/configuring-sqldb

		// MySQL's wait_timeout setting will automatically close any connections
		// that haven't been used for 8 hours (by default).

		// Set the maximum lifetime of a connection to 1 hour. Setting it to 0
		// means that there is no maximum lifetime and the connection is reused
		// forever (which is the default behavior).

		dbh.DB().SetConnMaxLifetime(time.Hour)

		// Set the maximum number of concurrently idle connections to 5. Setting this
		// to less than or equal to 0 will mean that no idle connections are retained.

		dbh.DB().SetMaxIdleConns(2)

		// Set the maximum number of concurrently open connections to 5. Setting this
		// to less than or equal to 0 will mean there is no maximum limit (which
		// is also the default setting).

		dbh.DB().SetMaxOpenConns(3)

		// Set the number of open and idle connection to a maximum total of (idle:2 + open:3) = 5

		//save it
		dh.connections = append(dh.connections, dbh)
	}

	log.Println("Connection Established#", len(dh.connections))

	return dh
}

// GetConnection get 1 connection
func (dh *DBHandle) GetConnection() *gorm.DB {
	if len(dh.connections) <= 0 {
		log.Println("ERROR: db connection empty")
		return nil
	}
	p := rand.Intn(len(dh.connections))
	if len(dh.connections) <= 0 || dh.connections[p] == nil {
		log.Println("ERROR: db connection empty")
		return nil
	}
	return dh.connections[p]
}