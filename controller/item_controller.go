package controller

import (
	"fmt"
	"log"

	"github.com/jun2900/indo-dispo/dal"
	"github.com/jun2900/indo-dispo/entity"
	"github.com/jun2900/indo-dispo/model"
	"github.com/jun2900/indo-dispo/services"

	"github.com/gofiber/fiber/v2"
)

type ItemController struct {
	itemService       services.ItemService
	wholesalerService services.WholesalerService
}

func NewItemController(itemService services.ItemService, wholesalerService services.WholesalerService) *ItemController {
	return &ItemController{
		itemService:       itemService,
		wholesalerService: wholesalerService,
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

	if len(input.Code) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "item code cannot be empty",
		})
	}

	item, _ := i.itemService.GetItemByItemName(input.Name)
	if item != nil {
		_, err := i.itemService.GetItemWithItemIdAndSupplierId(item.ItemID, input.SupplierId)
		if err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     "a supplier cannot have the same item",
			})
		}
	}
	item, _, err := i.itemService.CreateItem(&model.Item{
		SupplierID:        input.SupplierId,
		ItemCode:          input.Code,
		ItemName:          input.Name,
		ItemDescription:   input.Description,
		ItemPurchasePrice: input.PurchasePrice,
		ItemSellPrice:     input.SellPrice,
		ItemUnit:          input.Unit,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on registering item, details = %v", err),
		})
	}

	if len(input.WholeSalers) > 0 {
		var wholeSalers []model.Wholesaler
		for _, ws := range input.WholeSalers {
			if ws.Price <= 0 {
				return c.Status(fiber.StatusInternalServerError).JSON(entity.ErrRespController{
					SourceFunction: functionName,
					ErrMessage:     "price cannot below or equals to zero",
				})
			}
			if ws.Qty <= 0 {
				return c.Status(fiber.StatusInternalServerError).JSON(entity.ErrRespController{
					SourceFunction: functionName,
					ErrMessage:     "quantity cannot below or equals to zero",
				})
			}
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

	return c.Status(fiber.StatusCreated).JSON(entity.StatusResponse{
		Status: "successfully created item",
	})
}

// @Summary Update Item
// @Tags Item
// @Accept  json
// @Produce  json
// @Param  itemId path int true "item id"
// @Param  input body entity.UpdateItemReq true "update item request"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /item/{itemId} [put]
func (i *ItemController) UpdateItem(c *fiber.Ctx) error {
	functionName := "UpdateItem"

	var input entity.UpdateItemReq

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing item input, details = %v", err),
		})
	}

	itemId, err := c.ParamsInt("itemId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing params item id, details = %v", err),
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

	if len(input.Description) < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "item description cannot be empty",
		})
	}

	_, _, err = i.itemService.UpdateItem(int32(itemId), &model.Item{
		ItemName:          input.Name,
		ItemDescription:   &input.Description,
		ItemPurchasePrice: input.PurchasePrice,
		ItemSellPrice:     input.SellPrice,
		ItemUnit:          input.Unit,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on updating item, details = %v", err),
		})
	}

	if len(input.WholeSalers) > 0 {
		_, _, err := i.wholesalerService.DeleteWholesalerByItemId(int32(itemId))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on deleting wholesaler by item id, details = %v", err),
			})
		}

		var modelWholeSalers []model.Wholesaler
		for _, ws := range input.WholeSalers {
			modelWholeSalers = append(modelWholeSalers, model.Wholesaler{
				ItemID:          int32(itemId),
				WholesalerQty:   ws.Qty,
				WholesalerPrice: ws.Price,
			})
		}
		_, _, err = i.wholesalerService.CreateWholesaler(modelWholeSalers)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on updating wholesaler, details = %v", err),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(entity.StatusResponse{
		Status: "successfully updated item",
	})
}

// @Summary Delete Item
// @Tags Item
// @Accept  json
// @Produce  json
// @Param  itemId path int true "item id"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /item/{itemId} [delete]
func (i *ItemController) DeleteItem(c *fiber.Ctx) error {
	functionName := "DeleteItem"

	itemId, err := c.ParamsInt("itemId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing params item id, details = %v", err),
		})
	}

	_, _, err = i.itemService.DeleteItem(int32(itemId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on deleting item, details = %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.StatusResponse{
		Status: "successfully deleted item",
	})
}

// @Summary Get Item by Item Code
// @Tags Item
// @Accept  json
// @Produce  json
// @Param  itemCode path string true "item code"
// @Success 200 {object} model.Item
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /item/{itemCode} [get]
func (i *ItemController) GetItemByCodeItem(c *fiber.Ctx) error {
	functionName := "GetItemByCodeItem"
	log.Println("get item by code item")

	item, err := dal.Item.Where(dal.Item.ItemCode.Eq(c.Params("itemCode"))).First()
	if err != nil {
		log.Printf("error on retrieving item, details = %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("cannot retrieve item, details = %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(item)
}
