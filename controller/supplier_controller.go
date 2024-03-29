package controller

import (
	"fmt"
	"strings"

	"github.com/jun2900/indo-dispo/entity"
	"github.com/jun2900/indo-dispo/model"
	"github.com/jun2900/indo-dispo/services"

	"github.com/gofiber/fiber/v2"
)

type SupplierController struct {
	supplierService   services.SupplierService
	wholesalerService services.WholesalerService
	itemService       services.ItemService
}

func NewSupplierController(supplierService services.SupplierService, wholesalerService services.WholesalerService, itemService services.ItemService) *SupplierController {
	return &SupplierController{
		supplierService:   supplierService,
		wholesalerService: wholesalerService,
		itemService:       itemService,
	}
}

// @Summary Get All Suplier
// @Description get suppliers by their type
// @Tags Supplier
// @Accept  json
// @Produce  json
// @Param       page     				query    int    false "page requested (defaults to 0)"
// @Param       pagesize 				query    int    false "number of records in a page  (defaults to 20)"
// @Param       order    				query    string false "asc / desc"
// @Param       name    				query    string false "supplier name"
// @Param       email    				query    string false "supplier email"
// @Param       address    				query    string false "supplier address"
// @Param  supplierType query string true "supplier type (vendor or customer)"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /suppliers [get]
func (s *SupplierController) GetAllSupplierByType(c *fiber.Ctx) error {
	functionName := "GetAllSupplierByType"

	page := c.QueryInt("page", 0)
	pagesize := c.QueryInt("pagesize", 20)
	order := c.Query("order", "")
	name := c.Query("name", "")
	email := c.Query("email", "")
	address := c.Query("address", "")

	suppliers, totalRows, err := s.supplierService.GetAllSupplierByType(strings.ToLower(c.Query("supplierType", "")), page, pagesize, order, name, email, address)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing supplier input: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.PagedResults{
		Page:         page,
		PageSize:     pagesize,
		Data:         suppliers,
		TotalRecords: int(totalRows),
	})
}

// @Summary Register Supplier
// @Description register supplier (vendor or customer)
// @Tags Supplier
// @Accept  json
// @Produce  json
// @Param  input body entity.SupplierAddReq true "supplier request"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /supplier [post]
func (s *SupplierController) RegisterSupplier(c *fiber.Ctx) error {
	var input entity.SupplierAddReq

	functionName := "RegisterSupplier"

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing supplier input: %v", err),
		})
	}

	if len(input.Name) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "name cannot be empty",
		})
	}

	_, _, err := s.supplierService.GetSupplierBySupplierName(input.Name)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "supplier already exist",
		})
	}

	if strings.ToLower(input.Type) != "vendor" && strings.ToLower(input.Type) != "customer" {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "supplier type must be vendor or customer",
		})
	}

	s.supplierService.CreateSupplier(&model.Supplier{
		SupplierName:        input.Name,
		SupplierEmail:       input.Email,
		SupplierTelephone:   input.Telephone,
		SupplierWeb:         input.Web,
		SupplierNpwp:        input.Npwp,
		SupplierAddress:     input.Address,
		SupplierType:        strings.ToUpper(input.Type),
		SupplierWhatsapp:    input.Whatsapp,
		SupplierDescription: input.Description,
		SupplierCity:        input.City,
		SupplierState:       input.State,
		SupplierZipCode:     input.ZipCode,
	})

	return c.Status(fiber.StatusCreated).JSON(entity.StatusResponse{
		Status: "supplier registered",
	})
}

// @Summary Get All Items By Supplier Id
// @Tags Supplier
// @Accept  json
// @Produce  json
// @Param  supplierId path int true "supplier id"
// @Success 200 {object} []entity.ListItemBySupplierResp
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /supplier/{supplierId}/items [get]
func (s *SupplierController) GetItemsBySupplierId(c *fiber.Ctx) error {
	functionName := "GetItemsBySupplierId"

	supplierId, err := c.ParamsInt("supplierId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing supplier id, details = %v", err),
		})
	}

	itemSuppliers, _, err := s.itemService.GetItemBySupplierId(int32(supplierId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting item supplier by supplier id, details = %v", err),
		})
	}

	var resp []entity.ListItemBySupplierResp
	for _, is := range itemSuppliers {
		item, err := s.itemService.GetItem(is.ItemID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on getting item, details = %v", err),
			})
		}

		wholeSalerRec, _, err := s.wholesalerService.GetWholesalerByItemId(item.ItemID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on getting wholesalers, details = %v", err),
			})
		}

		var wholeSalers []entity.ListWholeSaler
		if len(wholeSalerRec) > 0 {
			for _, ws := range wholeSalerRec {
				wholeSalers = append(wholeSalers, entity.ListWholeSaler{
					Id:    ws.WholesalerID,
					Qty:   ws.WholesalerQty,
					Price: ws.WholesalerPrice,
				})
			}
		}

		resp = append(resp, entity.ListItemBySupplierResp{
			Id:            is.ItemID,
			Code:          is.ItemCode,
			PurchasePrice: is.ItemPurchasePrice,
			SellPrice:     is.ItemSellPrice,
			Unit:          is.ItemUnit,
			Name:          item.ItemName,
			Description:   item.ItemDescription,
			WholeSaler:    wholeSalers,
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// @Summary Get Details Supplier
// @Tags Supplier
// @Accept  json
// @Produce  json
// @Param  supplierId path int true "supplier id"
// @Success 200 {object} model.Supplier
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /supplier/details/{supplierId} [get]
func (s *SupplierController) GetSupplierDetail(c *fiber.Ctx) error {
	functionName := "GetSupplierDetail"

	supplierId, err := c.ParamsInt("supplierId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing supplier id, details = %v", err),
		})
	}

	supplier, _, err := s.supplierService.GetSupplier(int32(supplierId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting supplier, details = %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(supplier)
}

// @Summary Get Supplier and Item Detail by Item Name
// @Tags Supplier
// @Accept  json
// @Produce  json
// @Param  input body entity.ListSupplierByItemReq true "list supplier by item req"
// @Success 200 {object} entity.ListSupplierByItemResp
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /supplier/items [post]
func (s *SupplierController) ListSupplierByItemName(c *fiber.Ctx) error {
	functionName := "ListSupplierByItemName"

	var input entity.ListSupplierByItemReq
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing item input, details = %v", err),
		})
	}

	items, _, err := s.itemService.GetAllItemByItemName(input.ItemName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting item, details = %v", err),
		})
	}

	var itemDetails []entity.ListSupplierByItemDetails

	if len(items) > 0 {
		for _, item := range items {
			supplier, _, err := s.supplierService.GetSupplier(item.SupplierID)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
					SourceFunction: functionName,
					ErrMessage:     fmt.Sprintf("error on getting supplier by item, details = %v", err),
				})
			}
			itemDetails = append(itemDetails, entity.ListSupplierByItemDetails{
				SupplierName:  supplier.SupplierName,
				Description:   supplier.SupplierDescription,
				PurchasePrice: item.ItemPurchasePrice,
				SellPrice:     item.ItemSellPrice,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(entity.ListSupplierByItemResp{
		ItemName: input.ItemName,
		Data:     itemDetails,
	})
}
