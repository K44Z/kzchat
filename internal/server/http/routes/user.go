package routes

import (
	h "github.com/K44Z/kzchat/internal/server/http/handlers"
	"github.com/K44Z/kzchat/internal/server/http/middleware"
	"github.com/K44Z/kzchat/internal/server/services"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(router fiber.Router, s *services.Services) {
	router.Get("/:username", middleware.Jwt, h.SearchForUserByUsernameHandler(s))
	router.Get("/usernames/all", middleware.Jwt, h.GetAllUsersHandler(s))
}
