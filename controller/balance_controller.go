package controller

import (
	"anara/entity"
	"anara/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type BalanceController struct {
	balanceService services.BalanceService
	billService    services.BillService
}

func NewBalanceController(balanceService services.BalanceService, billService services.BillService) *BalanceController {
	return &BalanceController{
		balanceService: balanceService,
		billService:    billService,
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
	functionName := ""

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
