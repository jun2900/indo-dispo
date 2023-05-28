package controller

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jun2900/indo-dispo/entity"
	"github.com/jun2900/indo-dispo/model"
	"github.com/jun2900/indo-dispo/services"
	"gorm.io/gorm"
)

type InvoiceController struct {
	invoiceService    services.InvoiceService
	supplierService   services.SupplierService
	itemService       services.ItemService
	attachmentService services.AttachmentService
	itemSellService   services.ItemSellService
}

func NewInvoiceController(invoiceService services.InvoiceService, supplierService services.SupplierService,
	itemService services.ItemService, attachmentService services.AttachmentService, itemSellService services.ItemSellService) *InvoiceController {
	return &InvoiceController{
		invoiceService:    invoiceService,
		supplierService:   supplierService,
		itemService:       itemService,
		attachmentService: attachmentService,
		itemSellService:   itemSellService,
	}
}

// @Summary List All Invoices
// @Tags Invoice
// @Accept  json
// @Produce  json
// @Param       page     				query    int    false "page requested (defaults to 0)"
// @Param       pagesize 				query    int    false "number of records in a page  (defaults to 20)"
// @Param       order    				query    string false "asc / desc"
// @Param status query string false "filter by invoice status"
// @Param customer query string false "filter by customer name"
// @Param invoiceType query string false "filter by invoice type"
// @Param   dateFrom    query    string  false        "search lower limit start date time"
// @Param   dateTo    query    string  false        "search upper limit start date time"
// @Success 200 {object} entity.PagedResults{Data=[]model.VSupplierInvoice}
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /invoices [get]
func (i *InvoiceController) GetAllInvoices(c *fiber.Ctx) error {
	functionName := "GetAllInvoices"

	status := c.Query("status", "")
	customer := c.Query("customer", "")

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

	trxHandle := c.Locals("db_trx").(*gorm.DB)

	bills, totalRows, err := i.invoiceService.WithTrx(trxHandle).GetAllInvoices(page, pagesize, order, status, customer, dateFromTime, dateToTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting invoices, details = %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.PagedResults{
		Page:         page,
		PageSize:     pagesize,
		Data:         bills,
		TotalRecords: int(totalRows),
	})
}

// @Summary Get Invoice Details
// @Tags Invoice
// @Accept  json
// @Produce  json
// @Param  invoiceId path int true "invoice id"
// @Success 200 {object} entity.BillDetailsResp
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /invoice/{invoiceId} [get]
func (i *InvoiceController) GetInvoiceDetail(c *fiber.Ctx) error {
	functionName := "GetInvoiceDetail"

	invoiceId, err := c.ParamsInt("invoiceId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing invoice id, details = %v", err),
		})
	}

	trxHandle := c.Locals("db_trx").(*gorm.DB)
	invoice, _, err := i.invoiceService.WithTrx(trxHandle).GetInvoice(int32(invoiceId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting invoice, details = %v", err),
		})
	}

	var attachments []entity.Attachment
	attachmentRec, _, err := i.attachmentService.WithTrx(trxHandle).GetAttachmentByInvoiceId(invoice.InvoicesID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on attachment by invoice id, details = %v", err),
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

	itemSells, _, err := i.itemSellService.WithTrx(trxHandle).GetAllItemSellsByInvoiceId(invoice.InvoicesID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting item purchases with invoice id %d details = %v", invoice.InvoicesID, err),
		})
	}

	var itemInvoices []entity.ItemBill
	for _, ip := range itemSells {
		item, err := i.itemService.GetItem(ip.ItemID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on getting item purchases with invoice id %d details = %v", invoice.InvoicesID, err),
			})
		}

		//is, err := i.itemService.WithTrx(trxHandle).GetItemWithItemIdAndSupplierId(item.ItemID, invoice.SupplierID)
		//if err != nil {
		//	return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
		//		SourceFunction: functionName,
		//		ErrMessage:     fmt.Sprintf("error on getting item supplier with supplier id '%d' and item id '%d' details = %v", invoice.SupplierID, item.ItemID, err),
		//	})
		//}

		var amount float64
		if ip.ItemSellDiscount != nil {
			amount = item.ItemPurchasePrice*float64(ip.ItemSellQty) - item.ItemPurchasePrice*float64(ip.ItemSellQty)**ip.ItemSellDiscount/100
		} else {
			amount = item.ItemPurchasePrice * float64(ip.ItemSellQty)
		}

		itemPpn := false
		if ip.ItemSellPpn == 1 {
			itemPpn = true
		}
		itemInvoices = append(itemInvoices, entity.ItemBill{
			Id:          item.ItemID,
			Name:        item.ItemName,
			Description: item.ItemDescription,
			Qty:         ip.ItemSellQty,
			Price:       item.ItemPurchasePrice,
			Amount:      amount,
			ItemPpn:     itemPpn,
			ItemUnit:    ip.ItemSellUnit,
		})
	}

	total := 0
	subTotal := 0
	for _, ib := range itemInvoices {
		total += int(ib.Amount)
		subTotal += int(ib.Qty) * int(ib.Price)
	}

	return c.Status(fiber.StatusOK).JSON(entity.InvoiceDetailResp{
		StartDate:          invoice.InvoiceStartDate.Format(layoutTime),
		DueDate:            invoice.InvoiceDueDate.Format(layoutTime),
		InvoiceNumber:      invoice.InvoiceNumber,
		InvoiceOrderNumber: invoice.InvoiceOrderNumber,
		Attachments:        attachments,
		Items:              itemInvoices,
		InvoiceStatus:      invoice.InvoiceStatus,
		InvoiceSubTotal:    int64(subTotal),
		InvoiceTotal:       int64(total),
		Logo:               invoice.InvoiceLogo,
		Title:              invoice.InvoiceTitle,
		Subheading:         invoice.InvoiceSubheading,
		ShippingCost:       invoice.InvoiceShippingCost,
		BankName:           invoice.InvoiceBankName,
		Notes:              invoice.InvoiceNotes,
	})
}

