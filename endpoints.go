package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/syedomair/gokit_interface/models"
)

type Endpoints struct {
	PostUserEndpoint         endpoint.Endpoint
	PostAuthenticateEndpoint endpoint.Endpoint
	GetUserEndpoint          endpoint.Endpoint
	PutUserEndpoint          endpoint.Endpoint
	PatchBookEndpoint        endpoint.Endpoint
	GetUserBookEndpoint      endpoint.Endpoint
	PostBookEndpoint         endpoint.Endpoint
	GetBooksEndpoint         endpoint.Endpoint
	GetPublicBooksEndpoint   endpoint.Endpoint
	GetBookEndpoint          endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostUserEndpoint:         MakePostUserEndpoint(s),
		PostAuthenticateEndpoint: MakePostAuthenticateEndpoint(s),
		GetUserEndpoint:          MakeGetUserEndpoint(s),
		PutUserEndpoint:          MakePutUserEndpoint(s),
		PatchBookEndpoint:        MakePatchBookEndpoint(s),
		GetUserBookEndpoint:      MakeGetUserBookEndpoint(s),
		PostBookEndpoint:         MakePostBookEndpoint(s),
		GetBooksEndpoint:         MakeGetBooksEndpoint(s),
		GetPublicBooksEndpoint:   MakeGetPublicBooksEndpoint(s),
		GetBookEndpoint:          MakeGetBookEndpoint(s),
	}
}

func MakePostUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postUserRequest)
		e := s.PostUser(ctx, req.User)
		return successResponse(e), e
	}
}
func MakePostAuthenticateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postUserRequest)
		u, e := s.PostAuthenticate(ctx, req.User)
		return successResponse(u), e
	}
}
func MakeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		u, e := s.GetUser(ctx, req.ID)
		return successResponse(u), e
	}
}

func MakePutUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putUserRequest)
		e := s.PutUser(ctx, req.ID, req.User)
		return successResponse(e), e
	}
}

func MakePatchBookEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(patchBookRequest)
		e := s.PatchBook(ctx, req.ID, req.Book)
		return successResponse(e), e
	}
}

func MakeGetUserBookEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		books, offset, limit, count := s.GetUserBooks(ctx, req.ID, req.Offset, req.Limit, req.Orderby, req.Sort)
		return successResponseList(books, offset, limit, count), nil
	}
}

func MakePostBookEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postBookRequest)
		str, e := s.PostBook(ctx, req.Book)
		if str != "" {
			return errorResponse(str), nil
		} else {
			return successResponse(e), e
		}
	}
}
func MakeGetBooksEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		books, offset, limit, count := s.GetBooks(ctx, req.Offset, req.Limit, req.Orderby, req.Sort)
		return successResponseList(books, offset, limit, count), nil
	}
}

func MakeGetPublicBooksEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		books, offset, limit, count := s.GetPublicBooks(ctx, req.Offset, req.Limit, req.Orderby, req.Sort)
		return successResponseList(books, offset, limit, count), nil
	}
}
func MakeGetBookEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		b, e := s.GetBook(ctx, req.ID)
		return successResponse(b), e
	}
}

type postUserRequest struct {
	User models.User
}
type postBookRequest struct {
	Book models.Book
}
type putUserRequest struct {
	ID   string
	User models.User
}
type patchBookRequest struct {
	ID   string
	Book models.Book
}

type getRequest struct {
	ID      string
	Offset  string
	Limit   string
	Orderby string
	Sort    string
}

func errorAuthResponse(class interface{}) map[string]interface{} {
	return commonResponse(class, "error", "400")
}

func errorResponse(class interface{}) map[string]interface{} {
	return commonResponse(class, "error", "500")
}

func successResponse(class interface{}) map[string]interface{} {
	return commonResponse(class, "success", "200")
}

func commonResponse(class interface{}, result string, code string) map[string]interface{} {
	response := make(map[string]interface{})
	response["data"] = class
	response["result"] = result
	response["code"] = code
	return response
}

func successResponseList(class interface{}, offset string, limit string, count string) map[string]interface{} {
	tempResponse := make(map[string]interface{})
	tempResponse["offset"] = offset
	tempResponse["limit"] = limit
	tempResponse["count"] = count
	tempResponse["list"] = class
	return successResponse(tempResponse)
}
