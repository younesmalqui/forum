package handlers

type Response struct {
	Error   bool
	Success bool
	Message string
	Data    any
}
