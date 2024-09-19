package routers

import (
	"sanco_microservices/controllers"
	"sanco_microservices/middleware"

	"github.com/gin-gonic/gin"
)

var users = map[string]string{
	"admin":  "password",
	"editor": "secret",
}

func StartServer() *gin.Engine {
	router := gin.Default()
	// RoomController := controllers.NewRoomController()

	// roomTypeController := controllers.NewRoomTypeController()
	// FloorController := controllers.NewFloorController()
	GuestController := controllers.NewGuestController()
	SuppliersController := controllers.NewSuppliersController()

	// Public routes
	router.POST("/login", middleware.LoginHandler)
	router.POST("/book", GuestController.Book)

	router.POST("/api/master/supplier/store", SuppliersController.Create)
	router.GET("/api/master/supplier/detail/:id", SuppliersController.GetByID)
	router.PUT("/api/master/supplier/update/:id", SuppliersController.Update)
	router.DELETE("/api/master/supplier/delete/:id", SuppliersController.Delete)
	router.GET("/api/master/supplier/show", SuppliersController.Get)
	// Protected routes
	protected := router.Group("/protected")
	protected.Use(middleware.AuthMiddleware())
	{

		// protected.GET("/guest/:id", RoomController.GetByID)
		// protected.PUT("/checkIn/:id", GuestController.CheckIn)
		// protected.PUT("/checkOut/:id", GuestController.CheckOut)
		// protected.GET("/unbookedRom", RoomController.GetUnBooked)

		// protected.POST("/room", RoomController.Create)
		// protected.GET("/room/:id", RoomController.GetByID)
		// protected.PUT("/room/:id", RoomController.Update)
		// protected.DELETE("/room/:id", RoomController.Delete)
		// protected.GET("/room", RoomController.Get)

		// protected.POST("/roomtypes", roomTypeController.CreateRoomType)
		// protected.GET("/roomtypes/:id", roomTypeController.GetRoomTypeByID)
		// protected.PUT("/roomtypes/:id", roomTypeController.UpdateRoomType)
		// protected.DELETE("/roomtypes/:id", roomTypeController.DeleteRoomType)
		// protected.GET("/roomtypes", roomTypeController.GetAllRoomTypes)

		// protected.POST("/floors", FloorController.CreateFloor)
		// protected.GET("/floors/:id", FloorController.GetFloorByID)
		// protected.PUT("/floors/:id", FloorController.UpdateFloor)
		// protected.DELETE("/floors/:id", FloorController.DeleteFloor)
		// protected.GET("/floors", FloorController.GetAllFloors)

		// protected.POST("/api/master/supplier/store", SuppliersController.Create)
		// protected.GET("/api/master/supplier/detail/:id", SuppliersController.GetByID)
		// protected.PUT("/api/master/supplier/update/:id", SuppliersController.Update)
		// protected.DELETE("/room/:id", SuppliersController.Delete)
		// protected.GET("/api/master/supplier/show", SuppliersController.Get)

	}

	return router
}
