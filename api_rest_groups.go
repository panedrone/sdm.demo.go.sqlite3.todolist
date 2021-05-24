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
		respondWithBabRequestError(w, fmt.Sprintf("Invalid request: %s", err.Error()))
		return
	}
	rd := bytes.NewReader(bodyBytes) // === panedrone: r.Body became unavailable
	decoder := json.NewDecoder(rd)
	var inGroup Group
	err = decoder.Decode(&inGroup)
	if err != nil {
		respondWithBabRequestError(w, fmt.Sprintf("JSON decoder FAIL: %s. Input: %s", err.Error(), bodyBytes))
		return
	}
	name := inGroup.GName
	if len(name) == 0 {
		respondWithBabRequestError(w, fmt.Sprintf("Invalid input: %s", bodyBytes))
		return
	}
	dao := GroupsDao{ds: &ds}
	gr := Group{}
	gr.GName = name
	dao.CreateGroup(&gr)
}

func returnAllGroupsHandler(w http.ResponseWriter, _ *http.Request) {
	dao := GroupsDao{ds: &ds}
	groups := dao.GetGroups()
	err := json.NewEncoder(w).Encode(groups)
	if err != nil {
		panic(err)
	}
}

func groupUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gid := vars["g_id"]
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		respondWithBabRequestError(w, "Expected: g_id")
		return
	}
	bodyBytes, err := ioutil.ReadAll(r.Body) // === panedrone: store source for error handling
	if err != nil {
		respondWithBabRequestError(w, fmt.Sprintf("Invalid request: %s", err.Error()))
		return
	}
	rd := bytes.NewReader(bodyBytes) // === panedrone: r.Body became unavailable
	decoder := json.NewDecoder(rd)
	var inGroup Group
	err = decoder.Decode(&inGroup)
	if err != nil {
		respondWithBabRequestError(w, fmt.Sprintf("JSON decoder FAIL: %s. Input: %s", err.Error(), bodyBytes))
		return
	}
	name := inGroup.GName
	if len(name) == 0 {
		respondWithBabRequestError(w, fmt.Sprintf("Invalid input: %s", bodyBytes))
		return
	}
	dao := GroupsDao{ds: &ds}
	gr := dao.ReadGroup(gId)
	gr.GName = name
	dao.UpdateGroup(&gr)
}

func groupDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gid := vars["g_id"]
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		respondWithBabRequestError(w, "Expected: g_id")
		return
	}
	dao := GroupsDao{ds: &ds}
	dao.DeleteGroup(gId)
}

func returnGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gid := vars["g_id"]
	gId, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		respondWithBabRequestError(w, "Expected: g_id")
		return
	}
	dao := GroupsDao{ds: &ds}
	currGroup := dao.ReadGroup(gId)
	err = json.NewEncoder(w).Encode(currGroup)
	if err != nil {
		panic(err)
	}
}
