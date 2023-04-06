package controller

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/jun2900/indo-dispo/entity"
	"github.com/jun2900/indo-dispo/model"
	"github.com/jun2900/indo-dispo/services"

	"github.com/gofiber/fiber/v2"
)

var layoutTime = "2006-01-02"

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
// @Param vendor query string false "filter by supplier name"
// @Param billType query string false "filter by bill type"
// @Param   dateFrom    query    string  false        "search lower limit event time"
// @Param   dateTo    query    string  false        "search upper limit event time"
// @Success 200 {object} entity.PagedResults{Data=[]model.VSupplierBill}
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /bills [get]
func (b *BillController) GetAllBills(c *fiber.Ctx) error {
	functionName := "GetAllBills"

	status := c.Query("status", "")
	vendor := c.Query("vendor", "")
	billType := c.Query("billType", "")

	order := c.Query("order", "")
	page := c.QueryInt("page", 0)
	pagesize := c.QueryInt("pagesize", 20)

	dateFrom := c.Query("dateFrom", "")
	dateTo := c.Query("dateTo", "")

	var dateFromTime time.Time
	var dateToTime time.Time
	var err error

	if dateFrom != "" {
		dateFromTime, err = time.Parse(layoutTime, dateFrom)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on parsing date from time, details = %v", err),
			})
		}
	}

	if dateTo != "" {
		dateToTime, err = time.Parse(layoutTime, dateTo)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on parsing date to time, details = %v", err),
			})
		}
	}

	bills, totalRows, err := b.billService.GetAllBill(page, pagesize, order, time.Time{}, status, vendor, billType, dateFromTime, dateToTime)
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

	startDateTime, dueDateTime, total, err := b.validateBillData(struct {
		StartDate  string
		DueDate    string
		Items      []entity.ItemPurchase
		SupplierId int32
	}{
		StartDate:  input.StartDate,
		DueDate:    input.DueDate,
		Items:      input.Items,
		SupplierId: input.SupplierId,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     err.Error(),
		})
	}

	randNum, _ := rand.Int(rand.Reader, big.NewInt(9000))
	randNum = randNum.Add(randNum, big.NewInt(1000))

	billNumber := fmt.Sprintf("IDS/%s/%d", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""), randNum)

	bill, _, err := b.billService.CreateBill(&model.Bill{
		SupplierID:        input.SupplierId,
		BillStartDate:     startDateTime,
		BillDueDate:       dueDateTime,
		BillNumber:        billNumber,
		BillOrderNumber:   nil,
		BillTotal:         total,
		BillStatus:        "MENUNGGU PEMBAYARAN",
		BillType:          input.BillType,
		BillShippingCost:  input.ShippingCost,
		BillAccountNumber: input.AccountNumber,
		BillBankName:      input.BankName,
		BillNotes:         input.BillNote,
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
		itPurchasePpn := 0
		if item.ItemPpn {
			itPurchasePpn = 1
		}

		modelItemPurchases = append(modelItemPurchases, model.ItemPurchase{
			ItemID:               item.ItemId,
			BillID:               &bill.BillID,
			RecurringBillID:      nil,
			ItemPurchaseQty:      item.ItemQty,
			ItemPurchaseTime:     time.Now(),
			ItemPurchaseDiscount: item.ItemDiscount,
			ItemPurchasePpn:      int32(itPurchasePpn),
			ItemPurchaseUnit:     item.ItemUnit,
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

		is, err := b.itemService.GetItemWithItemIdAndSupplierId(item.ItemID, bill.SupplierID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on getting item supplier with supplier id '%d' and item id '%d' details = %v", bill.SupplierID, item.ItemID, err),
			})
		}

		var amount float64
		if ip.ItemPurchaseDiscount != nil {
			amount = is.ItemPurchasePrice*float64(ip.ItemPurchaseQty) - is.ItemPurchasePrice*float64(ip.ItemPurchaseQty)**ip.ItemPurchaseDiscount/100
		} else {
			amount = is.ItemPurchasePrice * float64(ip.ItemPurchaseQty)
		}

		itemPpn := false
		if ip.ItemPurchasePpn == 1 {
			itemPpn = true
		}
		itemBills = append(itemBills, entity.ItemBill{
			Name:        item.ItemName,
			Description: item.ItemDescription,
			Qty:         ip.ItemPurchaseQty,
			Price:       is.ItemPurchasePrice,
			Amount:      amount,
			ItemPpn:     itemPpn,
			ItemUnit:    ip.ItemPurchaseUnit,
		})
	}

	total := 0
	subTotal := 0
	for _, ib := range itemBills {
		total += int(ib.Amount)
		subTotal += int(ib.Qty) * int(ib.Price)
	}

	return c.Status(fiber.StatusOK).JSON(entity.BillDetailsResp{
		StartDate:         bill.BillStartDate.Format(layoutTime),
		DueDate:           bill.BillDueDate.Format(layoutTime),
		BillNumber:        bill.BillNumber,
		BillOrderNumber:   bill.BillOrderNumber,
		BillType:          bill.BillType,
		Attachments:       attachments,
		Items:             itemBills,
		BillStatus:        bill.BillStatus,
		BillSubTotal:      int64(subTotal),
		BillTotal:         int64(total),
		BillShippingCost:  bill.BillShippingCost,
		BillAccountNumber: bill.BillAccountNumber,
		BillBankName:      bill.BillBankName,
	})
}

