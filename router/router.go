package router

import (
	"github.com/gorilla/mux"
	"github.com/mangesh-shinde/learnera/handlers"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HomeHandler)

	//Course routes
	r.HandleFunc("/api/courses", handlers.GetAllCourses).Methods("GET")
	r.HandleFunc("/api/courses/{courseid}", handlers.GetCourseById).Methods("GET")
	r.HandleFunc("/api/courses/{courseid}", handlers.DeleteCourseById).Methods("DELETE")
	r.HandleFunc("/api/courses", handlers.AddCourse).Methods("POST")
	r.HandleFunc("/api/courses", handlers.DeleteAllCourses).Methods("DELETE")

	//User routes
	r.HandleFunc("/api/login", handlers.UserLogin).Methods("POST")
	r.HandleFunc("/api/signup", handlers.AddUser).Methods("POST")
	r.HandleFunc("/api/profiles", handlers.GetAllUserProfiles).Methods("GET")
	r.HandleFunc("/api/profiles/{username}", handlers.GetUserProfile).Methods("GET")
	r.HandleFunc("/api/profiles", handlers.AddUserProfile).Methods("POST")
	r.HandleFunc("/api/courses/add-to-cart/{courseid}", handlers.AddCourseToCart).Methods("POST")

	return r
}
