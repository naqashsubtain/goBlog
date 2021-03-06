package controllers

import "github.com/naqash/goBlog/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", //middlewares.SetMiddlewareAuthentication(
		middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route for admin users
	s.Router.HandleFunc("/admin/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	s.Router.HandleFunc("/createInvite", middlewares.SetMiddlewareAuthentication(s.CreateInviteToken))
	s.Router.HandleFunc("/invite/login", s.SetMiddlewareInviteToken(s.GetJobs))

	//Jobs routes
	s.Router.HandleFunc("/jobs", middlewares.SetMiddlewareJSON(s.CreateJob)).Methods("POST")
	s.Router.HandleFunc("/jobs", middlewares.SetMiddlewareJSON(s.GetJobs)).Methods("GET")
	s.Router.HandleFunc("/jobs/{id}", middlewares.SetMiddlewareJSON(s.GetJob)).Methods("GET")
	s.Router.HandleFunc("/jobs/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateJob))).Methods("PUT")
	s.Router.HandleFunc("/jobs/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteJob)).Methods("DELETE")
	s.Router.HandleFunc("/softDeletejob/{id}", middlewares.SetMiddlewareAuthentication(s.SoftDeleteJob)).Methods("DELETE")
	s.Router.HandleFunc("/jobsByDistance", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetJobsbyDistance))).Methods("POST")
}