// @Summary Register Invoice
// @Tags Invoice
// @Accept  json
// @Produce  json
// @Param  input body entity.AddInvoiceReq true "add invoice request"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /invoice [post]
func (i *InvoiceController) CreateInvoice(c *fiber.Ctx) error {
	var input entity.AddInvoiceReq

	functionName := "CreateInvoice"

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

	supplier, _, err := i.supplierService.GetSupplier(input.CustomerId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting customer, details = %v", err),
		})
	}

	if !strings.EqualFold(supplier.SupplierType, "customer") {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "is not a customer",
		})
	}

	var total float64
	total = 0

	for _, item := range input.Items {
		it, err := i.itemService.GetItem(item.ItemId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on getting item with item id '%d' details = %v", item.ItemId, err),
			})
		}
		total += it.ItemSellPrice*float64(item.ItemQty) - it.ItemSellPrice*float64(item.ItemQty)**item.ItemDiscount/100
		if item.ItemPpn {
			total += total * 11 / 100
		}
	}

	total += input.ShippingCost

	randNum, _ := rand.Int(rand.Reader, big.NewInt(9000))
	randNum = randNum.Add(randNum, big.NewInt(1000))

	invoiceNumber := fmt.Sprintf("IDS/%s/%d", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""), randNum)

	trxHandle := c.Locals("db_trx").(*gorm.DB)

	invoice, _, err := i.invoiceService.WithTrx(trxHandle).CreateInvoice(&model.Invoice{
		SupplierID:           input.CustomerId,
		InvoiceStartDate:     startDateTime,
		InvoiceDueDate:       dueDateTime,
		InvoiceNumber:        invoiceNumber,
		InvoiceTitle:         input.Title,
		InvoiceSubheading:    input.Subheading,
		InvoiceLogo:          input.Logo,
		InvoiceShippingCost:  input.ShippingCost,
		InvoiceAccountNumber: input.AccountNumber,
		InvoiceBankName:      input.BankName,
		InvoiceTotal:         total,
		InvoiceStatus:        "MENUNGGU PEMBAYARAN",
	})
	if err != nil {
		log.Println("testt create invoice err")
		trxHandle.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on creating invoice, details = %v", err),
		})
	}

	var modelAttachments []model.Attachment

	if len(input.Attachments) > 0 {
		for _, at := range input.Attachments {
			modelAttachments = append(modelAttachments, model.Attachment{
				BillID:         nil,
				InvoiceID:      &invoice.InvoicesID,
				AttachmentName: at.Name,
				AttachmentFile: at.File,
			})
		}

		_, _, err = i.attachmentService.WithTrx(trxHandle).CreateAttachments(modelAttachments)
		if err != nil {
			log.Println("testt create attachment err")
			trxHandle.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on creating attachments, details = %v", err),
			})
		}
	}

	var modelItemSells []model.ItemSell
	for _, item := range input.Items {
		itSellPpn := 0
		if item.ItemPpn {
			itSellPpn = 1
		}

		modelItemSells = append(modelItemSells, model.ItemSell{
			ItemID:           item.ItemId,
			InvoiceID:        invoice.InvoicesID,
			ItemSellQty:      item.ItemQty,
			ItemSellTime:     time.Now(),
			ItemSellDiscount: item.ItemDiscount,
			ItemSellPpn:      int32(itSellPpn),
			ItemSellUnit:     item.ItemUnit,
		})
	}

	_, _, err = i.itemSellService.WithTrx(trxHandle).CreateItemSell(modelItemSells)
	if err != nil {
		log.Println("testt create item sell err")
		trxHandle.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on creating item sell, details = %v", err),
		})
	}

	//if err := trxHandle.Commit().Error; err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
	//		SourceFunction: functionName,
	//		ErrMessage:     fmt.Sprintf("error on commiting transcation, details = %v", err),
	//	})
	//}

	return c.Status(fiber.StatusCreated).JSON(entity.StatusResponse{
		Status: "successfully created invoices",
	})
}

