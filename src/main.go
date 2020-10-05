package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/gcppubsub"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	project := os.Getenv("PUB_PROJECT")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/add-company", func(c echo.Context) error {
		newID := createUUID()
		sendTopic(c.Request().Context(), project, "add-company:"+newID)
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
	topic, err := pubsub.OpenTopic(ctx, fmt.Sprintf("gcppubsub://projects/%s/topics/my-test-topic", project))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := topic.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	if err := topic.Send(ctx, &pubsub.Message{Body: []byte(order)}); err != nil {
		log.Fatal(err)
	}
}
