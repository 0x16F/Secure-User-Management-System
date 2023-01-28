// Code generated by swaggo/swag. DO NOT EDIT
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
        "/auth/login": {
            "post": {
                "description": "login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "login",
                "operationId": "login",
                "parameters": [
                    {
                        "description": "credentials",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.AccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    }
                }
            }
        },
        "/auth/refresh": {
            "post": {
                "description": "refresh jwt access token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "refresh",
                "operationId": "refresh",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.AccessResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get users",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "get users",
                "operationId": "get users",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 50,
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "asc",
                            "desc"
                        ],
                        "type": "string",
                        "default": "asc",
                        "description": "order",
                        "name": "order",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "",
                        "description": "name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "",
                        "description": "surname",
                        "name": "surname",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "",
                        "description": "login",
                        "name": "login",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "",
                        "description": "permissions",
                        "name": "permissions",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "offset",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.FindUsersResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
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
                "description": "create user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "create user",
                "operationId": "create user",
                "parameters": [
                    {
                        "description": "create user info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.CreateUserDTO"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/users.CreateUserResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "get user",
                "operationId": "get user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.FindUserDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
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
                "description": "delete user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "delete user",
                "operationId": "delete user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "update user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "update user",
                "operationId": "update user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "update info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.UpdateUserDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.AppError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.AccessResponse": {
            "type": "object",
            "properties": {
                "access": {
                    "type": "string"
                }
            }
        },
        "auth.LoginRequest": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "response.AppError": {
            "type": "object",
            "properties": {
                "developer_message": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "params": {
                    "$ref": "#/definitions/response.Map"
                }
            }
        },
        "response.Map": {
            "type": "object",
            "additionalProperties": true
        },
        "response.SuccessResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "params": {
                    "$ref": "#/definitions/response.Map"
                }
            }
        },
        "user.CreateUserDTO": {
            "type": "object",
            "properties": {
                "birthday": {
                    "type": "integer"
                },
                "login": {
                    "type": "string",
                    "default": "Ivanov"
                },
                "name": {
                    "type": "string",
                    "default": "Иван"
                },
                "password": {
                    "type": "string",
                    "default": "password"
                },
                "permissions": {
                    "type": "string",
                    "enum": [
                        "admin",
                        "read-only",
                        "banned"
                    ]
                },
                "surname": {
                    "type": "string",
                    "default": "Иванов"
                }
            }
        },
        "user.FindUserDTO": {
            "type": "object",
            "properties": {
                "birthday": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "login": {
                    "type": "string",
                    "default": "Ivanov"
                },
                "name": {
                    "type": "string",
                    "default": "Иван"
                },
                "permissions": {
                    "type": "string",
                    "enum": [
                        "admin",
                        "read-only",
                        "banned"
                    ]
                },
                "surname": {
                    "type": "string",
                    "default": "Иванов"
                }
            }
        },
        "user.UpdateUserDTO": {
            "type": "object",
            "properties": {
                "birthday": {
                    "type": "integer"
                },
                "login": {
                    "type": "string",
                    "default": "Ivanov"
                },
                "name": {
                    "type": "string",
                    "default": "Иван"
                },
                "password": {
                    "type": "string",
                    "default": "password"
                },
                "permissions": {
                    "type": "string",
                    "enum": [
                        "admin",
                        "read-only",
                        "banned"
                    ]
                },
                "surname": {
                    "type": "string",
                    "default": "Иванов"
                }
            }
        },
        "users.CreateUserResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "users.FindUsersResponse": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/user.FindUserDTO"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Test Project",
	Description:      "Just solving a test task",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
