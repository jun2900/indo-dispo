package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/jun2900/indo-dispo/controller"
	"github.com/jun2900/indo-dispo/infrastructure"
	"github.com/jun2900/indo-dispo/model"
	"github.com/jun2900/indo-dispo/services"

	_ "github.com/jun2900/indo-dispo/docs"
	"github.com/jun2900/indo-dispo/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

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
	go scheduleDailyTask(18, 0, 0, func() {
		scriptCheckRecurring(DB)
	})

	//services
	supplierService := services.NewSupplierService(DB)
	itemService := services.NewItemService(DB)
	billService := services.NewBillService(DB)
	itemPurchaseService := services.NewItemPurchaseService(DB)
	attachmentService := services.NewAttachmentService(DB)
	adminService := services.NewAdminService(DB)
	wholesalerService := services.NewWholesalerService(DB)
	balanceService := services.NewBalanceService(DB)
	balanceLogService := services.NewBalanceLogService(DB)
	recurringBillService := services.NewRecurringBillService(DB)
	invoiceService := services.NewInvoiceService(DB)
	itemSellService := services.NewItemSellService(DB)

	//controllers
	supplierController := controller.NewSupplierController(supplierService, wholesalerService, itemService)
	itemController := controller.NewItemController(itemService, wholesalerService)
	billController := controller.NewBillController(supplierService, billService, itemPurchaseService, itemService, attachmentService)
	adminController := controller.NewAuthController(adminService)
	balanceController := controller.NewBalanceController(balanceService, billService, balanceLogService)
	recurringBillController := controller.NewRecurringBillController(recurringBillService, itemPurchaseService, supplierService, itemService, attachmentService)
	invoiceController := controller.NewInvoiceController(invoiceService, supplierService, itemService, attachmentService, itemSellService)

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
	app.Get("/supplier/details/:supplierId", supplierController.GetSupplierDetail)
	app.Post("/supplier/items", supplierController.ListSupplierByItemName)
	app.Post("/item", itemController.RegisterItem)
	app.Put("/item/:itemId", itemController.UpdateItem)
	app.Delete("/item/:itemId", itemController.DeleteItem)

	app.Post("/bill", billController.CreateBill)
	app.Get("/bills", billController.GetAllBills)
	app.Get("/bill/header", billController.GetBillHeader)
	app.Put("/bill/status/:billId", billController.UpdateBillStatus)
	app.Put("/bill/:billId", billController.UpdateBill)
	app.Get("/bill/:billId", billController.GetBillDetail)
	app.Delete("/bill/:billId", billController.DeleteBill)

	app.Get("/balance/header", balanceController.GetNetBalanceAmount)
	app.Put("/balance", balanceController.AddBalanceAmount)
	app.Get("/balance/logs", balanceController.GetAllBalanceLog)

	app.Get("/recurring_bills", recurringBillController.GetAllRecurringBill)
	app.Post("/recurring_bill", recurringBillController.AddRecurringBill)
	app.Put("recurring_bill/status/:recurringBillId", recurringBillController.UpdateStatusRecurringBill)

	app.Get("/invoices", middleware.DBTransactionMiddleware(DB), invoiceController.GetAllInvoices)
	app.Post("/invoice", middleware.DBTransactionMiddleware(DB), invoiceController.CreateInvoice)

	app.Post("/admin/login", adminController.Login)

	log.Fatal(app.Listen(os.Getenv("API_SERVICE_PORT")))

	select {}
}

func scriptCheckRecurring(db *gorm.DB) {
	recurringBill := []model.RecurringBill{}
	if err := db.Model(&model.RecurringBill{}).Find(&recurringBill).Error; err != nil {
		panic(fmt.Sprintf("error on finding recurring bills, details = %v", recurringBill))
	}

	if len(recurringBill) > 0 {
		for _, rb := range recurringBill {
			if time.Now().After(rb.StartDate) {
				if rb.EndDate == nil || time.Now().Before(*rb.EndDate) {
					switch strings.ToLower(rb.Frequency) {
					case "daily":
						processDailyTransaction(db, rb)
					case "weekly":
						processWeeklyTransaction(db, rb)
					case "monthly":
						processMonthlyTransaction(db, rb)
					}
				}
			}
		}
	}
}

func scheduleDailyTask(hour, min, sec int, task func()) {
	now := time.Now()
	// Calculate the time until the next occurrence of the desired time
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())
	if next.Before(now) {
		// If the desired time has already passed today, schedule it for tomorrow
		next = next.Add(24 * time.Hour)
	}

	// Calculate the duration until the next occurrence
	duration := next.Sub(now)
	// Schedule the task to run after the calculated duration
	time.AfterFunc(duration, func() {
		// Run the task
		task()
		// Schedule the task to run again in 24 hours
		scheduleDailyTask(hour, min, sec, task)
	})
}

func processDailyTransaction(db *gorm.DB, recurringTxn model.RecurringBill) {
	// Check if a transaction has already been recorded for today
	todaysTxn := model.Bill{}
	if err := db.Where("recurring_bill_id = ? AND bill_start_date = CURDATE()", recurringTxn.ID).First(&todaysTxn).Error; err == gorm.ErrRecordNotFound {
		// If no transaction has been recorded for today, create a new one
		insertingAttachmentAndItemPurchase(db, recurringTxn)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Daily transaction for Gym already recorded for today")
	}
}

