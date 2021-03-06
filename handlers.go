package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

var tpl = template.Must(template.ParseFiles("index.html"))

type PageVariables struct {
	Groups       []Group
	CurrentGroup *Group
	GroupTasks   []Task
	CurrentTask  *Task
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	ds := DataStore{}
	ds.open()
	defer ds.close()
	dao := GroupsDao{ds: &ds}
	groups := dao.getGroups()
	//ids := dao.getGroupsIds() // just test
	//fmt.Println(ids)
	vars := PageVariables{Groups: groups}
	tpl.Execute(w, vars)
}

func groupCreateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("g_name")
	ds := DataStore{}
	ds.open()
	defer ds.close()
	dao := GroupsDao{ds: &ds}
	gr := Group{}
	gr.GName = name
	dao.createGroup(&gr)
	http.Redirect(w, r, "/", http.StatusFound)
}

func groupUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	gid := r.Form.Get("g_id")
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		log.Panic(err)
		return
	}
	name := r.Form.Get("g_name")
	ds := DataStore{}
	ds.open()
	defer ds.close()
	dao := GroupsDao{ds: &ds}
	gr := dao.readGroup(gId)
	gr.GName = name
	dao.updateGroup(&gr)
	http.Redirect(w, r, "/", http.StatusFound)
}

func groupDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	gId := r.Form.Get("g_id")
	gid, err := strconv.ParseInt(gId, 10, 64)
	if err != nil {
		log.Panic(err)
		return
	}
	ds := DataStore{}
	ds.open()
	defer ds.close()
	dao := GroupsDao{ds: &ds}
	dao.deleteGroup(gid)
	http.Redirect(w, r, "/", http.StatusFound)
}

func groupTasksHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	gid := r.Form.Get("g_id")
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		log.Panic(err)
		return
	}
	ds := DataStore{}
	ds.open()
	defer ds.close()
	dao := GroupsDao{ds: &ds}
	groups := dao.getGroups()
	currGroup := dao.readGroup(gId)
	tDao := TasksDao{ds: &ds}
	tasks := tDao.getGroupTasks(gId)
	vars := PageVariables{Groups: groups, CurrentGroup: &currGroup, GroupTasks: tasks}
	tpl.Execute(w, vars)
}

func taskEditHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tid := r.Form.Get("t_id")
	tId, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		log.Panic(err)
		return
	}
	ds := DataStore{}
	ds.open()
	defer ds.close()
	gDao := GroupsDao{ds: &ds}
	groups := gDao.getGroups()
	tDao := TasksDao{ds: &ds}
	currTask := tDao.readTask(tId)
	currGroup := gDao.readGroup(currTask.GId)
	tasks := tDao.getGroupTasks(currTask.GId)
	vars := PageVariables{Groups: groups, CurrentGroup: &currGroup, GroupTasks: tasks, CurrentTask: &currTask}
	tpl.Execute(w, vars)
}

func taskHideHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	gId := r.Form.Get("g_id")
	http.Redirect(w, r, fmt.Sprintf("/group/tasks?g_id=%s", gId), http.StatusFound)
}

func taskCreateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	gid := r.Form.Get("g_id")
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		log.Panic(err)
		return
	}
	subject := r.Form.Get("t_subject")
	ds := DataStore{}
	ds.open()
	defer ds.close()
	tDao := TasksDao{ds: &ds}
	t := Task{}
	t.GId = gId
	t.TSubject = subject
	t.TPriority = 1
	currentTime := time.Now().Local()
	layoutISO := currentTime.Format("2006-01-02")
	t.TDate = layoutISO
	tDao.createTask(&t)
	http.Redirect(w, r, fmt.Sprintf("/group/tasks?g_id=%s", gid), http.StatusFound)
}

func taskUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tid := r.Form.Get("t_id")
	tId, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		log.Panic(err)
		return
	}
	date := r.Form.Get("t_date")
	_, err2 := time.Parse("2006-01-02", date)
	if err2 != nil {
		log.Panic(err)
		return
	}
	subject := r.Form.Get("t_subject")
	priority  := r.Form.Get("t_priority")
	tPriority, err := strconv.ParseInt(priority, 10, 64)
	if err != nil {
		log.Panic(err)
		return
	}
	comments  := r.Form.Get("t_comments")
	ds := DataStore{}
	ds.open()
	defer ds.close()
	tDao := TasksDao{ds: &ds}
	t := tDao.readTask(tId)
	t.TSubject = subject
	t.TPriority = tPriority
	t.TDate = date
	t.TComments = comments
	tDao.updateTask(&t)
	http.Redirect(w, r, fmt.Sprintf("/task/edit?t_id=%d", t.TId), http.StatusFound)
}

func taskDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tid := r.Form.Get("t_id")
	tId, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		log.Panic(err)
		return
	}
	ds := DataStore{}
	ds.open()
	defer ds.close()
	tDao := TasksDao{ds: &ds}
	t := tDao.readTask(tId)
	gId := t.GId
	tDao.deleteTask(tId);
	http.Redirect(w, r, fmt.Sprintf("/group/tasks?g_id=%d", gId), http.StatusFound)
}
