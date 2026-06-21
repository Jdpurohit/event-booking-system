package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a *App) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Role Role   `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	token := uuid.New().String()
	u := &User{Name: req.Name, Role: req.Role, Token: token}
	a.Store.CreateUser(u)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func (a *App) HandleListEvents(w http.ResponseWriter, r *http.Request) {
	events := a.Store.ListEvents()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (a *App) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Capacity    int    `json:"capacity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	e := &Event{OrganizerID: user.ID, Title: req.Title, Description: req.Description, Capacity: req.Capacity}
	a.Store.CreateEvent(e)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

func (a *App) HandleGetEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	e, err := a.Store.GetEvent(id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

func (a *App) HandleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	var req Event
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	// verify organizer owns the event
	ev, err := a.Store.GetEvent(id)
	if err != nil || ev.OrganizerID != user.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	updated, err := a.Store.UpdateEvent(id, &req)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	// enqueue notification to customers who booked this event
	a.Jobs.Enqueue(Job{Type: JobTypeEventUpdate, Payload: map[string]any{"event_id": id}})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (a *App) HandleBookEvent(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	b := &Booking{EventID: id, CustomerID: user.ID, Quantity: req.Quantity}
	booking, err := a.Store.CreateBooking(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// enqueue booking confirmation job
	a.Jobs.Enqueue(Job{Type: JobTypeBookingConfirmation, Payload: map[string]any{"booking_id": booking.ID}})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

func (a *App) HandleListBookings(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var out []*Booking
	switch user.Role {
	case RoleCustomer:
		out = a.Store.GetBookingsByUser(user.ID)
	case RoleOrganizer:
		out = a.Store.GetBookingsForOrganizer(user.ID)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}
