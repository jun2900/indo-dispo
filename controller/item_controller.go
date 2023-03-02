package controller

import (
	"anara/entity"
	"anara/model"
	"anara/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ItemController struct {
	itemService         services.ItemService
	itemSupplierService services.ItemSupplierService
	wholesalerService   services.WholesalerService
}

func NewItemController(itemService services.ItemService, itemSupplierService services.ItemSupplierService, wholesalerService services.WholesalerService) *ItemController {
	return &ItemController{
		itemService:         itemService,
		itemSupplierService: itemSupplierService,
		wholesalerService:   wholesalerService,
	}
}

// @Summary Register Item
// @Tags Item
// @Accept  json
// @Produce  json
// @Param  input body entity.AddItemReq true "add item request"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /item [post]
func (i *ItemController) RegisterItem(c *fiber.Ctx) error {
	var input entity.AddItemReq

	functionName := "RegisterItem"

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing item input, details = %v", err),
		})
	}

	if len(input.Name) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "item name cannot be empty",
		})
	}

	if len(input.Unit) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "item unit cannot be empty",
		})
	}

	item, err := i.itemService.GetItemByItemName(input.Name)
	if item != nil {
		_, _, err := i.itemSupplierService.GetItemSupplierByItemIdAndSupplierId(item.ItemID, input.SupplierId)
		if err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     "a supplier cannot have the same item",
			})
		}
	} else if err != nil {
		item, _, err = i.itemService.CreateItem(&model.Item{
			ItemName:        input.Name,
			ItemDescription: input.Description,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on registering item, details = %v", err),
			})
		}
	}

	if len(input.WholeSalers) > 0 {
		var wholeSalers []model.Wholesaler
		for _, ws := range input.WholeSalers {
			wholeSalers = append(wholeSalers, model.Wholesaler{
				ItemID:          item.ItemID,
				WholesalerQty:   ws.Qty,
				WholesalerPrice: ws.Price,
			})
		}
		_, _, err := i.wholesalerService.CreateWholesaler(wholeSalers)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on creating wholesaler, details = %v", err),
			})
		}
	}

	_, _, err = i.itemSupplierService.CreateItemSupplier(&model.ItemSupplier{
		ItemID:                    item.ItemID,
		SupplierID:                input.SupplierId,
		ItemSupplierPurchasePrice: input.PurchasePrice,
		ItemSupplierSellPrice:     input.SellPrice,
		ItemSupplierUnit:          input.Unit,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on registering item supplier, details = %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entity.StatusResponse{
		Status: "successfully created item",
	})
}
