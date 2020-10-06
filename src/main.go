package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	project := os.Getenv("PUB_PROJECT")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/add-school", func(c echo.Context) error {
		newID := createUUID()
		sendTopic(c.Request().Context(), project, "add-school:"+newID)
		return c.String(http.StatusOK, newID)
	})

	e.GET("/add-grade", func(c echo.Context) error {
		newID := createUUID()
		sendTopic(c.Request().Context(), project, "add-grade:"+newID)
		return c.String(http.StatusOK, newID)
	})

	e.GET("/add-class", func(c echo.Context) error {
		newID := createUUID()
		sendTopic(c.Request().Context(), project, "add-class:"+newID)
		return c.String(http.StatusOK, newID)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func createUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
	return u.String()
}

func sendTopic(ctx context.Context, project, order string) {
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatal(err)
	}

	topic := client.Topic("my-test-topic")
	defer topic.Stop()

	topic.EnableMessageOrdering = true

	message := &pubsub.Message{
		OrderingKey: fmt.Sprintf("%d", time.Now().UnixNano()),
		Data:        []byte(order),
	}
	r := topic.Publish(ctx, message)
	if r == nil {
		log.Fatal("failed to topic.Publish!")
	}
	log.Printf("%+v", r)
}
