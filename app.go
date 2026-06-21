package main

type App struct {
	Store *Store
	Jobs  *JobQueue
}

func NewApp() *App {
	store := NewStore()
	jobs := NewJobQueue()
	jobs.Start()
	return &App{Store: store, Jobs: jobs}
}
