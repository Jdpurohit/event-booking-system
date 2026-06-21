package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	app := NewApp()

	r := chi.NewRouter()

	r.Post("/signup", app.HandleSignup)

	r.Route("/events", func(r chi.Router) {
		r.Get("/", app.HandleListEvents)
		r.Post("/", app.AuthMiddleware(app.RequireRole(RoleOrganizer, app.HandleCreateEvent)))
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", app.HandleGetEvent)
			r.Put("/", app.AuthMiddleware(app.RequireRole(RoleOrganizer, app.HandleUpdateEvent)))
			r.Post("/book", app.AuthMiddleware(app.RequireRole(RoleCustomer, app.HandleBookEvent)))
		})
	})

	r.Get("/bookings", app.AuthMiddleware(app.HandleListBookings))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
