// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
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
        "/couriers/": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "courier-controller"
                ],
                "summary": "Получение профилей курьеров",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Максимальное количество курьеров в выдаче. Если параметр не передан, то значение по умолчанию равно 1.",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Количество курьеров, которое нужно пропустить для отображения текущей страницы. Если параметр не передан, то значение по умолчанию равно 0.",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.GetCouriersResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.BadRequestResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "courier-controller"
                ],
                "summary": "Создание профилей курьеров",
                "parameters": [
                    {
                        "description": "Couriers",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CreateCourierRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.CouriersCreateResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.BadRequestResponse"
                        }
                    }
                }
            }
        },
        "/couriers/assignments": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "courier-controller"
                ],
                "summary": "список распределенных заказов",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор курьера для получения списка распредленных заказов. Если не указан, возвращаются данные по всем курьерам.",
                        "name": "courier_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Дата распределения заказов. Если не указана, то используется текущий день",
                        "name": "date",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.OrderAssignResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.BadRequestResponse"
                        }
                    }
                }
            }
        },
        "/couriers/meta-info/{courier_id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "courier-controller"
                ],
                "summary": "Получение meta-информации о курьере.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Courier identifier",
                        "name": "courier_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Максимальное количество курьеров в выдаче. Если параметр не передан, то значение по умолчанию равно 1.",
                        "name": "startDate",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Количество курьеров, которое нужно пропустить для отображения текущей страницы. Если параметр не передан, то значение по умолчанию равно 0.",
                        "name": "endDate",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.GetCourierMetaInfoResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.BadRequestResponse"
                        }
                    }
                }
            }
        },
        "/couriers/{courier_id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "courier-controller"
                ],
                "summary": "Получение профиля курьера",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Courier identifier",
                        "name": "courier_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.CourierDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.BadRequestResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/model.BadRequestResponse"
                        }
                    }
                }
            }
        },
        "/orders/": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order-controller"
                ],
                "summary": "Получение заказов",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Максимальное количество заказов в выдаче. Если параметр не передан, то значение по умолчанию равно 1.",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Количество заказов, которое нужно пропустить для отображения текущей страницы. Если параметр не передан, то значение по умолчанию равно 0.",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.OrderDTO"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.BadRequestResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order-controller"
                ],
                "summary": "Создание заказов",
                "parameters": [
                    {
                        "description": "Orders",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CreateOrderRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.OrderDTO"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.BadRequestResponse"
                        }
                    }
                }
            }
        },
        "/orders/assign": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order-controller"
                ],
                "summary": "Распределение заказов по курьерам",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Дата распределения заказов. Если не указана, то используется текущий день",
                        "name": "date",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.OrderAssignResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.BadRequestResponse"
                        }
                    }
                }
            }
        },
        "/orders/complete": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order-controller"
                ],
                "summary": "Завершение заказов",
                "parameters": [
                    {
                        "description": "Orders",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CompleteOrderRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.OrderDTO"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.BadRequestResponse"
                        }
                    }
                }
            }
        },
        "/orders/{order_id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order-controller"
                ],
                "summary": "Получение информации о заказе",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Order identifier",
                        "name": "order_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.OrderDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.BadRequestResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/model.BadRequestResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.BadRequestResponse": {
            "type": "object"
        },
        "model.CompleteOrder": {
            "type": "object",
            "required": [
                "complete_time",
                "courier_id",
                "order_id"
            ],
            "properties": {
                "complete_time": {
                    "type": "string"
                },
                "courier_id": {
                    "type": "integer"
                },
                "order_id": {
                    "type": "integer"
                }
            }
        },
        "model.CompleteOrderRequest": {
            "type": "object",
            "required": [
                "complete_info"
            ],
            "properties": {
                "complete_info": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.CompleteOrder"
                    }
                }
            }
        },
        "model.CourierDTO": {
            "type": "object",
            "required": [
                "courier_id",
                "courier_type",
                "regions",
                "working_hours"
            ],
            "properties": {
                "courier_id": {
                    "type": "integer",
                    "example": 2
                },
                "courier_type": {
                    "type": "string",
                    "enum": [
                        "FOOT",
                        "BIKE",
                        "AUTO"
                    ],
                    "example": "AUTO"
                },
                "regions": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    },
                    "example": [
                        1,
                        2,
                        3
                    ]
                },
                "working_hours": {
                    "description": "WorkingHours is string slice of strings that represents time interval.\n\nString must be in HH:MM-HH:MM format where HH is hour (integer 0-23) and MM is minutes (integer 0-59).",
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "12:00-23:00",
                        "14:30-15:30"
                    ]
                }
            }
        },
        "model.CourierGroupOrders": {
            "type": "object",
            "properties": {
                "courier_id": {
                    "type": "integer"
                },
                "orders": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.GroupOrders"
                    }
                }
            }
        },
        "model.CouriersCreateResponse": {
            "type": "object",
            "properties": {
                "couriers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.CourierDTO"
                    }
                }
            }
        },
        "model.CreateCourierDTO": {
            "type": "object",
            "required": [
                "courier_type",
                "regions",
                "working_hours"
            ],
            "properties": {
                "courier_type": {
                    "type": "string",
                    "enum": [
                        "FOOT",
                        "BIKE",
                        "AUTO"
                    ],
                    "example": "AUTO"
                },
                "regions": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    },
                    "example": [
                        1,
                        2,
                        3
                    ]
                },
                "working_hours": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "12:00-23:00",
                        "14:30-15:30"
                    ]
                }
            }
        },
        "model.CreateCourierRequest": {
            "type": "object",
            "required": [
                "couriers"
            ],
            "properties": {
                "couriers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.CreateCourierDTO"
                    }
                }
            }
        },
        "model.CreateOrderDTO": {
            "type": "object",
            "required": [
                "cost",
                "delivery_hours",
                "regions",
                "weight"
            ],
            "properties": {
                "cost": {
                    "type": "integer"
                },
                "delivery_hours": {
                    "description": "DeliveryHours is string slice of strings that represents time interval.\n\nString must be in HH:MM-HH:MM format where HH is hour (integer 0-23) and MM is minutes (integer 0-59).",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "regions": {
                    "type": "integer"
                },
                "weight": {
                    "type": "integer"
                }
            }
        },
        "model.CreateOrderRequest": {
            "type": "object",
            "required": [
                "orders"
            ],
            "properties": {
                "orders": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.CreateOrderDTO"
                    }
                }
            }
        },
        "model.GetCourierMetaInfoResponse": {
            "type": "object",
            "required": [
                "courier_id",
                "courier_type",
                "regions",
                "working_hours"
            ],
            "properties": {
                "courier_id": {
                    "type": "integer",
                    "example": 1
                },
                "courier_type": {
                    "type": "string",
                    "enum": [
                        "FOOT",
                        "BIKE",
                        "AUTO"
                    ],
                    "example": "AUTO"
                },
                "earnings": {
                    "type": "integer"
                },
                "rating": {
                    "type": "integer"
                },
                "regions": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    },
                    "example": [
                        1,
                        3,
                        6
                    ]
                },
                "working_hours": {
                    "description": "WorkingHours is string slice of strings that represents time interval.\n\nString must be in HH:MM-HH:MM format where HH is hour (integer 0-23) and MM is minutes (integer 0-59).",
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "12:00-23:00",
                        "14:30-15:30"
                    ]
                }
            }
        },
        "model.GetCouriersResponse": {
            "type": "object",
            "properties": {
                "couriers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.CourierDTO"
                    }
                },
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                }
            }
        },
        "model.GroupOrders": {
            "type": "object",
            "required": [
                "orders"
            ],
            "properties": {
                "group_order_id": {
                    "type": "integer"
                },
                "orders": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.OrderDTO"
                    }
                }
            }
        },
        "model.OrderAssignResponse": {
            "type": "object",
            "properties": {
                "couriers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.CourierGroupOrders"
                    }
                },
                "date": {
                    "type": "string"
                }
            }
        },
        "model.OrderDTO": {
            "type": "object",
            "required": [
                "completed_time",
                "cost",
                "delivery_hours",
                "order_id",
                "regions",
                "weight"
            ],
            "properties": {
                "completed_time": {
                    "type": "string"
                },
                "cost": {
                    "type": "integer"
                },
                "delivery_hours": {
                    "description": "DeliveryHours is string slice of strings that represents time interval.\n\nString must be in HH:MM-HH:MM format where HH is hour (integer 0-23) and MM is minutes (integer 0-59).",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "order_id": {
                    "type": "integer"
                },
                "regions": {
                    "type": "integer"
                },
                "weight": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Yandex Lavka",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
