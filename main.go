package main

import (
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	routes(e)
	//起動
	e.Start(":8888")

}

func routes(e *echo.Echo) {
	e.POST("/", Send)
}
