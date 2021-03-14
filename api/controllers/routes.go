package controllers

import "github.com/agusdwis/rest-api/api/middlewares"

func (server *Server) initializeRoutes() {

	// Home Route
	server.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(server.Home)).Methods("GET")

	// Login Route
	server.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(server.Login)).Methods("POST")

	//Users routes
	server.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(server.CreateUser)).Methods("POST")
	server.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(server.GetUsers)).Methods("GET")
	server.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(server.GetUser)).Methods("GET")
	server.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateUser))).Methods("PUT")
	server.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(server.DeleteUser)).Methods("DELETE")

	//Posts routes
	server.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(server.CreatePost)).Methods("POST")
	server.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(server.GetPosts)).Methods("GET")
	server.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(server.GetPost)).Methods("GET")
	server.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdatePost))).Methods("PUT")
	server.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareAuthentication(server.DeletePost)).Methods("DELETE")

	//Likes route
	server.Router.HandleFunc("/likes/{id}", middlewares.SetMiddlewareJSON(server.GetLike)).Methods("GET")
	server.Router.HandleFunc("/likes/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.LikePost))).Methods("POST")
	server.Router.HandleFunc("/likes/{id}", middlewares.SetMiddlewareAuthentication(server.UnlikePost)).Methods("DELETE")

	//Comments route
	server.Router.HandleFunc("/comments/{id}", middlewares.SetMiddlewareJSON(server.GetComments)).Methods("GET")
	server.Router.HandleFunc("/comments/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.CreateComment))).Methods("POST")
	server.Router.HandleFunc("/comments/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateComment))).Methods("PUT")
	server.Router.HandleFunc("/comments/{id}", middlewares.SetMiddlewareAuthentication(server.DeleteComment)).Methods("DELETE")
}
