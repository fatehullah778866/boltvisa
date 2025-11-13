package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetOpenAPISpec returns OpenAPI 3.0 specification
func (h *Handlers) GetOpenAPISpec(c *gin.Context) {
	spec := `{
  "openapi": "3.0.0",
  "info": {
    "title": "Visa Help Center API",
    "version": "1.0.0",
    "description": "Comprehensive visa application management system API"
  },
  "servers": [
    {
      "url": "http://localhost:8080",
      "description": "Development server"
    },
    {
      "url": "https://api.boltvisa.com",
      "description": "Production server"
    }
  ],
  "components": {
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    },
    "schemas": {
      "User": {
        "type": "object",
        "properties": {
          "id": { "type": "integer" },
          "email": { "type": "string", "format": "email" },
          "first_name": { "type": "string" },
          "last_name": { "type": "string" },
          "role": { "type": "string", "enum": ["admin", "consultant", "applicant"] },
          "active": { "type": "boolean" },
          "created_at": { "type": "string", "format": "date-time" },
          "updated_at": { "type": "string", "format": "date-time" }
        }
      },
      "VisaApplication": {
        "type": "object",
        "properties": {
          "id": { "type": "integer" },
          "user_id": { "type": "integer" },
          "category_id": { "type": "integer" },
          "status": { "type": "string", "enum": ["draft", "submitted", "in_review", "approved", "rejected", "cancelled"] },
          "passport_number": { "type": "string" },
          "nationality": { "type": "string" },
          "created_at": { "type": "string", "format": "date-time" }
        }
      },
      "Error": {
        "type": "object",
        "properties": {
          "error": { "type": "string" }
        }
      }
    }
  },
  "paths": {
    "/api/v1/auth/register": {
      "post": {
        "summary": "Register a new user",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["email", "password", "first_name", "last_name"],
                "properties": {
                  "email": { "type": "string", "format": "email" },
                  "password": { "type": "string", "minLength": 8 },
                  "first_name": { "type": "string" },
                  "last_name": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "User created successfully",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "token": { "type": "string" },
                    "user": { "$ref": "#/components/schemas/User" }
                  }
                }
              }
            }
          },
          "400": { "$ref": "#/components/responses/BadRequest" },
          "409": { "$ref": "#/components/responses/Conflict" }
        }
      }
    },
    "/api/v1/auth/login": {
      "post": {
        "summary": "Login user",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["email", "password"],
                "properties": {
                  "email": { "type": "string", "format": "email" },
                  "password": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Login successful",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "token": { "type": "string" },
                    "user": { "$ref": "#/components/schemas/User" }
                  }
                }
              }
            }
          },
          "401": { "$ref": "#/components/responses/Unauthorized" }
        }
      }
    },
    "/api/v1/users/me": {
      "get": {
        "summary": "Get current user",
        "security": [{ "bearerAuth": [] }],
        "responses": {
          "200": {
            "description": "User data",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/User" }
              }
            }
          },
          "401": { "$ref": "#/components/responses/Unauthorized" }
        }
      }
    },
    "/api/v1/applications": {
      "get": {
        "summary": "List applications",
        "security": [{ "bearerAuth": [] }],
        "parameters": [
          {
            "name": "page",
            "in": "query",
            "schema": { "type": "integer", "default": 1 }
          },
          {
            "name": "page_size",
            "in": "query",
            "schema": { "type": "integer", "default": 20 }
          },
          {
            "name": "status",
            "in": "query",
            "schema": { "type": "string" }
          }
        ],
        "responses": {
          "200": {
            "description": "List of applications",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "array",
                      "items": { "$ref": "#/components/schemas/VisaApplication" }
                    },
                    "page": { "type": "integer" },
                    "page_size": { "type": "integer" },
                    "total": { "type": "integer" },
                    "total_pages": { "type": "integer" }
                  }
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create new application",
        "security": [{ "bearerAuth": [] }],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["category_id"],
                "properties": {
                  "category_id": { "type": "integer" },
                  "passport_number": { "type": "string" },
                  "nationality": { "type": "string" }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Application created",
            "content": {
              "application/json": {
                "schema": { "$ref": "#/components/schemas/VisaApplication" }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "responses": {
      "BadRequest": {
        "description": "Bad request",
        "content": {
          "application/json": {
            "schema": { "$ref": "#/components/schemas/Error" }
          }
        }
      },
      "Unauthorized": {
        "description": "Unauthorized",
        "content": {
          "application/json": {
            "schema": { "$ref": "#/components/schemas/Error" }
          }
        }
      },
      "Conflict": {
        "description": "Conflict",
        "content": {
          "application/json": {
            "schema": { "$ref": "#/components/schemas/Error" }
          }
        }
      }
    }
  }
}`

	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, spec)
}
