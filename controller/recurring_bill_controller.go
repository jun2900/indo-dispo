package controller

import (
	"anara/entity"
	"anara/model"
	"anara/services"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	enumFrequency = []string{
		"daily", "weekly", "monthly",
	}
)

type RecurringBillController struct {
	recurringBillService services.RecurringBillService
	itemPurchaseService  services.ItemPurchaseService
	supplierService      services.SupplierService
	itemService          services.ItemService
	attachmentService    services.AttachmentService
}

func NewRecurringBillController(recurringBilService services.RecurringBillService,
	itemPurchaseService services.ItemPurchaseService,
	supplierService services.SupplierService,
	itemService services.ItemService,
	attachmentService services.AttachmentService) *RecurringBillController {
	return &RecurringBillController{
		recurringBillService: recurringBilService,
		itemPurchaseService:  itemPurchaseService,
		supplierService:      supplierService,
		itemService:          itemService,
		attachmentService:    attachmentService,
	}
}

// @Summary Register Recurring Bill
// @Tags Recurring Bill
// @Accept  json
// @Produce  json
// @Param  input body entity.AddRecurringBillReq true "add recurring bill request"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /recurring_bill [post]
func (r *RecurringBillController) AddRecurringBill(c *fiber.Ctx) error {
	functionName := "AddRecurringBill"

	var input entity.AddRecurringBillReq

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing item input, details = %v", err),
		})
	}

	startDateTime, err := time.Parse(layoutTime, input.StartDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing start date, details = %v", err),
		})
	}

	var endDate *time.Time
	if input.EndDate != nil {
		endDateTime, err := time.Parse(layoutTime, *input.EndDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on parsing due date, details = %v", err),
			})
		}
		endDate = &endDateTime
	}

	if len(input.Items) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "items cannot be empty",
		})
	}

	if len(input.BankName) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "bank name cannot be empty",
		})
	}

	if len(input.AccountNumber) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "account number cannot be empty",
		})
	}

	if input.PaymentDue < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "payment due must be more than 0",
		})
	}

	if !stringInSlice(input.Frequency, enumFrequency) {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "frequency must be (daily, weeekly, monthly)",
		})
	}

	supplier, _, err := r.supplierService.GetSupplier(input.SupplierId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting supplier details = %v", err),
		})
	}

	if strings.ToLower(supplier.SupplierType) != "vendor" {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "supplier is not a vendor",
		})
	}

	var total float64
	total = 0
	for _, item := range input.Items {
		it, err := r.itemService.GetItemWithItemIdAndSupplierId(item.ItemId, input.SupplierId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on getting item with item id '%d' and supplier id  '%d' details = %v", item.ItemId, input.SupplierId, err),
			})
		}
		total += it.ItemPurchasePrice*float64(item.ItemQty) - it.ItemPurchasePrice*float64(item.ItemQty)**item.ItemDiscount/100
	}

	//randNum, _ := rand.Int(rand.Reader, big.NewInt(9000))
	//randNum = randNum.Add(randNum, big.NewInt(1000))

	//billNumber := fmt.Sprintf("IDS/%s/%d", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""), randNum)

	recurBill, _, err := r.recurringBillService.CreateRecurringBill(&model.RecurringBill{
		SupplierID:    input.SupplierId,
		Frequency:     input.Frequency,
		Total:         total,
		Notes:         input.Notes,
		ShippingCost:  input.ShippingCost,
		AccountNumber: input.AccountNumber,
		BankName:      input.BankName,
		StartDate:     startDateTime,
		EndDate:       endDate,
		PaymentDue:    int32(input.PaymentDue),
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on creating recurring bill, details = %v", err),
		})
	}

	var modelAttachments []model.Attachment
	if len(input.Attachments) > 0 {
		for _, at := range input.Attachments {
			modelAttachments = append(modelAttachments, model.Attachment{
				BillID:          nil,
				RecurringBillID: &recurBill.ID,
				InvoiceID:       nil,
				AttachmentName:  at.Name,
				AttachmentFile:  at.File,
			})
		}

		_, _, err = r.attachmentService.CreateAttachments(modelAttachments)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on creating attachments, details = %v", err),
			})
		}
	}

	var modelItemPurchases []model.ItemPurchase
	for _, item := range input.Items {
		modelItemPurchases = append(modelItemPurchases, model.ItemPurchase{
			ItemID:               item.ItemId,
			BillID:               nil,
			RecurringBillID:      &recurBill.ID,
			ItemPurchaseQty:      item.ItemQty,
			ItemPurchaseTime:     time.Now(),
			ItemPurchaseDiscount: item.ItemDiscount,
		})
	}
	_, _, err = r.itemPurchaseService.CreateItemPurchase(modelItemPurchases)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on creating item purchases, details = %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entity.StatusResponse{
		Status: "successfully created recurring bill",
	})
}

func stringInSlice(s string, slice []string) bool {
	for _, str := range slice {
		if strings.EqualFold(str, s) {
			return true
		}
	}
	return false
}