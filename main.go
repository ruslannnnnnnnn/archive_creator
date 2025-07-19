package main

import (
	"archive_creator/internal/api"
	"archive_creator/internal/config"
	"github.com/gin-gonic/gin"
	"strconv"
)

const configPath = "/home/ruslan/goprojects/gotest/config.yaml"

func main() {

	appConfig, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	api.InitRoutes(r, appConfig)

	r.Run(":" + strconv.Itoa(appConfig.Port))
}
