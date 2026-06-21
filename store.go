package main

import (
	"errors"
	"sync"
)

type Store struct {
	mu sync.RWMutex

	users    map[int]*User
	events   map[int]*Event
	bookings map[int]*Booking

	nextUserID    int
	nextEventID   int
	nextBookingID int
}

func NewStore() *Store {
	return &Store{
		users:         make(map[int]*User),
		events:        make(map[int]*Event),
		bookings:      make(map[int]*Booking),
		nextUserID:    1,
		nextEventID:   1,
		nextBookingID: 1,
	}
}

func (s *Store) CreateUser(u *User) *User {
	s.mu.Lock()
	defer s.mu.Unlock()
	u.ID = s.nextUserID
	s.nextUserID++
	s.users[u.ID] = u
	return u
}

func (s *Store) GetUserByToken(token string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if u.Token == token {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (s *Store) CreateEvent(e *Event) *Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	e.ID = s.nextEventID
	s.nextEventID++
	e.AvailableTickets = e.Capacity
	s.events[e.ID] = e
	return e
}

func (s *Store) GetEvent(id int) (*Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.events[id]
	if !ok {
		return nil, errors.New("event not found")
	}
	// return a copy to avoid races
	copy := *e
	return &copy, nil
}

func (s *Store) UpdateEvent(id int, upd *Event) (*Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.events[id]
	if !ok {
		return nil, errors.New("event not found")
	}
	if upd.Title != "" {
		e.Title = upd.Title
	}
	if upd.Description != "" {
		e.Description = upd.Description
	}
	if upd.Capacity > 0 {
		// adjust available tickets relative to previous capacity
		diff := upd.Capacity - e.Capacity
		e.Capacity = upd.Capacity
		e.AvailableTickets += diff
		if e.AvailableTickets < 0 {
			e.AvailableTickets = 0
		}
	}
	copy := *e
	return &copy, nil
}

func (s *Store) ListEvents() []*Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Event, 0, len(s.events))
	for _, e := range s.events {
		copy := *e
		out = append(out, &copy)
	}
	return out
}

func (s *Store) CreateBooking(b *Booking) (*Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.events[b.EventID]
	if !ok {
		return nil, errors.New("event not found")
	}
	if b.Quantity <= 0 {
		return nil, errors.New("invalid quantity")
	}
	if e.AvailableTickets < b.Quantity {
		return nil, errors.New("not enough tickets")
	}
	e.AvailableTickets -= b.Quantity
	b.ID = s.nextBookingID
	s.nextBookingID++
	s.bookings[b.ID] = b
	return b, nil
}

func (s *Store) GetBookingsByUser(userID int) []*Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []*Booking{}
	for _, b := range s.bookings {
		if b.CustomerID == userID {
			copy := *b
			out = append(out, &copy)
		}
	}
	return out
}

func (s *Store) GetBookingsByEvent(eventID int) []*Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []*Booking{}
	for _, b := range s.bookings {
		if b.EventID == eventID {
			copy := *b
			out = append(out, &copy)
		}
	}
	return out
}

func (s *Store) GetBookingsForOrganizer(organizerID int) []*Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []*Booking{}
	for _, b := range s.bookings {
		e, ok := s.events[b.EventID]
		if ok && e.OrganizerID == organizerID {
			copy := *b
			out = append(out, &copy)
		}
	}
	return out
}