// @Summary Get Bill Header For Raw Only
// @Description get bill overdue open and draft stats
// @Tags Bill
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.BillHeaderResp
// @Router /bill/header [get]
func (b *BillController) GetBillHeader(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(entity.BillHeaderResp{
		Overdue:   b.billService.GetAllOverdueBillTotalWithBillType("raw"),
		Open:      b.billService.GetAllOpenBillTotal(),
		BillDraft: b.billService.GetAllMenungguPembayaranBillTotalWithBillType("raw"),
	})
}

// @Summary Update Bill Status
// @Tags Bill
// @Accept  json
// @Produce  json
// @Param  billId path int true "bill id"
// @Param  input body entity.BillUpdateStatusReq true "update bill status request (paid/cancelled)"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /bill/status/{billId} [put]
func (b *BillController) UpdateBillStatus(c *fiber.Ctx) error {
	var input entity.BillUpdateStatusReq

	functionName := "UpdateBillStatus"

	billId, err := c.ParamsInt("billId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing bill id, details = %v", err),
		})
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing item input, details = %v", err),
		})
	}

	if strings.ToLower(input.Status) != "paid" && strings.ToLower(input.Status) != "cancelled" {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "status must be paid or cancelled",
		})
	}

	if strings.ToLower(input.Status) == "paid" {
		_, _, err := b.billService.UpdateBillStatus(int32(billId), "SUDAH DIBAYAR")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on updating bill status, details = %v", err),
			})
		}
	} else {
		_, _, err := b.billService.UpdateBillStatus(int32(billId), strings.ToUpper(input.Status))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on updating bill status, details = %v", err),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(entity.StatusResponse{
		Status: "success on updating bill",
	})
}

