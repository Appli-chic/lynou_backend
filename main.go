package main

import (
	config2 "github.com/applichic/lynou/config"
	"github.com/applichic/lynou/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// LoadConfiguration configurations
	config2.LoadConfiguration()
	util.LoginToStorage()

	// Init database
	db, err := config2.InitDB()
	defer db.Close()

	if err != nil {
		panic(err)
	}

	router := InitRouter()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://lynou.com"}
	router.Use(cors.New(config))

	err = router.Run() // listen and serve on 0.0.0.0:8080

	if err != nil {
		panic(err)
	}
}
