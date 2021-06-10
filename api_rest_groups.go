package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func groupCreateHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body) // === panedrone: store source for error handling
	if err != nil {
		respondWithUriError(w, r, err)
		return
	}
	rd := bytes.NewReader(bodyBytes) // === panedrone: r.Body became unavailable
	decoder := json.NewDecoder(rd)
	var inGroup Group
	err = decoder.Decode(&inGroup)
	if err != nil {
		respondWithBadRequestError(w, fmt.Sprintf("JSON decoder FAIL: %s. Input: %s", err.Error(), bodyBytes))
		return
	}
	name := inGroup.GName
	if len(name) == 0 {
		respondWithBadRequestError(w, fmt.Sprintf("Invalid input: %s", bodyBytes))
		return
	}
	dao := GroupsDao{ds: &ds}
	gr := Group{}
	gr.GName = name
	err = dao.CreateGroup(&gr)
	if err != nil {
		respondWith500(w, err.Error())
		return
	}
}

func returnAllGroupsHandler(w http.ResponseWriter, _ *http.Request) {
	dao := GroupsDao{ds: &ds}
	groups, err := dao.GetGroups()
	if err != nil {
		respondWith500(w, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(groups)
	if err != nil {
		respondWith500(w, err.Error())
	}
}

func groupUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gid, ok := vars["g_id"]
	if !ok {
		respondWithBadURI(w, r)
		return
	}
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
	var inGroup Group
	err = decoder.Decode(&inGroup)
	if err != nil {
		respondWithBadRequestError(w, fmt.Sprintf("JSON decoder FAIL: %s. Input: %s", err.Error(), bodyBytes))
		return
	}
	name := inGroup.GName
	if len(name) == 0 {
		respondWithBadRequestError(w, fmt.Sprintf("Invalid input: %s", bodyBytes))
		return
	}
	dao := GroupsDao{ds: &ds}
	gr, err := dao.ReadGroup(gId)
	if err != nil {
		respondWith500(w, err.Error())
		return
	}
	gr.GName = name
	_, err = dao.UpdateGroup(&gr)
	if err != nil {
		respondWith500(w, err.Error())
		return
	}
}

func groupDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gid, ok := vars["g_id"]
	if !ok {
		respondWithBadURI(w, r)
		return
	}
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		respondWithUriError(w, r, err)
		return
	}
	dao := GroupsDao{ds: &ds}
	_, err = dao.DeleteGroup(gId)
	if err != nil {
		respondWith500(w, err.Error())
	}
}

func returnGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gid, ok := vars["g_id"]
	if !ok {
		respondWithBadURI(w, r)
		return
	}
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		respondWithUriError(w, r, err)
		return
	}
	dao := GroupsDao{ds: &ds}
	currGroup, err := dao.ReadGroup(gId)
	if err != nil {
		respondWith500(w, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(currGroup)
	if err != nil {
		respondWith500(w, err.Error())
	}
}
