package routers

import (
	"sanco_microservices/controllers"
	"sanco_microservices/middleware"
	"sanco_microservices/repository"
	"sanco_microservices/structs"

	"github.com/gin-gonic/gin"
)

var users = map[string]string{
	"admin":  "password",
	"editor": "secret",
}

func StartServer() *gin.Engine {
	router := gin.Default()

	// Initialize the general repository for Suppliers
	supplierRepo := repository.NewGeneralRepository[structs.Sanco_Suppliers]("Supplier")
	suppliersController := controllers.NewGeneralController(supplierRepo)

	// Public routes
	router.POST("/login", middleware.LoginHandler)

	// Supplier routes
	router.POST("/api/master/supplier/store", suppliersController.Create)
	router.GET("/api/master/supplier/detail/:id", suppliersController.GetByID)
	router.PUT("/api/master/supplier/update/:id", suppliersController.Update)
	router.DELETE("/api/master/supplier/delete/:id", suppliersController.Delete)
	router.GET("/api/master/supplier/show", suppliersController.Get)
	router.GET("/api/master/supplier/tables", suppliersController.GetTables)

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
