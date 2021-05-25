package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func returnTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tid, ok := vars["t_id"]
	if !ok {
		respondWithBadURI(w, r)
		return
	}
	tId, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		respondWithUriError(w, r, err)
		return
	}
	tDao := TasksDao{ds: &ds}
	currTask := tDao.ReadTask(tId)
	err = json.NewEncoder(w).Encode(currTask)
	if err != nil {
		panic(err)
	}
}

func returnGroupTasksHandler(w http.ResponseWriter, r *http.Request) {
	// https://stackoverflow.com/questions/45378566/gorilla-mux-optional-query-values/45378656
	// https://stackoverflow.com/questions/46045756/retrieve-optional-query-variables-with-gorilla-mux
	urlParams := r.URL.Query()
	gid := urlParams.Get("g_id")
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		respondWithUriError(w, r, err)
		return
	}
	tDao := TasksDao{ds: &ds}
	tasks := tDao.GetGroupTasks(gId)
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		panic(err)
	}
}

func taskCreateHandler(w http.ResponseWriter, r *http.Request) {
	// https://stackoverflow.com/questions/45378566/gorilla-mux-optional-query-values/45378656
	// https://stackoverflow.com/questions/46045756/retrieve-optional-query-variables-with-gorilla-mux
	urlParams := r.URL.Query()
	gid := urlParams.Get("g_id")
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		respondWithUriError(w, r, err)
		return
	}
	bodyBytes, err := ioutil.ReadAll(r.Body) // === panedrone: store source for error handling
	if err != nil {
		respondWithUriError(w, r, err)
		return
	}
	rd := bytes.NewReader(bodyBytes) // === panedrone: r.Body became unavailable
	decoder := json.NewDecoder(rd)
	var inTask Task
	err = decoder.Decode(&inTask)
	if err != nil {
		respondWithBadRequestError(w, fmt.Sprintf("JSON decoder FAIL: %s. Input: %s", err.Error(), bodyBytes))
		return
	}
	subject := inTask.TSubject
	if len(subject) == 0 {
		respondWithBadRequestError(w, fmt.Sprintf("Invalid input: %s", bodyBytes))
		return
	}
	tDao := TasksDao{ds: &ds}
	t := Task{}
	t.GId = gId
	t.TSubject = subject
	t.TPriority = 1
	currentTime := time.Now().Local()
	layoutISO := currentTime.Format("2006-01-02")
	t.TDate = layoutISO
	tDao.CreateTask(&t)
}

func taskDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tid, ok := vars["t_id"]
	if !ok {
		respondWithBadURI(w, r)
		return
	}
	tId, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		respondWithUriError(w, r, err)
		return
	}
	tDao := TasksDao{ds: &ds}
	tDao.DeleteTask(tId)
}

func taskUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tid, ok := vars["t_id"]
	if !ok {
		respondWithBadURI(w, r)
		return
	}
	tId, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		respondWithUriError(w, r, err)
		return
	}
	bodyBytes, err := ioutil.ReadAll(r.Body) // === panedrone: store source for error handling
	if err != nil {
		respondWithBadRequestError(w, fmt.Sprintf("Invalid request: %s", err.Error()))
		return
	}
	rd := bytes.NewReader(bodyBytes) // === panedrone: r.Body became unavailable
	decoder := json.NewDecoder(rd)
	var inTask Task
	err = decoder.Decode(&inTask)
	if err != nil {
		respondWithBadRequestError(w, fmt.Sprintf("JSON decoder FAIL: %s. Input: %s", err.Error(), bodyBytes))
		return
	}
	date := inTask.TDate
	_, err = time.Parse("2006-01-02", date)
	if err != nil {
		respondWithBadRequestError(w, fmt.Sprintf("Invalid input: %s -> %s", bodyBytes, err.Error()))
		return
	}
	subject := inTask.TSubject
	if len(subject) == 0 {
		respondWithBadRequestError(w, fmt.Sprintf("Invalid input: %s", bodyBytes))
		return
	}
	priority := inTask.TPriority
	comments := inTask.TComments
	tDao := TasksDao{ds: &ds}
	t := tDao.ReadTask(tId)
	t.TSubject = subject
	t.TPriority = priority
	t.TDate = date
	t.TComments = comments
	tDao.UpdateTask(&t)
}
