package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"sdm_demo_go_todolist/dal"
	"strconv"
	"time"
)

func ReturnTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	tDao := dal.CreateTasksDao()
	currTask, err := tDao.ReadTask(tId)
	if err != nil {
		respondWith500(w, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(currTask)
	if err != nil {
		respondWith500(w, err.Error())
	}
}

func ReturnGroupTasksHandler(w http.ResponseWriter, r *http.Request) {
	// https://stackoverflow.com/questions/45378566/gorilla-mux-optional-query-values/45378656
	// https://stackoverflow.com/questions/46045756/retrieve-optional-query-variables-with-gorilla-mux
	urlParams := r.URL.Query()
	gid := urlParams.Get("g_id")
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		respondWithUriError(w, r, err)
		return
	}
	tDao := dal.CreateTasksDao()
	tasks, err := tDao.GetGroupTasks(gId)
	if err != nil {
		respondWith500(w, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		respondWith500(w, err.Error())
	}
}

func TaskCreateHandler(w http.ResponseWriter, r *http.Request) {
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
	var inTask dal.Task
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
	tDao := dal.CreateTasksDao()
	t := dal.Task{}
	t.GId = gId
	t.TSubject = subject
	t.TPriority = 1
	currentTime := time.Now().Local()
	layoutISO := currentTime.Format("2006-01-02")
	t.TDate = layoutISO
	err = tDao.CreateTask(&t)
	if err != nil {
		respondWith500(w, err.Error())
		return
	}
}

func TaskDeleteHandler(w http.ResponseWriter, r *http.Request) {
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
	tDao := dal.CreateTasksDao()
	_, err = tDao.DeleteTask(tId)
	if err != nil {
		respondWith500(w, err.Error())
		return
	}
}

func TaskUpdateHandler(w http.ResponseWriter, r *http.Request) {
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
	var inTask dal.Task
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
	tDao := dal.CreateTasksDao()
	t, err := tDao.ReadTask(tId)
	if err != nil {
		respondWith500(w, err.Error())
		return
	}
	t.TSubject = subject
	t.TPriority = priority
	t.TDate = date
	t.TComments = comments
	_, err = tDao.UpdateTask(&t)
	if err != nil {
		respondWith500(w, err.Error())
		return
	}
}
