package controller

import (
	"anara/entity"
	"anara/model"
	"anara/services"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type SupplierController struct {
	supplierService services.SupplierService
}

func NewSupplierController(supplierService services.SupplierService) *SupplierController {
	return &SupplierController{
		supplierService: supplierService,
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

	suppliers, totalRows, err := s.supplierService.GetAllSupplierByType(c.Query("supplierType", ""), page, pagesize, order, name, email, address)
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

	if strings.ToLower(input.Type) != "vendor" && strings.ToLower(input.Type) != "customer" {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ErrRespController{
			SourceFunction: functionName,
			ErrMessage:     "supplier type must be vendor or customer",
		})
	}

	s.supplierService.CreateSupplier(&model.Supplier{
		SupplierName:      input.Name,
		SupplierEmail:     input.Email,
		SupplierTelephone: input.Telephone,
		SupplierWeb:       input.Web,
		SupplierNpwp:      input.Npwp,
		SupplierAddress:   input.Address,
		SupplierType:      input.Type,
	})

	return c.Status(fiber.StatusCreated).JSON(entity.StatusResponse{
		Status: "supplier registered",
	})
}
