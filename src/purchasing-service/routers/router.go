package routers

import (
	"purchasing_service/controllers"
	"purchasing_service/middleware"
	"purchasing_service/repository"
	"purchasing_service/structs"

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
	router.POST("/api/login", middleware.LoginHandler)
	router.POST("/api/master/supplier/store", suppliersController.Create)
	router.GET("/api/master/supplier/detail/:id", suppliersController.GetByID)
	router.PUT("/api/master/supplier/update/:id", suppliersController.Update)
	router.PUT("/api/master/supplier/update/state/:id", suppliersController.UpdateState)
	router.DELETE("/api/master/supplier/delete/:id", suppliersController.Delete)
	router.GET("/api/master/supplier/show", suppliersController.Get)
	router.GET("/api/master/supplier/tables", suppliersController.GetTables)
	// Supplier routes

	// Protected routes
	protected := router.Group("/protected")
	protected.Use(middleware.AuthMiddleware())
	{

		protected.POST("/api/master/supplier/store", suppliersController.Create)
		protected.GET("/api/master/supplier/detail/:id", suppliersController.GetByID)
		protected.PUT("/api/master/supplier/update/:id", suppliersController.Update)
		protected.PUT("/api/master/supplier/update/state/:id", suppliersController.UpdateState)
		protected.DELETE("/api/master/supplier/delete/:id", suppliersController.Delete)
		protected.GET("/api/master/supplier/show", suppliersController.Get)
		protected.GET("/api/master/supplier/tables", suppliersController.GetTables)

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
