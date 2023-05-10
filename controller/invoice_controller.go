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
		SupplierID:         input.CustomerId,
		InvoiceStartDate:   startDateTime,
		InvoiceDueDate:     dueDateTime,
		InvoiceNumber:      invoiceNumber,
		InvoiceOrderNumber: nil,
		InvoiceTitle:       input.Title,
		InvoiceSubheading:  input.Subheading,
		InvoiceLogo:        input.Logo,
		InvoiceStatus:      "MENUNGGU PEMBAYARAN",
	})
	if err != nil {
		log.Println("testt create invoice err")
		trxHandle.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on creating bill, details = %v", err),
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

	if err := trxHandle.Commit().Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on commiting transcation, details = %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entity.StatusResponse{
		Status: "successfully created invoices",
	})
}