// @Summary Get Invoice Header
// @Description get invoice overdue open and draft stats
// @Tags Invoice
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.InvoiceHeaderResp
// @Router /invoice/header [get]
func (i *InvoiceController) GetInvoiceHeader(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(entity.InvoiceHeaderResp{
		Overdue: i.invoiceService.GetAllOverdueInvoiceTotal(),
		Open:    i.invoiceService.GetAllMenungguPembayaranInvoiceTotal(),
	})
}

// @Summary Update Invoice Status
// @Tags Invoice
// @Accept  json
// @Produce  json
// @Param  invoiceId path int true "invoice id"
// @Param  input body entity.InvoiceUpdateStatusReq true "update invoice status request (paid/cancelled)"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /invoice/status/{invoiceId} [put]
func (i *InvoiceController) UpdateInvoiceStatus(c *fiber.Ctx) error {
	var input entity.InvoiceUpdateStatusReq

	functionName := "UpdateInvoiceStatus"

	invoiceId, err := c.ParamsInt("invoiceId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing invoice id, details = %v", err),
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
		_, _, err := i.invoiceService.UpdateInvoiceStatus(int32(invoiceId), "SUDAH DIBAYAR")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on updating invoice status, details = %v", err),
			})
		}
	} else {
		_, _, err := i.invoiceService.UpdateInvoiceStatus(int32(invoiceId), strings.ToUpper(input.Status))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on updating invoice status, details = %v", err),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(entity.StatusResponse{
		Status: "success on updating invoice",
	})
}

// @Summary Update Invoice
// @Description Update Invoice more detail but not status
// @Tags Invoice
// @Accept  json
// @Produce  json
// @Param  invoiceId path int true "invoice id"
// @Param  input body entity.UpdateInvoiceReq true "update invoice req"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /invoice/{invoiceId} [put]
func (i *InvoiceController) UpdateInvoice(c *fiber.Ctx) error {
	functionName := "UpdateInvoice"

	var input entity.UpdateInvoiceReq
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing item input, details = %v", err),
		})
	}

	invoiceId, err := c.ParamsInt("invoiceId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing invoice id, details = %v", err),
		})
	}

	trxHandle := c.Locals("db_trx").(*gorm.DB)

	invoice, err := i.deleteInvoiceRelatedThings(trxHandle, int32(invoiceId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     err.Error(),
		})
	}

	startDateTime, dueDateTime, total, err := i.validateInvoiceData(struct {
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

	total += input.ShippingCost

	_, _, err = i.invoiceService.UpdateInvoice(invoice.InvoicesID, &model.Invoice{
		SupplierID:           input.SupplierId,
		InvoiceStartDate:     startDateTime,
		InvoiceDueDate:       dueDateTime,
		InvoiceTitle:         input.Title,
		InvoiceSubheading:    input.Subheading,
		InvoiceLogo:          input.Logo,
		InvoiceShippingCost:  input.ShippingCost,
		InvoiceAccountNumber: input.AccountNumber,
		InvoiceBankName:      input.BankName,
		InvoiceTotal:         total,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on creating invoice, details = %v", err),
		})
	}

	var modelAttachments []model.Attachment

	if len(input.Attachments) > 0 {
		for _, at := range input.Attachments {
			modelAttachments = append(modelAttachments, model.Attachment{
				BillID:         nil,
				InvoiceID:      &invoice.InvoicesID,
				AttachmentName: at.Name,
				AttachmentFile: at.File,
			})
		}

		_, _, err = i.attachmentService.CreateAttachments(modelAttachments)
		if err != nil {
			trxHandle.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on creating attachments, details = %v", err),
			})
		}
	}

	var modelItemSells []model.ItemSell
	for _, item := range input.Items {
		itPurchasePpn := 0
		if item.ItemPpn {
			itPurchasePpn = 1
		}

		modelItemSells = append(modelItemSells, model.ItemSell{
			ItemID:           item.ItemId,
			InvoiceID:        invoice.InvoicesID,
			ItemSellQty:      item.ItemQty,
			ItemSellTime:     time.Now(),
			ItemSellDiscount: item.ItemDiscount,
			ItemSellPpn:      int32(itPurchasePpn),
			ItemSellUnit:     item.ItemUnit,
		})
	}
	_, _, err = i.itemSellService.CreateItemSell(modelItemSells)
	if err != nil {
		trxHandle.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on creating item sells, details = %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entity.StatusResponse{
		Status: "successfully updated invoice",
	})

}

