package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/syedomair/kit2/models"
)

type Endpoints struct {
	PostUserEndpoint    endpoint.Endpoint
	GetUserEndpoint     endpoint.Endpoint
	PutUserEndpoint     endpoint.Endpoint
	GetUserBookEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostUserEndpoint:    MakePostUserEndpoint(s),
		GetUserEndpoint:     MakeGetUserEndpoint(s),
		PutUserEndpoint:     MakePutUserEndpoint(s),
		GetUserBookEndpoint: MakeGetUserBookEndpoint(s),
	}
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

func MakePostUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postUserRequest)
		e := s.PostUser(ctx, req.User)
		//return postUserResponse{Err: e}, nil
		return successResponse(e), e
	}
}
func MakeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		u, e := s.GetUser(ctx, req.ID)
		return successResponse(u), e
	}
}

func MakeGetUserBookEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		b, o, l, c := s.GetUserBooks(ctx, req.ID, req.Offset, req.Limit, req.Orderby, req.Sort)
		//return bookResponse, offset, limit, count
		return successResponseList(b, o, l, c), nil
	}
}

func MakePutUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putUserRequest)
		e := s.PutUser(ctx, req.ID, req.User)
		return successResponse(e), e
		//return putUserResponse{Err: e}, nil
	}
}

type postUserRequest struct {
	User models.User
}
type getRequest struct {
	ID      string
	Offset  string
	Limit   string
	Orderby string
	Sort    string
}

type putUserRequest struct {
	ID   string
	User models.User
}
