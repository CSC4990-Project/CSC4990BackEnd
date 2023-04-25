package routes

import (
	"github.com/CSC4990-Project/CSC4990BackEnd/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Post("/api/Login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Get("/api/userTicket/:user", controllers.UserTicketView)
	app.Post("/api/logout", controllers.Logout)
	app.Get("/api/tickets", controllers.AdminTicketView)
	app.Get("/api/tickets/:id", controllers.DetailedView)
	app.Post("api/submit", controllers.SubmitTicket)
	app.Post("api/update/:id", controllers.UpdateTicket)

}