// @Summary Delete Invoice
// @Description delete invoice that are not paid yet
// @Tags Invoice
// @Accept  json
// @Produce  json
// @Param  invoiceId path int true "invoice id"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /invoice/{invoiceId} [delete]
func (i *InvoiceController) DeleteInvoice(c *fiber.Ctx) error {
	functionName := "DeleteInvoice"

	invoiceId, err := c.ParamsInt("invoiceId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing invoice id, details = %v", err),
		})
	}

	trxHandle := c.Locals("db_trx").(*gorm.DB)
	invoice, err := i.deleteInvoiceRelatedThings(trxHandle, int32(invoiceId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     err.Error(),
		})
	}

	if err := i.invoiceService.DeleteInvoice(invoice.InvoicesID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on deleting invoice, details = %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.StatusResponse{
		Status: "success on deleting invoice",
	})
}

func (i *InvoiceController) validateInvoiceData(input struct {
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

	supplier, _, err := i.supplierService.GetSupplier(input.SupplierId)
	if err != nil {
		return time.Time{}, time.Time{}, -1, fmt.Errorf("error on getting supplier details = %v", err)
	}

	if strings.ToLower(supplier.SupplierType) != "customer" {
		return time.Time{}, time.Time{}, -1, errors.New("supplier is not a vendor")
	}

	total = 0
	for _, item := range input.Items {
		it, err := i.itemService.GetItem(item.ItemId)
		if err != nil {
			return time.Time{}, time.Time{}, -1, fmt.Errorf("error on getting item with item id '%d' details = %v", item.ItemId, err)
		}
		total += it.ItemSellPrice*float64(item.ItemQty) - it.ItemSellPrice*float64(item.ItemQty)**item.ItemDiscount/100
		if item.ItemPpn {
			total += total * 11 / 100
		}
	}

	return startDateTime, dueDateTime, total, nil
}

func (i *InvoiceController) deleteInvoiceRelatedThings(trxHandle *gorm.DB, invoiceId int32) (invoice *model.Invoice, err error) {
	invoice, _, err = i.invoiceService.GetInvoice(invoiceId)
	if err != nil {
		return nil, fmt.Errorf("error on getting invoice, details = %v", err)
	}

	if strings.ToLower(invoice.InvoiceStatus) == "sudah dibayar" {
		return nil, errors.New("invoice already been paid")
	}

	if err := i.itemSellService.DeleteItemSellsByInvoiceId(invoice.InvoicesID); err != nil {
		trxHandle.Rollback()
		return nil, fmt.Errorf("error on deleting item purchases by invoice id, details = %v", err)
	}

	if err := i.attachmentService.DeleteAttachmentByInvoiceId(invoice.InvoicesID); err != nil {
		trxHandle.Rollback()
		return nil, fmt.Errorf("error on deleting item purchases by invoice id, details = %v", err)
	}

	return invoice, nil
}
