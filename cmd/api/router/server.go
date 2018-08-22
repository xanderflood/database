package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/xanderflood/database/lib/web"
	"github.com/xanderflood/database/pkg/dbi"
)

type Server struct {
	DB   dbi.Interface
	Vars web.VarsGetter
}

type createTableRequest struct {
	Name string `json:name`
}

func (server *Server) CreateTable(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		web.JSONStandardRespond(
			w,
			fmt.Sprintf("failed to read request body: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	request := createTableRequest{}
	err = json.Unmarshal(buf.Bytes(), &request)
	if err != nil {
		web.JSONStandardRespond(
			w,
			fmt.Sprintf("failed to unmarshal request body: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	err = server.DB.CreateTable(request.Name)
	if err != nil {
		web.JSONStandardRespond(
			w,
			fmt.Sprintf("failed to create table: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	web.JSONStandardRespond(
		w,
		fmt.Sprintf("successfully created table: %s", request.Name),
		http.StatusOK,
	)
}

func (server *Server) Index(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (server *Server) Insert(w http.ResponseWriter, r *http.Request) {
	vars := server.Vars.Get(r)
	table := strings.TrimSpace(vars["table"])
	if len(table) == 0 {
		web.JSONStandardRespond(
			w,
			fmt.Sprintf("invalid table name: %s", table),
			http.StatusBadRequest,
		)
	}

	buf := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		web.JSONStandardRespond(
			w,
			fmt.Sprintf("failed to read request body: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	payload := map[string]string{}
	err = json.Unmarshal(buf.Bytes(), &payload)
	if err != nil {
		web.JSONStandardRespond(
			w,
			fmt.Sprintf("failed to unmarshal request body: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	err = server.DB.Insert(table, payload)
	if err != nil {
		web.JSONStandardRespond(
			w,
			fmt.Sprintf("failed to insert: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	web.JSONStandardRespond(
		w,
		fmt.Sprintf("successfully inserted into table: %s", table),
		http.StatusOK,
	)
}
