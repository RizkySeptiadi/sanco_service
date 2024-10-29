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
	purchase_invoices_repo := repository.NewGeneralPurchaseInvoiceRepository[structs.Sanco_Purchase_Invoices]("PurchaseInvoice")
	purchase_invoices_detail_repo := repository.NewGeneralPurchaseInvoiceRepository[structs.Sanco_Purchase_Invoice_details]("PurchaseInvoice")
	warehouse_repo := repository.NewGeneralWarehouseRepository[structs.Incoming_Order]("Incoming_order")

	suppliersController := controllers.NewGeneralController(supplierRepo)

	purchase_invoice_controller := controllers.NewPurchaseInvoiceController(purchase_invoices_repo)
	warehouse_controller := controllers.NewWarehouseController(warehouse_repo)

	purchase_invoice_detail_controller := controllers.NewPurchaseInvoiceController(purchase_invoices_detail_repo)

	// Public routes
	router.POST("/api/login", middleware.LoginHandler)

	// Supplier routes

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/master/supplier/store", suppliersController.Create)
		protected.GET("/master/supplier/detail/:id", suppliersController.GetByID)
		protected.PUT("/master/supplier/update/:id", suppliersController.Update)
		protected.PUT("/master/supplier/update/state/:id", suppliersController.UpdateState)
		protected.DELETE("/master/supplier/delete/:id", suppliersController.Delete)
		protected.GET("/master/supplier/show", suppliersController.Get)
		protected.GET("/master/supplier/tables", suppliersController.GetTables)
		protected.GET("/purchase_invoice/tables", purchase_invoice_controller.GetTables)
		protected.GET("/purchase_invoice/show", suppliersController.Get)
		protected.GET("/purchase_invoice/get_parent/:id", purchase_invoice_controller.GetByID)
		protected.GET("/purchase_invoice/get_all_detail/:id", purchase_invoice_detail_controller.GetAllDataDetailByID)
		protected.GET("/purchase_invoice/get_detail/:id", purchase_invoice_detail_controller.GetDataDetailByID)

		protected.POST("/purchase_invoice/store", purchase_invoice_controller.Create)
		protected.PUT("/purchase_invoice/update/:id", purchase_invoice_controller.Update)
		protected.PUT("/purchase_invoice/update/state/:id", purchase_invoice_controller.UpdateState)

		protected.PUT("/warehouse_request/create/:id", warehouse_controller.Create)
		protected.DELETE("/purchase_invoice/delete/:id", purchase_invoice_controller.Delete)

		// protected.GET("/api/master/purchase_invoice/tables", purchase_invoice_controller.GetTables)
		// protected.GET("/api/master/purchase_invoice/show", purchase_invoice_controller.Get)
		// protected.GET("/api/master/purchase_invoice/tables", purchase_invoice_controller.GetTables)

		// protected.GET("/api/master/purchase_invoice/tables", purchase_invoice_controller.GetTables)

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
