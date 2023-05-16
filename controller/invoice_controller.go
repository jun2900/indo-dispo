package controller

import (
	"crypto/rand"
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

	randNum, _ := rand.Int(rand.Reader, big.NewInt(9000))
	randNum = randNum.Add(randNum, big.NewInt(1000))

	invoiceNumber := fmt.Sprintf("IDS/%s/%d", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""), randNum)

	trxHandle := c.Locals("db_trx").(*gorm.DB)

	invoice, _, err := i.invoiceService.WithTrx(trxHandle).CreateInvoice(&model.Invoice{
		SupplierID:        input.CustomerId,
		InvoiceStartDate:  startDateTime,
		InvoiceDueDate:    dueDateTime,
		InvoiceNumber:     invoiceNumber,
		InvoiceTitle:      input.Title,
		InvoiceSubheading: input.Subheading,
		InvoiceLogo:       input.Logo,
		InvoiceTotal:      total,
		InvoiceStatus:     "MENUNGGU PEMBAYARAN",
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
			ErrMessage:     fmt.Sprintf("error on creating item purchases, details = %v", err),
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
