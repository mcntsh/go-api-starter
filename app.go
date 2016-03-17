package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/mcntsh/go-api-starter/helpers/service"
	"github.com/mcntsh/go-api-starter/user_service"
)

func main() {
	app := service.NewApp()

	userServ, err := userserv.StartService()
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Could not load user service – %s", err))
	}

	app.NewService("/users", userServ)

	app.Listen(":8080")
}
