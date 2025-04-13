package route

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"payton/controller"
)

func RegisterRoute(e *echo.Echo) {
	e.Any("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<strong>Hello, World!</strong>")
	})
	ton := e.Group("/api/v1")
	ton.POST("/sendTon", controller.Ctrl.SendTon)
}
