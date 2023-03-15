package controller

import (
	"anara/entity"
	"anara/model"
	"anara/services"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BalanceController struct {
	balanceService    services.BalanceService
	billService       services.BillService
	balanceLogService services.BalanceLogService
}

func NewBalanceController(balanceService services.BalanceService, billService services.BillService, balanceLogService services.BalanceLogService) *BalanceController {
	return &BalanceController{
		balanceService:    balanceService,
		billService:       billService,
		balanceLogService: balanceLogService,
	}
}

// @Summary Get Balance Header
// @Description get balance in net amount (balance - paid bills)
// @Tags Balance
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.GetBalanceAmountResp
// @Failure 400 {object} entity.ErrRespController
// @Router /balance/header [get]
func (b *BalanceController) GetNetBalanceAmount(c *fiber.Ctx) error {
	functionName := "GetNetBalanceAmount"

	balance, _, err := b.balanceService.GetBalance(1)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting balance, details = %v", err),
		})
	}

	total := balance.BalanceAmount - b.billService.GetAllPaidAndOperasionalBillTotal()

	return c.Status(fiber.StatusOK).JSON(entity.GetBalanceAmountResp{
		Amount: total,
	})
}

// @Summary Add Balance Amount
// @Tags Balance
// @Accept  json
// @Produce  json
// @Param  input body entity.UpdateBalanceAmountReq true "add balance req"
// @Success 200 {object} entity.StatusResponse
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /balance [put]
func (b *BalanceController) AddBalanceAmount(c *fiber.Ctx) error {
	functionName := "AddBalanceAmount"

	var input entity.UpdateBalanceAmountReq

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing item input, details = %v", err),
		})
	}

	dateAddedTime, err := time.Parse(layoutTime, input.DateAdded)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on parsing date added, details = %v", err),
		})
	}

	balance, _, err := b.balanceService.GetBalance(1)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting balance, details = %v", err),
		})
	}

	total := balance.BalanceAmount + input.Amount

	_, _, err = b.balanceService.UpdateBalanceAmount(1, total)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on updating balance amount, details = %v", err),
		})
	}

	_, _, err = b.balanceLogService.CreateBalanceLog(&model.BalanceLog{
		BalanceLogAmount:     input.Amount,
		BalanceLogNotes:      input.Notes,
		BalanceLogAttachment: input.Attachment,
		BalanceLogTimeAdded:  dateAddedTime,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on adding balance log, details = %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entity.StatusResponse{
		Status: "successfully added balance amount",
	})
}
