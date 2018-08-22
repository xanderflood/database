package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/xanderflood/database/lib/web"
	"github.com/xanderflood/database/pkg/dbi"
)

type Server struct {
	DB dbi.Interface
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
	//TODO
}
