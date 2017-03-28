package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

type testCaseType struct {
	method         string
	url            string
	requestBody    string
	responseResult string
	responseData   string
}

var testCases []testCaseType
var sh http.Handler

func init() {

	var logger log.Logger
	logger = log.NewNopLogger()
	var s Service
	{
		s = DBService()
	}
	var h http.Handler
	{
		h = MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}
	sh = securityMiddleware(s, h)
}

func TestServer(t *testing.T) {

	srv := httptest.NewServer(sh)
	defer srv.Close()

	testCases = []testCaseType{
		{"GET", srv.URL + "/user/1", ``, `"success"`, `{"email":"john@gmail.com","first_name":"John","id":1,"last_name":"Smith"}`},
		{"POST", srv.URL + "/book", `{"name":"Test Book changed","description":"Desc","publish":true, "user_id":1}`, `"success"`, `null`},
		{"PATCH", srv.URL + "/book/1", `{"name":"Test Book changed","description":"Desc changed","publish":true, "user_id":1}`, `"success"`, `null`},
		{"GET", srv.URL + "/book/1", ``, `"success"`, `{"book_name":"Test Book changed","description":"Desc changed","first_name":"John","id":1,"last_name":"Smith","publish":true,"user_id":1}`},
		{"GET", srv.URL + "/my-books/1", ``, `"success"`, `list`},
		{"GET", srv.URL + "/books", ``, `"success"`, `list`},
		{"GET", srv.URL + "/public/books", ``, `"success"`, `list`},
		{"POST", srv.URL + "/book", `{"name":"","description":"Desc","publish":true, "user_id":1}`, `"error"`, `"name is a requird field"`},
		{"POST", srv.URL + "/book", `{"name":"Test Book","description":"Desc","publish":true}`, `"error"`, `"user_id is a requird field"`},
	}

	commonTest(testCases, t, "dHb%e@Bg0f8-API_KEY-&bE71jKoH=2", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImpvaG5AZ21haWwuY29tIiwicGFzc3dvcmQiOiJNVEl6TkRVPSIsImlzcyI6InRlc3QifQ.Fk0pm1AmEl7nNGH_7Xcs93r5V2nhPGnBCfWbIwkhTHk")
}

func TestServerWithoutHeader(t *testing.T) {

	srv := httptest.NewServer(sh)
	defer srv.Close()

	testCases = []testCaseType{
		{"GET", srv.URL + "/user/1", ``, `"error"`, `"Header missing: x-key "`},
	}

	commonTest(testCases, t, "", "")
}
func TestServerInvalidHeader(t *testing.T) {

	srv := httptest.NewServer(sh)
	defer srv.Close()

	testCases = []testCaseType{
		{"GET", srv.URL + "/user/1", ``, `"error"`, `"Invalid JWT Signature"`},
	}

	commonTest(testCases, t, "dHb%e@Bg0f8-API_KEY-&bE71jKoH=2", "TeyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImpvaG5AZ21haWwuY29tIiwicGFzc3dvcmQiOiJNVEl6TkRVPSIsImlzcyI6InRlc3QifQ.Fk0pm1AmEl7nNGH_7Xcs93r5V2nhPGnBCfWbIwkhTHk")
}

func commonTest(testCases []testCaseType, t *testing.T, xkey string, xjwt string) {

	i := 0
	for _, testCase := range testCases {
		req, _ := http.NewRequest(testCase.method, testCase.url, strings.NewReader(testCase.requestBody))
		if xkey != "" {
			req.Header.Set("x-key", xkey)
		}
		if xjwt != "" {
			req.Header.Set("x-jwt", xjwt)
		}
		resp, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)

		var bodyInterface map[string]interface{}
		json.Unmarshal(body, &bodyInterface)
		jsonData, _ := json.Marshal(bodyInterface["data"])
		jsonResult, _ := json.Marshal(bodyInterface["result"])
		fmt.Println(strconv.Itoa(i) + " " + testCase.method + " " + string(testCase.url))
		//fmt.Println(string(jsonData))
		//fmt.Println(string(jsonResult))

		if string(jsonData) != testCase.responseData {
			if testCase.responseData != "list" {
				t.Error("Expected:" + testCase.responseData + " got:" + string(jsonData))
			}
		}
		if string(jsonResult) != testCase.responseResult {
			t.Error("Expected:" + testCase.responseResult + " got:" + string(jsonResult))
		}
		i++
	}
}
