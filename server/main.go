package main

import (
	"flag"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/tmaiaroto/discfg/config"
	"log"
)

var Options = config.Options{StorageInterfaceName: "dynamodb", Version: "0.6.0"}

func main() {
	// TODO: remove
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	port := *flag.String("port", "8899", "API port")
	apiVersion := *flag.String("version", "v1", "API version")
	region := *flag.String("region", "us-east-1", "AWS region")

	Options.Storage.AWS.Region = region

	e := echo.New()

	// Middleware
	//e.Use(mw.Logger())
	e.Use(mw.Gzip())
	e.Use(mw.Recover())

	// Routes
	switch apiVersion {
	default:
	case "v1":
		v1Routes(e)
	}

	// Start server
	e.Run(":" + port)
}
