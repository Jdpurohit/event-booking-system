package main

import (
	"fmt"
	"sync"
	"time"
)

type JobType string

const (
	JobTypeBookingConfirmation JobType = "booking_confirmation"
	JobTypeEventUpdate         JobType = "event_update_notification"
)

type Job struct {
	Type    JobType
	Payload map[string]any
}

type JobQueue struct {
	ch   chan Job
	once sync.Once
}

func NewJobQueue() *JobQueue {
	return &JobQueue{ch: make(chan Job, 100)}
}

func (q *JobQueue) Start() {
	q.once.Do(func() {
		go func() {
			for j := range q.ch {
				q.process(j)
			}
		}()
	})
}

func (q *JobQueue) Enqueue(j Job) {
	select {
	case q.ch <- j:
	default:
		// fallback: drop job if queue full
		fmt.Println("job queue full, dropping job", j.Type)
	}
}

func (q *JobQueue) process(j Job) {
	// simulate small delay
	time.Sleep(100 * time.Millisecond)
	switch j.Type {
	case JobTypeBookingConfirmation:
		idv, _ := j.Payload["booking_id"]
		fmt.Printf("[Job] Sending booking confirmation for booking_id=%v\n", idv)
	case JobTypeEventUpdate:
		evID, _ := j.Payload["event_id"]
		fmt.Printf("[Job] Notifying customers about event update event_id=%v\n", evID)
	default:
		fmt.Println("[Job] unknown job type", j.Type)
	}
}
