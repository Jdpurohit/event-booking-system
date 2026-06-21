package main

type Role string

const (
	RoleOrganizer Role = "organizer"
	RoleCustomer  Role = "customer"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Role  Role   `json:"role"`
	Token string `json:"token"`
}

type Event struct {
	ID               int    `json:"id"`
	OrganizerID      int    `json:"organizer_id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Capacity         int    `json:"capacity"`
	AvailableTickets int    `json:"available_tickets"`
}

type Booking struct {
	ID         int `json:"id"`
	EventID    int `json:"event_id"`
	CustomerID int `json:"customer_id"`
	Quantity   int `json:"quantity"`
}
