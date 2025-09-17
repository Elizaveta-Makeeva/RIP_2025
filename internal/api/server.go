package api

import (
  "github.com/gin-gonic/gin"
  "github.com/sirupsen/logrus"
  "log"
  "lab1/internal/app/revenueHandler"
  "lab1/internal/app/revenueModel"
)

func StartServer() {
  log.Println("Starting server")

  rm, err := revenueModel.NewRevenueModel()
  if err != nil {
    logrus.Error("ошибка инициализации репозитория")
  }

  revenueHandler := revenueHandler.NewRevenueHandler(rm)

  r := gin.Default()
  r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

  r.GET("/home", revenueHandler.GetPeriods)
  r.GET("/period/:id", revenueHandler.GetPeriod)
  r.GET("/application/:id", revenueHandler.GetPeriodsApplication)

  r.Run() 
  log.Println("Server down")
}