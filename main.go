package main

import (
	"bitresume/config"
	"os"
	// "bitresume/jobs"
	"bitresume/routes"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	config.InitOAuth()
	config.InitDB()
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Hello from Render!"})
    })
	r.Static("/uploads", "./uploads")
	corsConfig := cors.Config{
		AllowOrigins:     []string{"https://student-smart-hub.web.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}	
	r.Use(cors.New(corsConfig))
	routes.RegisterRoutes(r)
	c := cron.New(cron.WithSeconds())
	// _, errCron := c.AddFunc("0 37 11 * * *", jobs.CallDailyTasksForAllDates)
	// Schedule the job to run every day at 11:50 AM(seconds minute hour dayOfMonth month dayOfWeek)		
	// if errCron != nil {
	// 	panic("Failed to schedule cron job: " + errCron.Error())
	// }
	c.Start()
	r.Run(":" + os.Getenv("PORT"))
}