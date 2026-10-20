package main

import (
	"log"
	"os"

	"github.com/HisakeyT/backpacker-platform/internal/auth"
	"github.com/HisakeyT/backpacker-platform/internal/database"
	"github.com/HisakeyT/backpacker-platform/internal/handler"
	"github.com/HisakeyT/backpacker-platform/internal/repository/gorm"
	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtManager := auth.NewJWTManager(jwtSecret)

	userRepository := gormrepo.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, jwtManager)
	userHandler := handler.NewHandler(userUsecase)

	travelRepository := gormrepo.NewTravelRepository(db)
	travelUsecase := usecase.NewTravelUsecase(travelRepository)
	travelHandler := handler.NewTravelHandler(travelUsecase)

	travelPlanRepository := gormrepo.NewGormPlanRepository(db)
	travelPlanUsecase := usecase.NewTravelPlanUsecase(travelRepository, travelPlanRepository)
	travelPlanHandler := handler.NewTravelPlanHandler(travelPlanUsecase)
	planRecordRepository := gormrepo.NewPlanRecordRepository(db)
	planRecordUsecase := usecase.NewPlanRecordUsecase(travelRepository, travelPlanRepository, planRecordRepository)
	planRecordHandler := handler.NewPlanRecordHandler(planRecordUsecase)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	protected := r.Group("")
	protected.Use(auth.AuthMiddleware(jwtManager, userRepository))

	protected.GET("/users/me", userHandler.Me)
	protected.PATCH("/users/me", userHandler.UpdateMe)

	protected.POST("/travels", travelHandler.CreateTravel)
	protected.PATCH("/travels/:travel_id", travelHandler.UpdateTravel)
	protected.POST("/travels/:travel_id/copy", travelHandler.CopyTravel)
	protected.GET("/travels", travelHandler.GetTravels)
	protected.GET("/travels/:travel_id", travelHandler.GetTravel)
	protected.DELETE("/travels/:travel_id", travelHandler.DeleteTravel)
	r.GET("/travels/public", travelHandler.GetPublicTravels)
	r.GET("/travels/public/:travel_id", travelHandler.GetPublicTravel)

	protected.POST("/travels/:travel_id/plans", travelPlanHandler.CreateTravelPlan)
	protected.PUT("/travels/:travel_id/plans/order", travelPlanHandler.ReorderTravelPlans)
	protected.PATCH("/travels/:travel_id/plans/:travel_plan_id", travelPlanHandler.UpdateTravelPlan)
	protected.GET("/travels/:travel_id/plans", travelPlanHandler.GetTravelPlans)
	protected.DELETE("/travels/:travel_id/plans/:travel_plan_id", travelPlanHandler.DeleteTravelPlan)
	r.GET("/travels/public/:travel_id/plans", travelPlanHandler.GetPublicTravelPlan)

	protected.GET("/travels/:travel_id/plans/:travel_plan_id/records/:kind", planRecordHandler.List)
	protected.POST("/travels/:travel_id/plans/:travel_plan_id/records/:kind", planRecordHandler.Create)
	protected.PATCH("/travels/:travel_id/plans/:travel_plan_id/records/:kind/:record_id", planRecordHandler.Update)
	protected.DELETE("/travels/:travel_id/plans/:travel_plan_id/records/:kind/:record_id", planRecordHandler.Delete)
	r.GET("/travels/public/:travel_id/plans/:travel_plan_id/records/:kind", planRecordHandler.ListPublic)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
