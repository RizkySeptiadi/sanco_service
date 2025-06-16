package routers

import (
	"accounting_service/controllers"
	"accounting_service/middleware"
	"accounting_service/repository"
	"accounting_service/structs"

	"github.com/gin-gonic/gin"
)

var users = map[string]string{
	"admin":  "password",
	"editor": "secret",
}

func StartServer() *gin.Engine {
	router := gin.Default()

	// Initialize the general repository for Suppliers
	supplierRepo := repository.NewGeneralRepository[structs.Suppliers]("Suppliers")
	customersRepo := repository.NewGeneralRepository[structs.Customers]("Customers")
	companiesRepo := repository.NewGeneralRepository[structs.Companies]("Companies")

	sales_invoice_repo := repository.NewGeneralSalesInvoiceRepository[structs.SalesInvoice]("SalesInvoice")
	purchase_invoice_repo := repository.NewGeneralPurchaseInvoiceRepository[structs.PurchaseInvoice]("PurchaseInvoice")

	// purchase_invoices_detail_repo := repository.NewGeneralSalesInvoiceRepository[structs.Sanco_Purchase_Invoice_details]("PurchaseInvoice")

	suppliersController := controllers.NewGeneralController(supplierRepo)
	customersController := controllers.NewGeneralController(customersRepo)
	conpaniesController := controllers.NewGeneralController(companiesRepo)

	sales_invoice_controller := controllers.NewSalesInvoiceController(sales_invoice_repo)

	purchase_invoice_controller := controllers.NewPurchaseInvoiceController(purchase_invoice_repo)

	// purchase_invoice_detail_controller := controllers.NewSalesInvoiceController(purchase_invoices_detail_repo)

	// Public routes
	router.POST("/api/login", middleware.LoginHandler)

	// Supplier routes

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// protected.POST("/master/supplier/store", suppliersController.Create)
		// protected.POST("/master/supplier/update", suppliersController.Create)
		protected.POST("/master/supplier/upsert", suppliersController.UpsertBatch)

		// protected.POST("/master/customer/store", customersController.UpsertBatch)
		// protected.POST("/master/customer/update", customersController.UpsertBatch)
		protected.POST("/master/customer/upsert", customersController.UpsertBatch)

		// protected.POST("/master/company/store", conpaniesController.Update)
		protected.POST("/master/company/upsert", conpaniesController.UpsertBatch)
		// protected.POST("/master/company/update", conpaniesController.Update)

		protected.POST("/sales_invoice/store", sales_invoice_controller.Create)
		protected.POST("/purchase_invoice/store", purchase_invoice_controller.Create)

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
