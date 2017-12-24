package main

import (
	"fmt"

	"net/http"

	"io/ioutil"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.POST("/webhook/fb/", webhookHandler)
	e.Logger.Fatal(e.Start(":13964"))
}

func webhookHandler(c echo.Context) error {
	fp, err := c.FormParams()
	if err != nil {
		fmt.Printf("%#v\n", err)
		return c.JSON(http.StatusOK, nil)
	}
	fmt.Printf("[formParams]%#v\n", fp)
	fmt.Printf("[path]%#v\n", c.Path())
	fmt.Printf("[realIP]%#v\n", c.RealIP())
	fmt.Printf("[queryString]%#v\n", c.QueryString())

	req := c.Request()
	if req == nil {
		fmt.Println("req is nil")
		return c.JSON(http.StatusOK, nil)
	}

	fmt.Printf("[request]%#v\n", req)

	ioBody, err := req.GetBody()
	if err != nil {
		fmt.Printf("%#v\n", err)
		return c.JSON(http.StatusOK, nil)
	}
	bdary, err := ioutil.ReadAll(ioBody)
	if err != nil {
		fmt.Printf("%#v\n", err)
		return c.JSON(http.StatusOK, nil)
	}
	fmt.Printf("[body]%#v\n", bdary)

	return c.JSON(http.StatusOK, nil)
}
