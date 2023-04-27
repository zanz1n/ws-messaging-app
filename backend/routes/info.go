package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zanz1n/ws-messaging-app/utils"
)

type RouteInfo struct {
	Path   string   `json:"path"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

var validMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

func GetRoot() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		routes := c.App().GetRoutes()

		info := []RouteInfo{}

		for _, r := range routes {
			if utils.Includes(validMethods, r.Method) {
				info = append(info, RouteInfo{
					Path:   r.Path,
					Method: r.Method,
					Params: r.Params,
				})
			}
		}

		return c.JSON(fiber.Map{
			"message": "this route is empty, here is a lis of all routes",
			"data":    info,
		})
	}
}
