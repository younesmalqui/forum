package handlers

import "forum/config"

type Page struct {
	Title string
	User  *config.Session
	Data  any
}

func NewPageStruct(title string, session string, data any) *Page {
	return &Page{
		Title: title,
		User:  config.IsAuth(session),
		Data:  data,
	}
}
