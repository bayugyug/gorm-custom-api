package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/bayugyug/gorm-custom-api/api/routes"
	"github.com/bayugyug/gorm-custom-api/configs"
)

//init internal system initialize
func init() {
	//uniqueness
	rand.Seed(time.Now().UnixNano())
}

func main() {
	start := time.Now()
	log.Println(configs.APIVersion)

	//init
	appcfg := configs.NewAppSettings()

	//check
	if appcfg.Config == nil {
		log.Fatal("Oops! Config missing")
	}
	//init service
	service, err := routes.NewAPIRouter(
		routes.WithSvcOptConfig(appcfg.Config),
	)
	if err != nil {
		log.Fatal("Oops! config might be missing", err)
	}
	//run service
	service.Run()
	log.Println("Since", time.Since(start))
	log.Println("Done")
}
