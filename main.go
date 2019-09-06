package main

import (
	"github.com/applichic/lynou/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// LoadConfiguration configurations
	util.LoadConfiguration()

	// Init database
	db, err := util.InitDB()
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

	router.Run() // listen and serve on 0.0.0.0:8080
}
