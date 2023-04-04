package controller

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/jun2900/indo-dispo/entity"
	"github.com/jun2900/indo-dispo/model"
	"github.com/jun2900/indo-dispo/services"

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

// @Summary Get All Balance Logs
// @Tags Balance
// @Accept  json
// @Produce  json
// @Param       page     				query    int    false "page requested (defaults to 0)"
// @Param       pagesize 				query    int    false "number of records in a page  (defaults to 20)"
// @Param       order    				query    string false "asc / desc"
// @Param       log_start_time    	query    string false "balance log time start (lower bound)"
// @Param       log_end_time    	query    string false "balance log time end (upper bound)"
// @Success 200 {object} entity.PagedResults{Data=[]model.BalanceLog}
// @Failure 400 {object} entity.ErrRespController
// @Failure 500 {object} entity.ErrRespController
// @Router /balance/logs [get]
func (b *BalanceController) GetAllBalanceLog(c *fiber.Ctx) error {
	functionName := "GetAllBalanceLog"

	page := c.QueryInt("page", 0)
	pagesize := c.QueryInt("pagesize", 20)

	order := c.Query("order", "")
	log_start_time := c.Query("log_start_time", "")
	log_end_time := c.Query("log_end_time", "")

	var b_log_start_time time.Time
	var b_log_end_time time.Time

	var err error
	if log_start_time != "" {
		b_log_start_time, err = time.Parse(layoutTime, log_start_time)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on parsing time, details = %v", err),
			})
		}
	}
	if log_end_time != "" {
		b_log_end_time, err = time.Parse(layoutTime, log_end_time)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on parsing time, details = %v", err),
			})
		}
	}

	balanceLogs, totalRows, err := b.balanceLogService.GetAllBalanceLog(page, pagesize, order, b_log_start_time, b_log_end_time)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     fmt.Sprintf("error on getting balance logs, details = %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.PagedResults{
		Page:         page,
		PageSize:     pagesize,
		Data:         balanceLogs,
		TotalRecords: int(totalRows),
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

	var attach *[]byte

	if input.Attachment != nil {
		at, err := base64.StdEncoding.DecodeString(*input.Attachment)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
				SourceFunction: functionName,
				ErrMessage:     fmt.Sprintf("error on decoding attachment, details = %v", err),
			})
		}
		attach = &at
	}

	_, _, err = b.balanceLogService.CreateBalanceLog(&model.BalanceLog{
		BalanceLogAmount:     input.Amount,
		BalanceLogNotes:      input.Notes,
		BalanceLogAttachment: attach,
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
