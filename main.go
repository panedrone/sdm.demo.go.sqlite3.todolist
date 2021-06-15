package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sdm_demo_go_todolist/api"
	"sdm_demo_go_todolist/dal"
)

func handleAssets(r *mux.Router) {
	staticFileDirectory := http.Dir("./assets/")
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for
	// "./assets/assets/index.html", and yield an error
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
}

func main() {
	err := dal.OpenDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		err = dal.CloseDB()
	}()

	myRouter := mux.NewRouter()
	myRouter.StrictSlash(true)
	handleAssets(myRouter)
	//fs := http.FileServer(http.Dir("assets"))
	//myRouter.Handle("/assets/", http.StripPrefix("/assets/", fs))
	////////////////////
	myRouter.HandleFunc("/groups", api.ReturnAllGroupsHandler).Methods("GET")
	myRouter.HandleFunc("/groups", api.GroupCreateHandler).Methods("POST")
	myRouter.HandleFunc("/groups/{g_id}", api.ReturnGroupHandler).Methods("GET")
	myRouter.HandleFunc("/groups/{g_id}", api.GroupUpdateHandler).Methods("PUT")
	myRouter.HandleFunc("/groups/{g_id}", api.GroupDeleteHandler).Methods("DELETE")
	////////////////////
	myRouter.HandleFunc("/tasks", api.ReturnGroupTasksHandler).Queries("g_id", "{g_id}").Methods("GET")
	myRouter.HandleFunc("/tasks", api.TaskCreateHandler).Queries("g_id", "{g_id}").Methods("POST")
	myRouter.HandleFunc("/tasks/{t_id}", api.ReturnTaskHandler).Methods("GET")
	myRouter.HandleFunc("/tasks/{t_id}", api.TaskUpdateHandler).Methods("PUT")
	myRouter.HandleFunc("/tasks/{t_id}", api.TaskDeleteHandler).Methods("DELETE")
	// log.Fatal
	// https://blog.scottlogic.com/2017/02/28/building-a-web-app-with-go.html
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
