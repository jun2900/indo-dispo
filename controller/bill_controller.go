package controller

import (
	"anara/entity"
	"anara/model"
	"anara/services"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

var layoutTime = "2006-01-02 15:04:05"

type BillController struct {
	supplierService     services.SupplierService
	billService         services.BillService
	itemPurchaseService services.ItemPurchaseService
	itemService         services.ItemService
	attachmentService   services.AttachmentService
}

func NewBillController(supplierService services.SupplierService, billService services.BillService, itemPurchaseService services.ItemPurchaseService, itemService services.ItemService, attachmentService services.AttachmentService) *BillController {
	return &BillController{
		supplierService:     supplierService,
		billService:         billService,
		itemPurchaseService: itemPurchaseService,
		itemService:         itemService,
		attachmentService:   attachmentService,
	}
}

// @Summary List All Bill
// @Tags Bill
// @Accept  json
// @Produce  json
// @Param       page     				query    int    false "page requested (defaults to 0)"
// @Param       pagesize 				query    int    false "number of records in a page  (defaults to 20)"
// @Param       order    				query    string false "asc / desc"
// @Param status query string false "filter by bill status"
// @Param vendor query string false "filter by bill vendor"
// @Success 200 {object} entity.PagedResults{Data=[]model.VSupplierBill}
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /bills [get]
func (b *BillController) GetAllBills(c *fiber.Ctx) error {
	functionName := "GetAllBills"

	status := c.Query("status", "")
	vendor := c.Query("vendor", "")
	order := c.Query("order", "")

	page := c.QueryInt("page", 0)
	pagesize := c.QueryInt("pagesize", 20)

	bills, totalRows, err := b.billService.GetAllBill(page, pagesize, order, time.Time{}, status, vendor)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting bills, details = %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.PagedResults{
		Page:         page,
		PageSize:     pagesize,
		Data:         bills,
		TotalRecords: int(totalRows),
	})
}

// @Summary Register Bill
// @Tags Bill
// @Accept  json
// @Produce  json
// @Param  input body entity.AddBillReq true "add bill request"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /bill [post]
func (b *BillController) CreateBill(c *fiber.Ctx) error {
	var input entity.AddBillReq

	functionName := "CreateBill"

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

	dueDateTime, err := time.Parse(layoutTime, input.DueDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing due date, details = %v", err),
		})
	}

	if len(input.BillNumber) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "bill number cannot be empty",
		})
	}

	if len(input.Items) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "items cannot be empty",
		})
	}

	if input.Discount != nil {
		if int(*input.Discount) < 0 && int(*input.Discount) > 100 {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     "discount cannot be below 0 or pass 100",
			})
		}
	}

	supplier, _, err := b.supplierService.GetSupplier(input.SupplierId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting supplier details = %v", err),
		})
	}

	if supplier.SupplierType != "vendor" {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "supplier is not a vendor",
		})
	}

	total := 0
	for _, item := range input.Items {
		it, err := b.itemService.GetItem(item.ItemId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on getting item with id %d details = %v", item.ItemId, err),
			})
		}
		total += int(*it.ItemPurchasePrice) * int(item.ItemQty)
	}

	bill, _, err := b.billService.CreateBill(&model.Bill{
		SupplierID:      input.SupplierId,
		BillStartDate:   startDateTime,
		BillDueDate:     dueDateTime,
		BillNumber:      input.BillNumber,
		BillOrderNumber: input.BillOrderNumber,
		BillDiscount:    input.Discount,
		BillTotal:       int32(total),
		BillStatus:      "DRAFT",
		BillType:        input.BillType,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on creating bill, details = %v", err),
		})
	}

	var modelAttachments []model.Attachment

	if len(input.Attachments) > 0 {
		for _, at := range input.Attachments {
			modelAttachments = append(modelAttachments, model.Attachment{
				BillID:         &bill.BillID,
				InvoiceID:      nil,
				AttachmentName: at.Name,
				AttachmentFile: at.File,
			})
		}

		_, _, err = b.attachmentService.CreateAttachments(modelAttachments)
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
			ItemID:           item.ItemId,
			BillID:           bill.BillID,
			ItemPurchaseQty:  item.ItemQty,
			ItemPurchaseTime: time.Now(),
		})
	}
	_, _, err = b.itemPurchaseService.CreateItemPurchase(modelItemPurchases)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on creating item purchases, details = %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entity.StatusResponse{
		Status: "successfully created bill",
	})
}

// @Summary Get Bill Details
// @Tags Bill
// @Accept  json
// @Produce  json
// @Param  billId path int true "bill id"
// @Success 200 {object} entity.BillDetailsResp
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /bill/{billId} [get]
func (b *BillController) GetBillDetail(c *fiber.Ctx) error {
	functionName := "GetBillDetail"

	billId, err := c.ParamsInt("billId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing bill id, details = %v", err),
		})
	}

	bill, _, err := b.billService.GetBill(int32(billId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting bill, details = %v", err),
		})
	}

	var attachments []entity.Attachment
	attachmentRec, _, err := b.attachmentService.GetAttachmentByBillId(bill.BillID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on attachment by bill id, details = %v", err),
		})
	}

	if len(attachmentRec) > 0 {
		for _, at := range attachmentRec {
			attachments = append(attachments, entity.Attachment{
				Name: at.AttachmentName,
				File: at.AttachmentFile,
			})
		}
	}

	itemPurchases, _, err := b.itemPurchaseService.GetAllItemPurchaseByBillId(bill.BillID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting item purchases with bill id %d details = %v", bill.BillID, err),
		})
	}

	var itemBills []entity.ItemBill
	for _, ip := range itemPurchases {
		item, err := b.itemService.GetItem(ip.ItemID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on getting item purchases with bill id %d details = %v", bill.BillID, err),
			})
		}
		itemBills = append(itemBills, entity.ItemBill{
			Name:        item.ItemName,
			Description: item.ItemDescription,
			Qty:         ip.ItemPurchaseQty,
			Price:       *item.ItemPurchasePrice,
			Amount:      *item.ItemPurchasePrice*ip.ItemPurchaseQty - *item.ItemPurchasePrice*ip.ItemPurchaseQty**bill.BillDiscount/100,
		})
	}

	total := 0
	subTotal := 0
	for _, ib := range itemBills {
		total += int(ib.Amount)
		subTotal += int(ib.Qty) * int(ib.Price)
	}

	return c.Status(fiber.StatusOK).JSON(entity.BillDetailsResp{
		StartDate:       bill.BillStartDate.Format(layoutTime),
		DueDate:         bill.BillDueDate.Format(layoutTime),
		BillNumber:      bill.BillNumber,
		BillOrderNumber: bill.BillOrderNumber,
		BillType:        bill.BillType,
		Attachments:     attachments,
		Items:           itemBills,
		BillStatus:      bill.BillStatus,
		BillSubTotal:    int64(subTotal),
		BillTotal:       int64(total),
		BillDiscount:    bill.BillDiscount,
	})
}

// @Summary Get Bill Header
// @Description get bill overdue open and draft stats
// @Tags Bill
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.BillHeaderResp
// @Router /bill/header [get]
func (b *BillController) GetBillHeader(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(entity.BillHeaderResp{
		Overdue:   b.billService.GetAllOverdueBillTotal(),
		Open:      b.billService.GetAllOpenBillTotal(),
		BillDraft: b.billService.GetAllDraftBillTotal(),
	})
}
