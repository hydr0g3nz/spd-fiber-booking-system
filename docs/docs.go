// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "support@example.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/bookings": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get a list of all bookings with optional sorting and filtering",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Get all bookings",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Sort by field (price or date)",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Filter high-value bookings (price \u003e 50,000)",
                        "name": "high-value",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of bookings",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Booking"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create a new booking with the provided details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Create a new booking",
                "parameters": [
                    {
                        "description": "Booking Information",
                        "name": "booking",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CreateBookingRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created booking",
                        "schema": {
                            "$ref": "#/definitions/models.Booking"
                        }
                    },
                    "400": {
                        "description": "Invalid request parameters",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/bookings/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get detailed information about a specific booking",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Get a booking by ID",
                "parameters": [
                    {
                        "minimum": 1,
                        "type": "integer",
                        "description": "Booking ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Booking details",
                        "schema": {
                            "$ref": "#/definitions/models.Booking"
                        }
                    },
                    "400": {
                        "description": "Invalid booking ID format",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Booking not found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Cancel an existing booking by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Cancel a booking",
                "parameters": [
                    {
                        "minimum": 1,
                        "type": "integer",
                        "description": "Booking ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Canceled booking details",
                        "schema": {
                            "$ref": "#/definitions/models.Booking"
                        }
                    },
                    "400": {
                        "description": "Invalid booking ID or cannot cancel",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Booking not found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.CreateBookingRequest": {
            "description": "Request payload for creating a new booking",
            "type": "object",
            "required": [
                "price",
                "service_id",
                "user_id"
            ],
            "properties": {
                "price": {
                    "type": "number",
                    "example": 30000
                },
                "service_id": {
                    "type": "integer",
                    "example": 456
                },
                "user_id": {
                    "type": "integer",
                    "example": 123
                }
            }
        },
        "models.Booking": {
            "description": "Booking entity representing a customer's service booking",
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2024-03-11T12:00:00Z"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "price": {
                    "type": "number",
                    "example": 30000
                },
                "service_id": {
                    "type": "integer",
                    "example": 456
                },
                "status": {
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.BookingStatus"
                        }
                    ],
                    "example": "pending"
                },
                "updated_at": {
                    "type": "string",
                    "format": "date-time",
                    "example": "2024-03-11T12:00:00Z"
                },
                "user_id": {
                    "type": "integer",
                    "example": 123
                }
            }
        },
        "models.BookingStatus": {
            "type": "string",
            "enum": [
                "pending",
                "confirmed",
                "rejected",
                "canceled"
            ],
            "x-enum-varnames": [
                "BookingStatusPending",
                "BookingStatusConfirmed",
                "BookingStatusRejected",
                "BookingStatusCanceled"
            ]
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "API key authentication",
            "type": "apiKey",
            "name": "X-API-Key",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Fiber Booking System API",
	Description:      "A booking system API built with Fiber framework",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
