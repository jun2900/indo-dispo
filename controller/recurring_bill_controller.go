package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/jun2900/indo-dispo/entity"
	"github.com/jun2900/indo-dispo/model"
	"github.com/jun2900/indo-dispo/services"

	"github.com/gofiber/fiber/v2"
)

var (
	enumFrequency = []string{
		"daily", "weekly", "monthly",
	}
	enumStatus = []string{
		"active", "inactive",
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
		Status:        "active",
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

// @Summary List All Recurring Bills
// @Tags Recurring Bill
// @Accept  json
// @Produce  json
// @Param       page     				query    int    false "page requested (defaults to 0)"
// @Param       pagesize 				query    int    false "number of records in a page  (defaults to 20)"
// @Param       order    				query    string false "asc / desc"
// @Success 200 {object} entity.PagedResults{Data=[]model.RecurringBill}
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /recurring_bills [get]
func (r *RecurringBillController) GetAllRecurringBill(c *fiber.Ctx) error {
	functionName := "GetAllRecurringBill"

	page := c.QueryInt("page", 0)
	pagesize := c.QueryInt("pagesize", 20)
	order := c.Query("order", "")

	recurBills, totalRecords, err := r.recurringBillService.GetAllRecurringBill(page, pagesize, order)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting recur bills, details = %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.PagedResults{
		Page:         page,
		PageSize:     pagesize,
		Data:         recurBills,
		TotalRecords: int(totalRecords),
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

// @Summary Update Recurring Bill Status
// @Tags Recurring Bill
// @Accept  json
// @Produce  json
// @Param  recurringBillId path int true "recurring bill id"
// @Param  input body entity.UpdateBillRecurringStatus true "update bill status request (active/inactive)"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /recurring_bill/status/{recurringBillId} [put]
func (r *RecurringBillController) UpdateStatusRecurringBill(c *fiber.Ctx) error {
	functionName := "UpdateRecurringBill"

	recurrBillId, err := c.ParamsInt("recurringBillId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing params item id, details = %v", err),
		})
	}

	var input entity.UpdateBillRecurringStatus
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing item input, details = %v", err),
		})
	}

	if !stringInSlice(input.Status, enumStatus) {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "status must be active or inactive",
		})
	}

	_, _, err = r.recurringBillService.UpdateRecurringBill(int32(recurrBillId), &model.RecurringBill{
		Status: input.Status,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on updating recurring bill, details = %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.StatusResponse{
		Status: "successfully update recurring bill",
	})
}
