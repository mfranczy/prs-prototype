package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// TODO: move dml to procedures

type Service struct {
	DB *sqlx.DB // is it a good idea??
}

type respFunc func(http.ResponseWriter, *http.Request) (interface{}, int, error)

func respHandler(fn respFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		data, status, err := fn(rw, req)
		if err != nil {
			log.Println("error", err)
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(status)

		err = json.NewEncoder(rw).Encode(data)
		if err != nil {
			log.Println("could not encode response:", err)
		}
	}
}

type EndpointEntry struct {
	Method      string
	Endpoint    string
	HandlerFunc http.HandlerFunc
}

func getEndpoints(s Service) []EndpointEntry {
	return []EndpointEntry{
		// packages
		{
			Method:      "GET",
			Endpoint:    "/packages/list",
			HandlerFunc: respHandler(s.listPackages),
		},
		// labels
		{
			Method:      "GET",
			Endpoint:    "/labels/list",
			HandlerFunc: respHandler(s.listLabels),
		},
		{
			Method:      "GET",
			Endpoint:    "/labels/list-packages/{id}",
			HandlerFunc: respHandler(s.listLabelPackages),
		},
		{
			Method:      "GET",
			Endpoint:    "/labels/get/{id}",
			HandlerFunc: respHandler(s.getLabel),
		},
		{
			Method:      "POST",
			Endpoint:    "/labels/add",
			HandlerFunc: respHandler(s.addLabel),
		},
		{
			Method:      "POST",
			Endpoint:    "/labels/attach-pkg",
			HandlerFunc: respHandler(s.attachLabel),
		},
	}
}

func Init(r *mux.Router, s Service) {
	subR := r.PathPrefix("/api/v1").Subrouter()

	for _, e := range getEndpoints(s) {
		subR.HandleFunc(e.Endpoint, e.HandlerFunc).Methods(e.Method)
	}
}

// translate to javascript language !! MOVE IT TO DB PACKAGE !!
type NullString sql.NullString

func (ns *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullString{s.String, false}
	} else {
		*ns = NullString{s.String, true}
	}

	return nil
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}
