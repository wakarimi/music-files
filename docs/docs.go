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
            "name": "Dmitry Kolesnikov (Zalimannard)",
            "email": "zalimannard@mail.ru"
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
        "/dirs": {
            "post": {
                "description": "Adds a new directory to the database for tracking",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Directories"
                ],
                "summary": "Add a new tracked directory",
                "parameters": [
                    {
                        "description": "Directory Data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dir_handler.trackRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dir_handler.trackResponse"
                        }
                    },
                    "400": {
                        "description": "Failed to decode request",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "404": {
                        "description": "Directory not found",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "409": {
                        "description": "Directory already tracked",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dir_handler.trackRequest": {
            "type": "object",
            "properties": {
                "path": {
                    "description": "Path to the directory on disk",
                    "type": "string"
                }
            }
        },
        "dir_handler.trackResponse": {
            "type": "object",
            "properties": {
                "dirId": {
                    "description": "Unique identifier of the directory in the database",
                    "type": "integer"
                },
                "name": {
                    "description": "Name of the directory",
                    "type": "string"
                }
            }
        },
        "responses.Error": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.4",
	Host:             "localhost:8022",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Wakarimi Music Files API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
