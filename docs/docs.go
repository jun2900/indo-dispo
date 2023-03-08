// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate_swagger = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin/login": {
            "post": {
                "description": "Login with email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin Auth"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "login user input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.LoginReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.LoginResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    }
                }
            }
        },
        "/bill": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bill"
                ],
                "summary": "Register Bill",
                "parameters": [
                    {
                        "description": "add bill request",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.AddBillReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.StatusResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    }
                }
            }
        },
        "/bill/header": {
            "get": {
                "description": "get bill overdue open and draft stats",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bill"
                ],
                "summary": "Get Bill Header",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.BillHeaderResp"
                        }
                    }
                }
            }
        },
        "/bill/{billId}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bill"
                ],
                "summary": "Get Bill Details",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "bill id",
                        "name": "billId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.BillDetailsResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bill"
                ],
                "summary": "Update Bill Status",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "bill id",
                        "name": "billId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "update bill status request (paid/cancelled)",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.BillUpdateStatusReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.StatusResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    }
                }
            }
        },
        "/bills": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bill"
                ],
                "summary": "List All Bill",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "page requested (defaults to 0)",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "number of records in a page  (defaults to 20)",
                        "name": "pagesize",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "asc / desc",
                        "name": "order",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "filter by bill status",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "filter by bill vendor",
                        "name": "vendor",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/entity.PagedResults"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "Data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/model.VSupplierBill"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    }
                }
            }
        },
        "/item": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Item"
                ],
                "summary": "Register Item",
                "parameters": [
                    {
                        "description": "add item request",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.AddItemReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.StatusResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    }
                }
            }
        },
        "/item/{itemId}": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Item"
                ],
                "summary": "Update Item",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "item id",
                        "name": "itemId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "update item request",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.UpdateItemReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.StatusResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Item"
                ],
                "summary": "Delete Item",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "item id",
                        "name": "itemId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.StatusResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "do ping",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Ping"
                ],
                "summary": "ping example",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token from login user (use Bearer in front of the jwt)",
                        "name": "x-access-token",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "pong",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/supplier": {
            "post": {
                "description": "register supplier (vendor or customer)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Supplier"
                ],
                "summary": "Register Supplier",
                "parameters": [
                    {
                        "description": "supplier request",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.SupplierAddReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.StatusResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    }
                }
            }
        },
        "/supplier/{supplierId}/items": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Supplier"
                ],
                "summary": "Get All Items By Supplier Id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "supplier id",
                        "name": "supplierId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.ListItemBySupplierResp"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    }
                }
            }
        },
        "/suppliers": {
            "get": {
                "description": "get suppliers by their type",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Supplier"
                ],
                "summary": "Get All Suplier",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "page requested (defaults to 0)",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "number of records in a page  (defaults to 20)",
                        "name": "pagesize",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "asc / desc",
                        "name": "order",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "supplier name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "supplier email",
                        "name": "email",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "supplier address",
                        "name": "address",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "supplier type (vendor or customer)",
                        "name": "supplierType",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.StatusResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrRespController"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.AddBillReq": {
            "type": "object",
            "properties": {
                "bill_account_number": {
                    "type": "string"
                },
                "bill_attachments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Attachment"
                    }
                },
                "bill_bank_name": {
                    "type": "string"
                },
                "bill_due_date": {
                    "type": "string"
                },
                "bill_items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.ItemPurchase"
                    }
                },
                "bill_notes": {
                    "type": "string"
                },
                "bill_number": {
                    "type": "string"
                },
                "bill_order_number": {
                    "type": "string"
                },
                "bill_shipping_cost": {
                    "type": "integer"
                },
                "bill_start_date": {
                    "type": "string"
                },
                "bill_type": {
                    "type": "string"
                },
                "supplier_id": {
                    "type": "integer"
                }
            }
        },
        "entity.AddItemReq": {
            "type": "object",
            "properties": {
                "item_description": {
                    "type": "string"
                },
                "item_name": {
                    "type": "string"
                },
                "item_purchase_price": {
                    "type": "integer"
                },
                "item_sell_price": {
                    "type": "integer"
                },
                "item_unit": {
                    "type": "string"
                },
                "item_wholesalers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.WholeSaler"
                    }
                },
                "supplier_id": {
                    "type": "integer"
                }
            }
        },
        "entity.Attachment": {
            "type": "object",
            "properties": {
                "attachment_file": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "attachment_name": {
                    "type": "string"
                }
            }
        },
        "entity.BillDetailsResp": {
            "type": "object",
            "properties": {
                "bill_attachments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Attachment"
                    }
                },
                "bill_due_date": {
                    "type": "string"
                },
                "bill_items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.ItemBill"
                    }
                },
                "bill_number": {
                    "type": "string"
                },
                "bill_order_number": {
                    "type": "string"
                },
                "bill_shipping_cost": {
                    "type": "integer"
                },
                "bill_start_date": {
                    "type": "string"
                },
                "bill_status": {
                    "type": "string"
                },
                "bill_subtotal": {
                    "type": "integer"
                },
                "bill_total": {
                    "type": "integer"
                },
                "bill_type": {
                    "type": "string"
                }
            }
        },
        "entity.BillHeaderResp": {
            "type": "object",
            "properties": {
                "bill_draft": {
                    "type": "integer"
                },
                "bill_open": {
                    "type": "integer"
                },
                "bill_overdue": {
                    "type": "integer"
                }
            }
        },
        "entity.BillUpdateStatusReq": {
            "type": "object",
            "properties": {
                "bill_status": {
                    "type": "string"
                }
            }
        },
        "entity.ErrRespController": {
            "type": "object",
            "properties": {
                "err_message": {
                    "type": "string"
                },
                "source_function": {
                    "type": "string"
                }
            }
        },
        "entity.ItemBill": {
            "type": "object",
            "properties": {
                "item_amount": {
                    "type": "integer"
                },
                "item_description": {
                    "type": "string"
                },
                "item_name": {
                    "type": "string"
                },
                "item_price": {
                    "type": "integer"
                },
                "item_qty": {
                    "type": "integer"
                }
            }
        },
        "entity.ItemPurchase": {
            "type": "object",
            "properties": {
                "item_discount": {
                    "type": "integer"
                },
                "item_id": {
                    "type": "integer"
                },
                "item_qty": {
                    "type": "integer"
                }
            }
        },
        "entity.ListItemBySupplierResp": {
            "type": "object",
            "properties": {
                "item_description": {
                    "type": "string"
                },
                "item_id": {
                    "type": "integer"
                },
                "item_name": {
                    "type": "string"
                },
                "item_purchase_price": {
                    "type": "integer"
                },
                "item_sell_price": {
                    "type": "integer"
                },
                "item_unit": {
                    "type": "string"
                },
                "wholesalers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.ListWholeSaler"
                    }
                }
            }
        },
        "entity.ListWholeSaler": {
            "type": "object",
            "properties": {
                "wholesaler_id": {
                    "type": "integer"
                },
                "wholesaler_price": {
                    "type": "integer"
                },
                "wholesaler_qty": {
                    "type": "integer"
                }
            }
        },
        "entity.LoginReq": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "entity.LoginResp": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                }
            }
        },
        "entity.PagedResults": {
            "type": "object",
            "properties": {
                "data": {},
                "page": {
                    "type": "integer"
                },
                "page_size": {
                    "type": "integer"
                },
                "total_records": {
                    "type": "integer"
                }
            }
        },
        "entity.StatusResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                }
            }
        },
        "entity.SupplierAddReq": {
            "type": "object",
            "properties": {
                "supplier_address": {
                    "type": "string"
                },
                "supplier_email": {
                    "type": "string"
                },
                "supplier_name": {
                    "type": "string"
                },
                "supplier_npwp": {
                    "type": "string"
                },
                "supplier_telephone": {
                    "type": "string"
                },
                "supplier_type": {
                    "type": "string"
                },
                "supplier_web": {
                    "type": "string"
                }
            }
        },
        "entity.UpdateItemReq": {
            "type": "object",
            "properties": {
                "item_description": {
                    "type": "string"
                },
                "item_name": {
                    "type": "string"
                },
                "item_purchase_price": {
                    "type": "integer"
                },
                "item_sell_price": {
                    "type": "integer"
                },
                "item_unit": {
                    "type": "string"
                },
                "item_wholesalers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.WholeSaler"
                    }
                }
            }
        },
        "entity.WholeSaler": {
            "type": "object",
            "properties": {
                "wholesaler_price": {
                    "type": "integer"
                },
                "wholesaler_qty": {
                    "type": "integer"
                }
            }
        },
        "model.VSupplierBill": {
            "type": "object",
            "properties": {
                "bill_due_date": {
                    "type": "string"
                },
                "bill_id": {
                    "type": "integer"
                },
                "bill_number": {
                    "type": "string"
                },
                "bill_order_number": {
                    "type": "string"
                },
                "bill_start_date": {
                    "type": "string"
                },
                "bill_status": {
                    "type": "string"
                },
                "bill_type": {
                    "type": "string"
                },
                "supplier_id": {
                    "type": "integer"
                },
                "supplier_name": {
                    "type": "string"
                },
                "supplier_type": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo_swagger holds exported Swagger Info so clients can modify it
var SwaggerInfo_swagger = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Test Accounting App",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate_swagger,
}

func init() {
	swag.Register(SwaggerInfo_swagger.InstanceName(), SwaggerInfo_swagger)
}
