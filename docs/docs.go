// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "url": "http://localhost/support"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/authorize": {
            "post": {
                "description": "Retrieve an JWT access token for supplied credentials",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Authorize",
                "operationId": "authorize",
                "parameters": [
                    {
                        "description": "Credentials",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.Credentials"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.TokenCredentials"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized (invalid credentials)",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/invoices": {
            "get": {
                "description": "List all Invoices",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "List Invoices",
                "operationId": "all_invoices",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "last seq id",
                        "name": "last_id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "number of invoices per page",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/invoice.Invoice"
                            }
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
            },
            "post": {
                "description": "Create an Invoice",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "Create Invoice",
                "operationId": "create_invoice",
                "parameters": [
                    {
                        "description": "Create Invoice Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/invoice.CreateInvoiceRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/invoice.Invoice"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized (invalid credentials)",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/invoices/user/{id}": {
            "get": {
                "description": "Get all Invoices for a specific User",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "Get Invoices For User",
                "operationId": "get_all_for_user_invoice",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/invoice.Invoice"
                            }
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
        "/invoices/{id}": {
            "get": {
                "description": "Get a specific Invoice",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "Get Invoice",
                "operationId": "get_invoice",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "return with Items (if there are any attached to the invoice)",
                        "name": "withItems",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/invoice.Invoice"
                            }
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
            },
            "put": {
                "description": "Update an Invoice",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "Update Invoice",
                "operationId": "update_invoice",
                "parameters": [
                    {
                        "description": "Update Invoice Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/invoice.UpdateInvoiceRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/invoice.Invoice"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized (invalid credentials)",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "description": "Remove a specific Invoice",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "Delete Invoice",
                "operationId": "delete_invoice",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/commons.DeleteResult"
                            }
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
        "/invoices/{id}/items": {
            "post": {
                "description": "Add Items to an Invoice",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "Add Items to Invoice",
                "operationId": "add_items_to_invoice",
                "parameters": [
                    {
                        "description": "Add Items to Invoice Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/invoice.ItemsToInvoiceRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/invoice.ItemsToInvoiceResponse"
                            }
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
        "/invoices/{id}/items/{itemId}": {
            "delete": {
                "description": "Remove a specific Item from a specific Invoice",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "invoice"
                ],
                "summary": "Remove Item From Invoice",
                "operationId": "remove_item_from_invoice",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/commons.DeleteResult"
                            }
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
        "/items": {
            "get": {
                "description": "List all Items",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "List Items",
                "operationId": "all_items",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "last seq id",
                        "name": "last_id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "number of items per page",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/item.Item"
                            }
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
            },
            "post": {
                "description": "Create an Item",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Create Item",
                "operationId": "create_item",
                "parameters": [
                    {
                        "description": "Create Item Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/item.CreateItemRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/item.Item"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized (invalid credentials)",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/items/{id}": {
            "get": {
                "description": "Get a specific Item",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Get Item",
                "operationId": "get_item",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/item.Item"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
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
            },
            "put": {
                "description": "Update an Item",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Update Item",
                "operationId": "update_item",
                "parameters": [
                    {
                        "description": "Update Item Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/item.UpdateItemRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/item.Item"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized (invalid credentials)",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "description": "Remove a specific Item",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "item"
                ],
                "summary": "Delete Item",
                "operationId": "delete_item",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/commons.DeleteResult"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
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
        "/persons": {
            "get": {
                "description": "List all Persons",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "person"
                ],
                "summary": "List Persons",
                "operationId": "all_persons",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "last seq id",
                        "name": "last_id",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "number of persons per page",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/person.Person"
                            }
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
            },
            "post": {
                "description": "Create an Person",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "person"
                ],
                "summary": "Create Person",
                "operationId": "create_person",
                "parameters": [
                    {
                        "description": "Create Person Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/person.CreatePersonRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/person.Person"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized (invalid credentials)",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/persons/{id}": {
            "get": {
                "description": "Get a specific Person",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "person"
                ],
                "summary": "Get Person",
                "operationId": "get_person",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/person.Person"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
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
            },
            "put": {
                "description": "Update an Person",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "person"
                ],
                "summary": "Update Person",
                "operationId": "update_person",
                "parameters": [
                    {
                        "description": "Update Person Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/person.UpdatePersonRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/person.Person"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized (invalid credentials)",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "description": "Remove a specific Person",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "person"
                ],
                "summary": "Delete Person",
                "operationId": "delete_person",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/commons.DeleteResult"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
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
        }
    },
    "definitions": {
        "auth.Credentials": {
            "type": "object",
            "properties": {
                "client_id": {
                    "type": "string"
                },
                "client_secret": {
                    "type": "string"
                }
            }
        },
        "commons.AuditInfo": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "created_by": {
                    "type": "string"
                },
                "last_change_by": {
                    "type": "string"
                },
                "last_update": {
                    "type": "string"
                }
            }
        },
        "commons.DeleteResult": {
            "type": "object",
            "properties": {
                "deleted": {
                    "type": "boolean"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "handlers.TokenCredentials": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "integer"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "invoice.CreateInvoiceRequest": {
            "type": "object",
            "properties": {
                "created_by": {
                    "type": "string"
                },
                "paid": {
                    "type": "boolean"
                },
                "total": {
                    "type": "number"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "invoice.Invoice": {
            "type": "object",
            "properties": {
                "audit_info": {
                    "$ref": "#/definitions/commons.AuditInfo"
                },
                "id": {
                    "type": "string"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/item.Item"
                    }
                },
                "paid": {
                    "type": "boolean"
                },
                "seq": {
                    "type": "integer"
                },
                "total": {
                    "type": "number"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "invoice.ItemsToInvoiceRequest": {
            "type": "object",
            "properties": {
                "invoice_id": {
                    "type": "string"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "invoice.ItemsToInvoiceResponse": {
            "type": "object",
            "properties": {
                "invoice_id": {
                    "type": "string"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "invoice.UpdateInvoiceRequest": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "last_changed_by": {
                    "type": "string"
                },
                "paid": {
                    "type": "boolean"
                },
                "total": {
                    "type": "number"
                }
            }
        },
        "item.CreateItemRequest": {
            "type": "object",
            "properties": {
                "created_by": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "unit_price": {
                    "type": "number"
                }
            }
        },
        "item.Item": {
            "type": "object",
            "properties": {
                "audit_info": {
                    "$ref": "#/definitions/commons.AuditInfo"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "seq": {
                    "type": "integer"
                },
                "unit_price": {
                    "type": "number"
                }
            }
        },
        "item.UpdateItemRequest": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "last_changed_by": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "unit_price": {
                    "type": "number"
                }
            }
        },
        "person.CreatePersonRequest": {
            "type": "object",
            "properties": {
                "created_by": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "person.Person": {
            "type": "object",
            "properties": {
                "audit_info": {
                    "$ref": "#/definitions/commons.AuditInfo"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "seq": {
                    "type": "integer"
                }
            }
        },
        "person.UpdatePersonRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "last_changed_by": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Inventory Service API",
	Description:      "This is an implementation of the imventory-service in Go.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
