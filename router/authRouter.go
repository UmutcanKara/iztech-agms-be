package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	throttle "github.com/s12i/gin-throttle"
	"iztech-agms/db"
	"iztech-agms/internal/auth"
	"iztech-agms/router/middleware"
	"time"
)

func AuthRouter(blackListHosts map[string]struct{}) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	conn, err := db.NewDatabase()
	if err != nil {
		panic(err)
	}
	maxEventPerSec := 1000
	maxBurstSize := 20

	r.Use(throttle.Throttle(maxEventPerSec, maxBurstSize))

	// Configure CORS middleware
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://172.17.0.2:5173"}, // Replace with your allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Apply the CORS middleware to the router
	r.Use(cors.New(corsConfig))

	r.Use(throttle.Throttle(maxEventPerSec, maxBurstSize))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Auth Service Pong!"})
	})
	r.Use(middleware.Security(blackListHosts))

	repo := auth.NewRepository(conn.GetDB())
	service := auth.NewService(repo)
	handler := auth.NewHandler(service)

	publicGroup := r.Group("/")
	{
		publicGroup.POST("/login", handler.Login)
		publicGroup.POST("/register", handler.Register)
		publicGroup.POST("/create", handler.CreateUsers)
	}
	protectedGroup := r.Group("/")
	jwtMiddleware := middleware.NewJWTMiddleware()
	{
		protectedGroup.Use(jwtMiddleware.Authorize())
		protectedGroup.GET("/changePwd", handler.ChangePwd)
		protectedGroup.POST("/logout", handler.Logout)
	}

	return r
}