func processWeeklyTransaction(db *gorm.DB, recurringTxn model.RecurringBill) {
	// Check if a transaction has already been recorded for this week
	weekStart := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -int(time.Now().Weekday()))
	weekEnd := weekStart.AddDate(0, 0, 7).Add(-time.Nanosecond)
	weeklyTxn := model.Bill{}
	if err := db.Where("recurring_bill_id = ? AND bill_start_date >= ? AND bill_start_date <= ?", recurringTxn.ID, weekStart.Format("2006-01-02"), weekEnd.Format("2006-01-02")).First(&weeklyTxn).Error; err == gorm.ErrRecordNotFound {
		// If no transaction has been recorded for this week, create a new one
		insertingAttachmentAndItemPurchase(db, recurringTxn)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Weekly transaction for Gym already recorded for this week")
	}
}

func processMonthlyTransaction(db *gorm.DB, recurringTxn model.RecurringBill) {
	// Check if a transaction has already been recorded for this month
	monthStart := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -time.Now().Day()+1)
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)
	monthlyTxn := model.Bill{}
	if err := db.Where("recurring_bill_id = ? AND bill_start_date >= ? AND bill_start_date <= ?", recurringTxn.ID, monthStart.Format("2006-01-02"), monthEnd.Format("2006-01-02")).First(&monthlyTxn).Error; err == gorm.ErrRecordNotFound {
		insertingAttachmentAndItemPurchase(db, recurringTxn)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Monthly transaction for Gym already recorded for this month")
	}
}

func insertingAttachmentAndItemPurchase(db *gorm.DB, recurringTxn model.RecurringBill) {
	randNum, _ := rand.Int(rand.Reader, big.NewInt(9000))
	randNum = randNum.Add(randNum, big.NewInt(1000))
	billNumber := fmt.Sprintf("IDS/%s/%d", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""), randNum)

	newTxn := model.Bill{
		SupplierID:        recurringTxn.SupplierID,
		BillStartDate:     time.Now(),
		BillDueDate:       time.Now().AddDate(0, 0, int(recurringTxn.PaymentDue)),
		BillNumber:        billNumber,
		BillOrderNumber:   nil,
		BillTotal:         recurringTxn.Total,
		BillStatus:        "MENUNGGU PEMBAYARAN",
		BillType:          strings.ToLower("operasional"),
		BillShippingCost:  recurringTxn.ShippingCost,
		BillAccountNumber: recurringTxn.AccountNumber,
		BillBankName:      recurringTxn.BankName,
		BillNotes:         recurringTxn.Notes,
		RecurringBillID:   &recurringTxn.ID,
	}
	if err := db.Create(&newTxn).Error; err != nil {
		panic(fmt.Sprintf("error on creating new bill, details = %v", err))
	}

	var modelAttachments []model.Attachment
	var inputModelAttachment []model.Attachment

	if err := db.Model(&model.Attachment{}).Where("recurring_bill_id = ?", recurringTxn.ID).Find(&modelAttachments).Error; err != nil {
		panic(fmt.Sprintf("error on getting attachments, details = %v", err))
	}
	if len(modelAttachments) > 0 {
		for _, ma := range modelAttachments {
			inputModelAttachment = append(inputModelAttachment, model.Attachment{
				BillID:          &newTxn.BillID,
				RecurringBillID: &recurringTxn.ID,
				InvoiceID:       nil,
				AttachmentName:  ma.AttachmentName,
				AttachmentFile:  ma.AttachmentFile,
			})
		}
		if err := db.Model(&model.Attachment{}).Save(&inputModelAttachment).Error; err != nil {
			panic(fmt.Sprintf("error on saving attachments, details = %v", err))
		}
	}

	var modelItemPurchases []model.ItemPurchase
	var inputModelItemPurchases []model.ItemPurchase
	if err := db.Model(&model.ItemPurchase{}).Where("recurring_bill_id = ?", recurringTxn.ID).Find(&modelItemPurchases).Error; err != nil {
		panic(fmt.Sprintf("error on getting item purchases, details = %v", err))
	}
	for _, mip := range modelItemPurchases {
		inputModelItemPurchases = append(inputModelItemPurchases, model.ItemPurchase{
			ItemID:               mip.ItemID,
			BillID:               &newTxn.BillID,
			RecurringBillID:      &recurringTxn.ID,
			ItemPurchaseQty:      mip.ItemPurchaseQty,
			ItemPurchaseTime:     mip.ItemPurchaseTime,
			ItemPurchaseDiscount: mip.ItemPurchaseDiscount,
			ItemPurchasePpn:      mip.ItemPurchasePpn,
			ItemPurchaseUnit:     mip.ItemPurchaseUnit,
		})
	}
	if err := db.Model(&model.ItemPurchase{}).Save(&inputModelItemPurchases).Error; err != nil {
		panic(fmt.Sprintf("error on saving item purchases, details = %v", err))
	}

	fmt.Println("Processed bill transaction")
}
