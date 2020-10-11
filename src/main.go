package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	project := os.Getenv("PUB_PROJECT")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/order01", handler(project, "order01"))
	e.GET("/order02", handler(project, "order02"))
	e.GET("/order03", handler(project, "order03"))
	e.GET("/order04", handler(project, "order04"))
	e.GET("/order05", handler(project, "order05"))

	e.Logger.Fatal(e.Start(":8080"))
}

func handler(project, path string) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		operationSequence := createOperationSequence()

		client, err := pubsub.NewClient(ctx, project)
		if err != nil {
			log.Fatal(err)
		}

		topic := client.Topic("my-normal-topic")
		defer topic.Stop()

		topic.EnableMessageOrdering = true

		message := &pubsub.Message{
			OrderingKey: operationSequence,
			Data:        []byte(path + ":" + operationSequence),
		}
		r := topic.Publish(ctx, message)
		if r == nil {
			log.Fatal("failed to topic.Publish!")
		}
		log.Printf("%+v", r)

		return c.String(http.StatusOK, path+":"+operationSequence)
	}
}

func createOperationSequence() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
