package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/syedomair/gokit_interface/models"
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)

	// POST    /Book
	// POST    /authenticate
	// GET     /books
	// GET     /public/books
	// GET     /book/:id
	// GET     /user/:id
	// GET     /my-books/:id
	// POST    /user
	// PUT     /user/:id
	// PATCH   /book/:id

	r.Methods("POST").Path("/authenticate").Handler(httptransport.NewServer(
		e.PostAuthenticateEndpoint,
		decodePostUserRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/book").Handler(httptransport.NewServer(
		e.PostBookEndpoint,
		decodePostBookRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/books").Handler(httptransport.NewServer(
		e.GetBooksEndpoint,
		decodeGetListRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/public/books").Handler(httptransport.NewServer(
		e.GetPublicBooksEndpoint,
		decodeGetListRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/book/{id}").Handler(httptransport.NewServer(
		e.GetBookEndpoint,
		decodeGetRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/my-books/{id}").Handler(httptransport.NewServer(
		e.GetUserBookEndpoint,
		decodeGetListIdRequest,
		encodeResponse,
	))
	r.Methods("PUT").Path("/user/{id}").Handler(httptransport.NewServer(
		e.PutUserEndpoint,
		decodePutUserRequest,
		encodeResponse,
	))
	r.Methods("PATCH").Path("/book/{id}").Handler(httptransport.NewServer(
		e.PatchBookEndpoint,
		decodePatchBookRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/user").Handler(httptransport.NewServer(
		e.PostUserEndpoint,
		decodePostUserRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/user/{id}").Handler(httptransport.NewServer(
		e.GetUserEndpoint,
		decodeGetRequest,
		encodeResponse,
	))
	return r
}

func decodePostUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postUserRequest
	if e := json.NewDecoder(r.Body).Decode(&req.User); e != nil {
		return nil, e
	}
	return req, nil
}

func decodePostBookRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postBookRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Book); e != nil {
		return nil, e
	}
	return req, nil
}
func decodeGetRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, nil
	}
	return getRequest{ID: id}, nil
}

func decodePutUserRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, nil
	}
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, err
	}
	return putUserRequest{
		ID:   id,
		User: user,
	}, nil
}

func decodePatchBookRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, nil
	}
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		return nil, err
	}
	return patchBookRequest{
		ID:   id,
		Book: book,
	}, nil
}

func decodeGetListIdRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, nil
	}
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	orderby := r.URL.Query().Get("orderby")
	sort := r.URL.Query().Get("sort")

	if offset == "" {
		offset = "0"
	}
	if limit == "" {
		limit = "10000"
	}
	if orderby == "" {
		orderby = "id"
	}
	if sort == "" {
		sort = "asc"
	}

	return getRequest{ID: id, Offset: offset, Limit: limit, Orderby: orderby, Sort: sort}, nil
}

func decodeGetListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	orderby := r.URL.Query().Get("orderby")
	sort := r.URL.Query().Get("sort")

	if offset == "" {
		offset = "0"
	}
	if limit == "" {
		limit = "10000"
	}
	if orderby == "" {
		orderby = "id"
	}
	if sort == "" {
		sort = "asc"
	}

	return getRequest{Offset: offset, Limit: limit, Orderby: orderby, Sort: sort}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}
