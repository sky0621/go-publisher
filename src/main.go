package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
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

	e.GET("/add-region", regionHandler(project, "add-region"))
	e.GET("/edit-region/:rid", regionHandler(project, "edit-region"))

	e.GET("/add-school/:rid", schoolHandler(project, "add-school"))
	e.GET("/edit-school/:rid/:sid", schoolHandler(project, "edit-school"))

	e.GET("/add-grade/:rid/:sid", handler(project, "add-grade"))
	e.GET("/add-class/:rid/:sid", handler(project, "add-class"))
	e.GET("/add-teacher/:rid/:sid", handler(project, "add-teacher"))
	e.GET("/add-student/:rid/:sid", handler(project, "add-student"))

	e.Logger.Fatal(e.Start(":8080"))
}

func regionHandler(project, path string) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		regionID := c.Param("rid")
		if regionID == "" {
			regionID = createUUID()
		}
		operationSequence := createOperationSequence()

		client, err := firestore.NewClient(ctx, project)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if client != nil {
				if err := client.Close(); err != nil {
					log.Fatal(err)
				}
			}
		}()

		_, err = client.Collection("region").Doc(regionID).
			Set(ctx, map[string]interface{}{
				"operationSequence": operationSequence,
				"order":             path + ":" + regionID,
				"isDone":            false,
			}, firestore.MergeAll)
		if err != nil {
			log.Fatal(err)
		}

		sendTopic(ctx, project, path+":"+regionID, operationSequence)
		return c.String(http.StatusOK, regionID)
	}
}

func schoolHandler(project, path string) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		regionID := c.Param("rid")
		if regionID == "" {
			log.Fatal("no regionID")
		}
		schoolID := c.Param("sid")
		if schoolID == "" {
			schoolID = createUUID()
		}
		operationSequence := createOperationSequence()

		client, err := firestore.NewClient(ctx, project)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if client != nil {
				if err := client.Close(); err != nil {
					log.Fatal(err)
				}
			}
		}()

		_, err = client.Collection("region").Doc(regionID).
			Collection("school").Doc(schoolID).
			Collection("operation").Doc(operationSequence).
			Set(ctx, map[string]interface{}{
				"operationSequence": operationSequence,
				"order":             path + ":" + regionID + ":" + schoolID,
				"isDone":            false,
			}, firestore.MergeAll)
		if err != nil {
			log.Fatal(err)
		}

		sendTopic(ctx, project, path+":"+regionID+":"+schoolID, operationSequence)
		return c.String(http.StatusOK, regionID+":"+schoolID)
	}
}

func handler(project, path string) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		regionID := c.Param("rid")
		if regionID == "" {
			log.Fatal("no regionID")
		}
		schoolID := c.Param("sid")
		if schoolID == "" {
			log.Fatal("no schoolID")
		}
		newID := createUUID()
		operationSequence := createOperationSequence()

		client, err := firestore.NewClient(ctx, project)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if client != nil {
				if err := client.Close(); err != nil {
					log.Fatal(err)
				}
			}
		}()

		_, err = client.Collection("region").Doc(regionID).
			Collection("school").Doc(schoolID).
			Collection("operation").Doc(operationSequence).
			Set(ctx, map[string]interface{}{
				"operationSequence": operationSequence,
				"order":             path + ":" + regionID + ":" + schoolID + ":" + newID,
				"isDone":            false,
			}, firestore.MergeAll)
		if err != nil {
			log.Fatal(err)
		}

		sendTopic(ctx, project, path+":"+regionID+":"+schoolID+":"+newID, operationSequence)
		return c.String(http.StatusOK, regionID+":"+schoolID+":"+newID)
	}
}

func createUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
	return u.String()
}

func sendTopic(ctx context.Context, project, order, serializeKey string) {
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatal(err)
	}

	topic := client.Topic("my-test-topic")
	defer topic.Stop()

	topic.EnableMessageOrdering = true

	message := &pubsub.Message{
		OrderingKey: serializeKey,
		Data:        []byte(order + ":" + serializeKey),
	}
	r := topic.Publish(ctx, message)
	if r == nil {
		log.Fatal("failed to topic.Publish!")
	}
	log.Printf("%+v", r)
}

func createOperationSequence() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
