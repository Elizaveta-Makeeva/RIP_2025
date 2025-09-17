package pkg

import (
   "fmt"

   "github.com/gin-gonic/gin"
   "github.com/sirupsen/logrus"
   "lab2/internal/app/config"
   "lab2/internal/app/revenueHandler"
)

type Application struct {
   Config  *config.Config
   Router  *gin.Engine
   Handler *revenueHandler.RevenueHandler
}

func NewApp(c *config.Config, r *gin.Engine, h *revenueHandler.RevenueHandler) *Application {
   return &Application{
      Config:  c,
      Router:  r,
      Handler: h,
   }
}

func (a *Application) RunApp() {
   logrus.Info("Server start up")

   a.Handler.RegisterRevenueHandler(a.Router)
   a.Handler.RegisterStatic(a.Router)

   serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
   if err := a.Router.Run(serverAddress); err != nil {
      logrus.Fatal(err)
   }
   logrus.Info("Server down")
}