// @Summary Update Bill
// @Description Update Bill more detail but not status
// @Tags Bill
// @Accept  json
// @Produce  json
// @Param  billId path int true "bill id"
// @Param  input body entity.UpdateBillReq true "update bill req"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /bill/{billId} [put]
func (b *BillController) UpdateBill(c *fiber.Ctx) error {
	functionName := "UpdateBill"

	var input entity.UpdateBillReq
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing item input, details = %v", err),
		})
	}

	billId, err := c.ParamsInt("billId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing bill id, details = %v", err),
		})
	}

	bill, err := b.deleteRelatedBillThings(int32(billId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     err.Error(),
		})
	}

	startDateTime, dueDateTime, total, err := b.validateBillData(struct {
		StartDate  string
		DueDate    string
		Items      []entity.ItemPurchase
		SupplierId int32
	}{
		StartDate:  input.StartDate,
		DueDate:    input.DueDate,
		Items:      input.Items,
		SupplierId: input.SupplierId,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     err.Error(),
		})
	}

	_, _, err = b.billService.UpdateBill(bill.BillID, &model.Bill{
		SupplierID:        input.SupplierId,
		BillStartDate:     startDateTime,
		BillDueDate:       dueDateTime,
		BillTotal:         total,
		BillShippingCost:  input.ShippingCost,
		BillAccountNumber: input.AccountNumber,
		BillBankName:      input.BankName,
		BillNotes:         input.BillNote,
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
		itPurchasePpn := 0
		if item.ItemPpn {
			itPurchasePpn = 1
		}

		modelItemPurchases = append(modelItemPurchases, model.ItemPurchase{
			ItemID:               item.ItemId,
			BillID:               &bill.BillID,
			ItemPurchaseQty:      item.ItemQty,
			ItemPurchaseTime:     time.Now(),
			ItemPurchaseDiscount: item.ItemDiscount,
			ItemPurchasePpn:      int32(itPurchasePpn),
			ItemPurchaseUnit:     item.ItemUnit,
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
		Status: "successfully updated bill",
	})

}

// @Summary Delete Bill
// @Description delete bill that are not paid yet
// @Tags Bill
// @Accept  json
// @Produce  json
// @Param  billId path int true "bill id"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /bill/{billId} [delete]
func (b *BillController) DeleteBill(c *fiber.Ctx) error {
	functionName := "DeleteBill"

	billId, err := c.ParamsInt("billId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing bill id, details = %v", err),
		})
	}

	bill, err := b.deleteRelatedBillThings(int32(billId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     err.Error(),
		})
	}

	if err := b.billService.DeleteBill(bill.BillID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on deleting bill, details = %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.StatusResponse{
		Status: "success on deleting bill",
	})
}

func (b *BillController) deleteRelatedBillThings(billId int32) (bill *model.Bill, err error) {
	bill, _, err = b.billService.GetBill(int32(billId))
	if err != nil {
		return nil, fmt.Errorf("error on getting bill, details = %v", err)
	}

	if strings.ToLower(bill.BillStatus) == "sudah dibayar" {
		return nil, errors.New("bill already been paid")
	}

	if err := b.itemPurchaseService.DeleteItemPurchasesByBillId(bill.BillID); err != nil {
		return nil, fmt.Errorf("error on deleting item purchases by bill id, details = %v", err)
	}

	if err := b.attachmentService.DeleteAttachmentByBillId(bill.BillID); err != nil {
		return nil, fmt.Errorf("error on deleting item purchases by bill id, details = %v", err)
	}

	return bill, nil
}

func (b *BillController) validateBillData(input struct {
	StartDate  string
	DueDate    string
	Items      []entity.ItemPurchase
	SupplierId int32
}) (startDateTime time.Time, dueDateTime time.Time, total float64, err error) {
	startDateTime, err = time.Parse(layoutTime, input.StartDate)
	if err != nil {
		return time.Time{}, time.Time{}, -1, fmt.Errorf("error on parsing start date, details = %v", err)
	}

	dueDateTime, err = time.Parse(layoutTime, input.DueDate)
	if err != nil {
		return time.Time{}, time.Time{}, -1, fmt.Errorf("error on parsing due date, details = %v", err)
	}

	if len(input.Items) < 1 {
		return time.Time{}, time.Time{}, -1, errors.New("items cannot be empty")
	}

	supplier, _, err := b.supplierService.GetSupplier(input.SupplierId)
	if err != nil {
		return time.Time{}, time.Time{}, -1, fmt.Errorf("error on getting supplier details = %v", err)
	}

	if strings.ToLower(supplier.SupplierType) != "vendor" {
		return time.Time{}, time.Time{}, -1, errors.New("supplier is not a vendor")
	}

	total = 0
	for _, item := range input.Items {
		it, err := b.itemService.GetItemWithItemIdAndSupplierId(item.ItemId, input.SupplierId)
		if err != nil {
			return time.Time{}, time.Time{}, -1, fmt.Errorf("error on getting item with item id '%d' and supplier id  '%d' details = %v", item.ItemId, input.SupplierId, err)
		}
		total += it.ItemPurchasePrice*float64(item.ItemQty) - it.ItemPurchasePrice*float64(item.ItemQty)**item.ItemDiscount/100
		if item.ItemPpn {
			total += total * 11 / 100
		}
	}

	return startDateTime, dueDateTime, total, nil
}
