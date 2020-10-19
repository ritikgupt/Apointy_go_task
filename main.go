package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"strconv"
	"time"
	"go.mongodb.org/mongo-driver/bson"
 	"go.mongodb.org/mongo-driver/bson/primitive"
 	"go.mongodb.org/mongo-driver/mongo"
 	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

 // this is the code to create the meeting structure
type Meeting struct {
	Title        string `json:"title"`
	Start_time string `json:"start"`
	End_time          string `json:"end"`
	creation       string `json:"creation"`
	Participants       []Participant    `json:"participants"`
}

// this is the code to create the participant structure
type Participant struct {
	Name string `json:"name"`
	Email string `json:"email"`
	RSVP string `json:"rsvp"`
}

//this function takes the request from /meetings route,identifies the request method accoringly i.e whether it is post or get request and redirects accordingly
func meetings(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	fmt.Println(r.Method)
	switch r.Method {
	case "POST":
		post(w,r)
		return
	case "GET":
		getspecificmeeting(w,r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

//this fuction takes the request from the /meeting/ route,identifies the request method accordingly and redirects accordingly

func getmeeting(w http.ResponseWriter,r *http.Request){
	fmt.Println(r)
	fmt.Println(r.Method)
	switch r.Method {
	case "GET":
		get(w,r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

//this fuction takes the request from the /meeting route,identifies the request method accordingly and redirects accordingly
func allmeetings(w http.ResponseWriter,r *http.Request){
	switch r.Method {
	case "GET":
		getparticipants(w,r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

//with the help of this function we can get all the meeting details corresponding to the meeting id

func get(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	fmt.Println(parts)
	fmt.Println(parts[2])
	id, _ := primitive.ObjectIDFromHex(parts[2])
	fmt.Println(id)
	var meeting Meeting
	collection := client.Database("meeting_scheduler").Collection("Meeting")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&meeting)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(meeting)
}

//with the help of this function we can store the meeting details of the participants

func  post(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")
		fmt.Println(r.Body);
		var meeting Meeting
		_ = json.NewDecoder(r.Body).Decode(&meeting)
		x:=meeting.Start_time
		y:=meeting.End_time
		layout := "Mon Jan 02 2006 15:04:05 GMT-0700"
		t1, _ := time.Parse(layout, x)
		t2,_:=time.Parse(layout,y)
		fmt.Println(t1.Format(layout))
		fmt.Println(t2.Format(layout))

		collection := client.Database("meeting_scheduler").Collection("Meeting")
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		result, _ := collection.InsertOne(ctx, meeting)
		json.NewEncoder(w).Encode(result)
}

//with the help of this function we can get the specific meetings accoring to the query parameters and the meeting details will be shown 5 at a time as pagination has been implemented here

func getspecificmeeting(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "application/json")
	start_time := r.FormValue("start")
	end_time := r.FormValue("end")
	page ,err:= strconv.Atoi(r.FormValue("page"))
	a :=5
	x := page*a

	fmt.Println(start_time)
	fmt.Println(end_time)
	var meetings []Meeting
	collection := client.Database("meeting_scheduler").Collection("Meeting")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var meeting Meeting
		cursor.Decode(&meeting)
		meetings = append(meetings, meeting)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	var output []Meeting
	fmt.Println(meetings)
	for i :=x;i<x+5;i++{
		if i>=len(meetings){
			break;
		} else{
		output=append(output,meetings[i])
		}
	}
	json.NewEncoder(w).Encode(output)
}

// with the help of this function we can get the details of the all the meetings of the person with the specific email id

func getparticipants(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "application/json")
	var participants [][]Participant
	var meetings []Meeting
	collection := client.Database("meeting_scheduler").Collection("Meeting")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	pemailid := r.FormValue("participant")
	fmt.Println(pemailid);
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var meeting Meeting
		cursor.Decode(&meeting)
		fmt.Println(meeting.Participants)
		for i := 0; i < len(meeting.Participants); i++ {
			x := meeting.Participants[i]
			if x.Email == pemailid {
				participants = append(participants,meeting.Participants)
				meetings = append(meetings, meeting)
			}
		}
		
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(meetings)
}



func main() {
	fmt.Println("Server has started")
	// meetingHandlers := newmeetingHandlers()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	http.HandleFunc("/meetings", meetings)
	http.HandleFunc("/meeting/",getmeeting)
	http.HandleFunc("/meeting",allmeetings)
	// http.HandleFunc("/meetings/", meetingHandlers.getmeeting)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}