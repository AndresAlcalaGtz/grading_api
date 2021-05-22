package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var grades = make(map[string]map[string]float32)

type Score struct {
	Student string
	Subject string
	Grade   float32
}

func main() {
	http.HandleFunc("/grading", grade)
	http.HandleFunc("/grading/", grade_student)

	fmt.Println("RESTful API is running...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func grade(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		var score Score
		err := json.NewDecoder(request.Body).Decode(&score)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		json_response := create_grade(score)
		response.Header().Set("Content-Type", "application/json")
		response.Write(json_response)

	case "GET":
		json_response, err := read_grades()
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		response.Write(json_response)

	case "PUT":
		var score Score
		err := json.NewDecoder(request.Body).Decode(&score)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		json_response := update_grade(score)
		response.Header().Set("Content-Type", "application/json")
		response.Write(json_response)
	}
}

func grade_student(response http.ResponseWriter, request *http.Request) {
	student := strings.TrimPrefix(request.URL.Path, "/grading/")

	switch request.Method {
	case "GET":
		json_response, err := read_grade(student)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
		}
		response.Header().Set("Content-Type", "application/json")
		response.Write(json_response)

	case "DELETE":
		json_response := delete_student(student)
		response.Header().Set("Content-Type", "application/json")
		response.Write(json_response)
	}
}

func create_grade(score Score) []byte {
	_, exists := grades[score.Student][score.Subject]
	if exists {
		return []byte(`{"code": "ERROR: the grade ​has already been registered"}`)
	}

	_, registered := grades[score.Student]
	if !registered {
		grades[score.Student] = map[string]float32{score.Subject: score.Grade}
	} else {
		grades[score.Student][score.Subject] = score.Grade
	}
	return []byte(`{"code": "SUCCESS: the grade ​has been registered"}`)
}

func read_grades() ([]byte, error) {
	json_data, err := json.MarshalIndent(grades, "", "    ")
	if err != nil {
		return json_data, nil
	}
	return json_data, err
}

func read_grade(student string) ([]byte, error) {
	json_data := []byte(`{}`)

	stu, registered := grades[student]
	if !registered {
		return json_data, nil
	}

	json_data, err := json.MarshalIndent(stu, "", "    ")
	if err != nil {
		return json_data, err
	}
	return json_data, nil
}

func update_grade(score Score) []byte {
	_, registered := grades[score.Student][score.Subject]
	if !registered {
		return []byte(`{"code": "ERROR: the student or the subject has not been registered"}`)
	}

	grades[score.Student][score.Subject] = score.Grade
	return []byte(`{"code": "SUCCESS: the grade ​has been updated"}`)
}

func delete_student(student string) []byte {
	_, registered := grades[student]
	if !registered {
		return []byte(`{"code": "ERROR: the student has not been registered"}`)
	}

	delete(grades, student)
	return []byte(`{"code": "SUCCESS: the student has been deleted"}`)
}
