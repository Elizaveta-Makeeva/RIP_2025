package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab2/internal/app/config"
	"lab2/internal/app/dsn"
	"lab2/internal/app/revenueHandler"
	"lab2/internal/app/revenueModel"
	"lab2/internal/pkg"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := revenueModel.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing revenueModel: %v", errRep)
	}

	hand := revenueHandler.NewRevenueHandler(rep)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}