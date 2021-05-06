package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var tpl = template.Must(template.ParseFiles("index.html"))

type PageVariables struct {
	Groups       []*Group
	CurrentGroup *Group
	GroupTasks   []*Task
	CurrentTask  *Task
}

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	ds := DataStore{}
	ds.Open()
	defer ds.Close()
	dao := GroupsDao{ds: &ds}
	groups := dao.GetGroups()
	//ids := dao.getGroupsIds() // just test
	//fmt.Println(ids)
	vars := PageVariables{Groups: groups}
	err := tpl.Execute(w, vars)
	if err != nil {
		panic(err)
	}
}

func groupCreateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	name := r.Form.Get("g_name")
	ds := DataStore{}
	ds.Open()
	defer ds.Close()
	dao := GroupsDao{ds: &ds}
	gr := Group{}
	gr.GName = name
	dao.CreateGroup(&gr)
	http.Redirect(w, r, "/", http.StatusFound)
}

func groupUpdateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	gid := r.Form.Get("g_id")
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		panic(err)
	}
	name := r.Form.Get("g_name")
	ds := DataStore{}
	ds.Open()
	defer ds.Close()
	dao := GroupsDao{ds: &ds}
	gr := dao.ReadGroup(gId)
	gr.GName = name
	dao.UpdateGroup(&gr)
	http.Redirect(w, r, fmt.Sprintf("/group/tasks?g_id=%d", gId), http.StatusFound)
}

func groupDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	gId := r.Form.Get("g_id")
	gid, err := strconv.ParseInt(gId, 10, 64)
	if err != nil {
		panic(err)
	}
	ds := DataStore{}
	ds.Open()
	defer ds.Close()
	dao := GroupsDao{ds: &ds}
	dao.DeleteGroup(gid)
	http.Redirect(w, r, "/", http.StatusFound)
}

func groupTasksHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	gid := r.Form.Get("g_id")
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		panic(err)
	}
	ds := DataStore{}
	ds.Open()
	defer ds.Close()
	dao := GroupsDao{ds: &ds}
	groups := dao.GetGroups()
	currGroup := dao.ReadGroup(gId)
	tDao := TasksDao{ds: &ds}
	tasks := tDao.GetGroupTasks(gId)
	vars := PageVariables{Groups: groups, CurrentGroup: &currGroup, GroupTasks: tasks}
	err = tpl.Execute(w, vars)
	if err != nil {
		panic(err)
	}
}

func taskEditHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	tid := r.Form.Get("t_id")
	tId, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		panic(err)
	}
	ds := DataStore{}
	ds.Open()
	defer ds.Close()
	gDao := GroupsDao{ds: &ds}
	groups := gDao.GetGroups()
	tDao := TasksDao{ds: &ds}
	currTask := tDao.ReadTask(tId)
	currGroup := gDao.ReadGroup(currTask.GId)
	tasks := tDao.GetGroupTasks(currTask.GId)
	vars := PageVariables{Groups: groups, CurrentGroup: &currGroup, GroupTasks: tasks, CurrentTask: &currTask}
	err = tpl.Execute(w, vars)
	if err != nil {
		panic(err)
	}
}

func taskHideHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	gId := r.Form.Get("g_id")
	http.Redirect(w, r, fmt.Sprintf("/group/tasks?g_id=%s", gId), http.StatusFound)
}

func taskCreateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	gid := r.Form.Get("g_id")
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		panic(err)
	}
	subject := r.Form.Get("t_subject")
	ds := DataStore{}
	ds.Open()
	defer ds.Close()
	tDao := TasksDao{ds: &ds}
	t := Task{}
	t.GId = gId
	t.TSubject = subject
	t.TPriority = 1
	currentTime := time.Now().Local()
	layoutISO := currentTime.Format("2006-01-02")
	t.TDate = layoutISO
	tDao.CreateTask(&t)
	http.Redirect(w, r, fmt.Sprintf("/group/tasks?g_id=%s", gid), http.StatusFound)
}

func taskUpdateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	tid := r.Form.Get("t_id")
	tId, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		panic(err)
	}
	date := r.Form.Get("t_date")
	_, err2 := time.Parse("2006-01-02", date)
	if err2 != nil {
		panic(err)
	}
	subject := r.Form.Get("t_subject")
	priority := r.Form.Get("t_priority")
	tPriority, err := strconv.ParseInt(priority, 10, 64)
	if err != nil {
		panic(err)
	}
	comments := r.Form.Get("t_comments")
	ds := DataStore{}
	ds.Open()
	defer ds.Close()
	tDao := TasksDao{ds: &ds}
	t := tDao.ReadTask(tId)
	t.TSubject = subject
	t.TPriority = tPriority
	t.TDate = date
	t.TComments = comments
	tDao.UpdateTask(&t)
	http.Redirect(w, r, fmt.Sprintf("/task/edit?t_id=%d", t.TId), http.StatusFound)
}

func taskDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	tid := r.Form.Get("t_id")
	tId, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		panic(err)
	}
	ds := DataStore{}
	ds.Open()
	defer ds.Close()
	tDao := TasksDao{ds: &ds}
	t := tDao.ReadTask(tId)
	gId := t.GId
	tDao.DeleteTask(tId)
	http.Redirect(w, r, fmt.Sprintf("/group/tasks?g_id=%d", gId), http.StatusFound)
}
