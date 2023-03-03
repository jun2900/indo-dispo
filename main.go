package main

import (
	"anara/controller"
	"anara/infrastructure"
	"anara/services"
	"log"
	"os"

	_ "anara/docs"
	"anara/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	swagger "github.com/arsmn/fiber-swagger/v2"
)

func readEnvironmentFile() {
	//Environment file Load --------------------------------
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(3)
	}
}

// @title Test Accounting App
// @version 1.0

// Pinger ping
// @Summary ping example
// @Description do ping
// @Tags Ping
// @Accept json
// @Produce json
// @Param x-access-token header string true "token from login user (use Bearer in front of the jwt)"
// @Success 200 {string} string "pong"
// @Failure 400 {string} string entity.ErrRespController
// @Failure 500 {string} string entity.ErrRespController
// @Router /ping [get]
func main() {
	readEnvironmentFile()

	DB := infrastructure.OpenDbConnection()

	//services
	supplierService := services.NewSupplierService(DB)
	itemService := services.NewItemService(DB)
	billService := services.NewBillService(DB)
	itemPurchaseService := services.NewItemPurchaseService(DB)
	attachmentService := services.NewAttachmentService(DB)
	adminService := services.NewAdminService(DB)
	wholesalerService := services.NewWholesalerService(DB)

	//controllers
	supplierController := controller.NewSupplierController(supplierService, wholesalerService, itemService)
	itemController := controller.NewItemController(itemService, wholesalerService)
	billController := controller.NewBillController(supplierService, billService, itemPurchaseService, itemService, attachmentService)
	adminController := controller.NewAuthController(adminService)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders: "x-access-token, Content-type",
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/ping", middleware.JWTMiddleware(), middleware.GetDataFromJWT, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"resp": c.Locals("admin_id").(float64)})
	})

	app.Get("/suppliers", supplierController.GetAllSupplierByType)

	app.Post("/supplier", supplierController.RegisterSupplier)
	app.Get("/supplier/:supplierId/items", supplierController.GetItemsBySupplierId)
	app.Post("/item", itemController.RegisterItem)
	app.Post("/bill", billController.CreateBill)

	app.Get("/bills", billController.GetAllBills)
	app.Get("/bill/header", billController.GetBillHeader)
	app.Get("/bill/:billId", billController.GetBillDetail)

	app.Post("/admin/login", adminController.Login)

	log.Fatal(app.Listen(os.Getenv("API_SERVICE_PORT")))
}
