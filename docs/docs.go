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
        "/covers/{coverId}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Covers"
                ],
                "summary": "Fetch data about a cover by its ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Cover Identifier",
                        "name": "coverId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully fetched cover data",
                        "schema": {
                            "$ref": "#/definitions/cover.readResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid coverId format",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch cover",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/covers/{coverId}/download": {
            "get": {
                "produces": [
                    "application/octet-stream"
                ],
                "tags": [
                    "Covers"
                ],
                "summary": "Download a cover by its ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Cover Identifier",
                        "name": "coverId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully downloaded cover file",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Invalid coverId format",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch cover",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/dirs": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Directories"
                ],
                "summary": "Get all added directories for scanning",
                "responses": {
                    "200": {
                        "description": "Successfully fetched all directories",
                        "schema": {
                            "$ref": "#/definitions/directory.readAllResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch directories",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Adds a new directory",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Directories"
                ],
                "summary": "Create a directory",
                "parameters": [
                    {
                        "description": "Details for the new directory",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/directory.createRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully created directory",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Failed to encode request or Invalid input",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to create directory",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/dirs/scan-all": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Directories"
                ],
                "summary": "Scan all directories",
                "responses": {
                    "200": {
                        "description": "Directories scanned successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to scan directories",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/dirs/{dirId}": {
            "delete": {
                "description": "Deletes a directory using its unique identifier",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Directories"
                ],
                "summary": "Delete a directory",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Directory ID",
                        "name": "dirId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Successfully deleted directory",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid dirId format",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to delete directory",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/dirs/{dirId}/scan": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Directories"
                ],
                "summary": "Scan a directory by its ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Directory Identifier",
                        "name": "dirId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Directory scanned successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid dirId format",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to scan directory",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/tracks": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tracks"
                ],
                "summary": "Retrieve all tracks",
                "responses": {
                    "200": {
                        "description": "Successfully retrieved all tracks",
                        "schema": {
                            "$ref": "#/definitions/track.readAllResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch all tracks",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/tracks/{trackId}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tracks"
                ],
                "summary": "Retrieve a single track by its ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Track ID",
                        "name": "trackId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully retrieved the track",
                        "schema": {
                            "$ref": "#/definitions/track.readResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid trackId format",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch track",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/tracks/{trackId}/download": {
            "get": {
                "produces": [
                    "application/octet-stream"
                ],
                "tags": [
                    "Tracks"
                ],
                "summary": "Download a track by its ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Track Identifier",
                        "name": "trackId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully downloaded track file",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Invalid trackId format",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Failed to fetch track",
                        "schema": {
                            "$ref": "#/definitions/types.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "cover.readResponse": {
            "type": "object",
            "properties": {
                "coverId": {
                    "type": "integer"
                },
                "format": {
                    "type": "string"
                },
                "heightPx": {
                    "type": "integer"
                },
                "size": {
                    "type": "integer"
                },
                "widthPx": {
                    "type": "integer"
                }
            }
        },
        "directory.createRequest": {
            "description": "Request structure to create a new directory",
            "type": "object",
            "properties": {
                "path": {
                    "type": "string"
                }
            }
        },
        "directory.readAllResponse": {
            "description": "List of directories",
            "type": "object",
            "properties": {
                "directories": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/directory.readAllResponseItem"
                    }
                }
            }
        },
        "directory.readAllResponseItem": {
            "description": "Directory details",
            "type": "object",
            "properties": {
                "dirId": {
                    "type": "integer"
                },
                "lastScanned": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                }
            }
        },
        "track.readAllResponse": {
            "description": "Response structure containing details of all tracks.",
            "type": "object",
            "properties": {
                "tracks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/track.readAllResponseItem"
                    }
                }
            }
        },
        "track.readAllResponseItem": {
            "description": "Single track item structure in the response of ReadAll endpoint.",
            "type": "object",
            "properties": {
                "audioCodec": {
                    "type": "string"
                },
                "coverId": {
                    "type": "integer"
                },
                "durationMs": {
                    "type": "integer"
                },
                "hashSha256": {
                    "type": "string"
                },
                "size": {
                    "type": "integer"
                },
                "trackId": {
                    "type": "integer"
                }
            }
        },
        "track.readResponse": {
            "description": "Response structure containing details of a single track.",
            "type": "object",
            "properties": {
                "audioCodec": {
                    "type": "string"
                },
                "bitrateKbps": {
                    "type": "integer"
                },
                "channels": {
                    "type": "integer"
                },
                "coverId": {
                    "type": "integer"
                },
                "durationMs": {
                    "type": "integer"
                },
                "hashSha256": {
                    "type": "string"
                },
                "sampleRateHz": {
                    "type": "integer"
                },
                "size": {
                    "type": "integer"
                },
                "trackId": {
                    "type": "integer"
                }
            }
        },
        "types.ErrorResponse": {
            "description": "Standard error response",
            "type": "object",
            "required": [
                "error"
            ],
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.2",
	Host:             "localhost:8022",
	BasePath:         "/api/music-files-service",
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
