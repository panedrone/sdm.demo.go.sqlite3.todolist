package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/group/create", groupCreateHandler)
	mux.HandleFunc("/group/update", groupUpdateHandler)
	mux.HandleFunc("/group/delete", groupDeleteHandler)
	mux.HandleFunc("/group/tasks", groupTasksHandler)
	mux.HandleFunc("/task/edit", taskEditHandler)
	mux.HandleFunc("/task/hide", taskHideHandler)
	mux.HandleFunc("/task/create", taskCreateHandler)
	mux.HandleFunc("/task/update", taskUpdateHandler)
	mux.HandleFunc("/task/delete", taskDeleteHandler)
	// log.Fatal
	// https://blog.scottlogic.com/2017/02/28/building-a-web-app-with-go.html
	log.Fatal(http.ListenAndServe(":8080", mux))
}
