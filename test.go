package main

import(
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
)

func TestGetMeetings(t *testing.T) {
	req, err := http.NewRequest("GET", "/meetings", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(meetings)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[
		{
			"title": "Rreifdf",
			"start": "Sat Sep 23 2017 15:38:22 GMT+0630",
			"end": "Sat Sep 23 2017 17:38:22 GMT+0630",
			"participants": [
				{
					"name": "Ritik",
					"email": "ritik.gupta2018@vitstudent.ac.in",
					"rsvp": "YES"
				},
				{
					"name": "Ritik",
					"email": "ritik.gupta2018@vitstudent.ac.in",
					"rsvp": "YES"
				}
			]
		},
		{
			"title": "Rreifdf",
			"start": "Sat Sep 23 2017 04:38:22 GMT+0630",
			"end": "Sat Sep 23 2017 08:38:22 GMT+0630",
			"participants": [
				{
					"name": "Ritik",
					"email": "ritik.gupta2018@vitstudent.ac.in",
					"rsvp": "YES"
				},
				{
					"name": "Ritik",
					"email": "ritik.gupta2018@vitstudent.ac.in",
					"rsvp": "YES"
				}
			]
		},
		{
			"title": "Rreifdf",
			"start": "1Sat Sep 23 2017 16:38:22 GMT+0630",
			"end": "Sat Sep 23 2017 16:38:22 GMT+0630",
			"participants": [
				{
					"name": "Ritik",
					"email": "ritik.gupta2018@vitstudent.ac.in",
					"rsvp": "YES"
				},
				{
					"name": "Ritik",
					"email": "ritik.gupta2018@vitstudent.ac.in",
					"rsvp": "YES"
				}
			]
		},
		{
			"title": "Rreifdf",
			"start": "Sat Sep 23 2017 15:38:22 GMT+0630",
			"end": "1Sat Sep 23 2017 17:38:22 GMT+0630",
			"participants": [
				{
					"name": "Ritik",
					"email": "ritik.gupta2018@vitstudent.ac.in",
					"rsvp": "YES"
				},
				{
					"name": "Ritik",
					"email": "ritik.gupta2018@vitstudent.ac.in",
					"rsvp": "YES"
				}
			]
		},
		{
			"title": "Rreifdf",
			"start": "Sat Sep 23 2017 15:38:22 GMT+0630",
			"end": "Sat Sep 23 2017 20:38:22 GMT+0630",
			"participants": [
				{
					"name": "Ritik",
					"email": "ritik.gupta2018@vitstudent.ac.in",
					"rsvp": "YES"
				},
				{
					"name": "Ritik",
					"email": "ritik.gupta2018@vitstudent.ac.in",
					"rsvp": "YES"
				}
			]
		}
	]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetMeetingByID(t *testing.T) {

	req, err := http.NewRequest("GET", "/meetings", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("5f8d430cc7f8bdbfafbb99b6")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getmeeting)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"title":"Rreifdf",
	"start":"Sat Sep 23 2017 15:38:22 GMT+0630",
	"end":"Sat Sep 23 2017 20:38:22 GMT+0630",
	"participants":[{"name":"Ritik","email":"vbv","rsvp":"YES"},
	{"name":"Ritik","email":"ritik.gupta2018@nb.ac.in","rsvp":"YES"}]}
	`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}


func TestCreateMeeting(t *testing.T) {

	var jsonStr = []byte(`{
    
        "title":"Rreifdf",
        "end":"Sat Sep 23 2017 20:38:22 GMT+0630",
        "start":"Sat Sep 23 2017 15:38:22 GMT+0630",
        "participants":[{
            "name":"Ritik",
            "email":"vbv",
            "rsvp":"YES"
        },{
            "name":"Ritik",
            "email":"ritik.gupta2018@nb.ac.in",
            "rsvp":"YES"
        }
        ]


}`)

	req, err := http.NewRequest("POST", "/meetings", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(meetings)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{
		"InsertedID": "5f8d4980d827390ce5f28a47"
	  }`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}