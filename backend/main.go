package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"forum/db"
	"forum/handlers"
	"forum/utils/middlewares"
)

func main() {
	// Initialize database
	if err := db.InitDB("../database/forum.db"); err != nil {
		log.Fatal("Database initialization failed:", err)
	}
	// to close db when panic
	defer func() {
		if err := recover(); err != nil {
			db.DB.Close()
			log.Fatal("Error: ", err)
		}
	}()

	// to close db when ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		db.DB.Close()
		fmt.Println()
		os.Exit(0)
	}()

	// Apply migrations
	if err := db.RunMigrations(); err != nil {
		panic("Migrations failed:" + err.Error())
	}

	// Create a rate limiter
	rl := middlewares.NewRateLimiter()

	// Start a goroutine to clean up old entries every second
	go rl.CleanupOldEntries()

	// Configure routes
	http.Handle("/check-auth", middlewares.RateLimit(rl, http.HandlerFunc(handlers.CheckAuthHandler)))
	http.Handle("/register", middlewares.RateLimit(rl, middlewares.ForbidnMiddleware(http.HandlerFunc(handlers.RegisterHandler))))
	http.Handle("/login", middlewares.RateLimit(rl, middlewares.ForbidnMiddleware(http.HandlerFunc(handlers.LoginHandler))))
	http.Handle("/logout", middlewares.RateLimit(rl, middlewares.AuthMiddleware(http.HandlerFunc(handlers.LogoutHandler))))
	http.Handle("/posts", middlewares.RateLimit(rl, middlewares.AuthMiddleware(http.HandlerFunc(handlers.GetPostsHandler))))
	http.Handle("/posts/store", middlewares.RateLimit(rl, middlewares.AuthMiddleware(http.HandlerFunc(handlers.CreatePostHandler))))
	http.Handle("/comments", middlewares.RateLimit(rl, middlewares.AuthMiddleware(http.HandlerFunc(handlers.GetCommentsHandler))))
	http.Handle("/comments/store", middlewares.RateLimit(rl, middlewares.AuthMiddleware(http.HandlerFunc(handlers.CreateCommentHandler))))
	http.Handle("/likes", middlewares.RateLimit(rl, middlewares.AuthMiddleware(http.HandlerFunc(handlers.GetLikesHandler))))
	http.Handle("/likes/store", middlewares.RateLimit(rl, middlewares.AuthMiddleware(http.HandlerFunc(handlers.CreateLikesHandler))))
	http.Handle("/messages", middlewares.RateLimit(rl, middlewares.AuthMiddleware(http.HandlerFunc(handlers.GetMessageHistoryHandler))))
	http.Handle("/ws", middlewares.RateLimit(rl, middlewares.AuthMiddleware(http.HandlerFunc(handlers.WebSocketHandler))))

	// Configure routes
	http.Handle("/", http.HandlerFunc(handlers.ServeFilesHandler))

	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("Server failed to start:" + err.Error())
	}
}
