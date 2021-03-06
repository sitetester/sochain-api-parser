// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
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
        "/block/{network}/{blockNumberOrHash}": {
            "get": {
                "description": "Show block by network \u0026 number/hash",
                "tags": [
                    "blocks"
                ],
                "summary": "Show block",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Network",
                        "name": "network",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "block number or hash",
                        "name": "blockNumberOrHash",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.DesiredBlockResponseData"
                        }
                    }
                }
            }
        },
        "/tx/{network}/{hash}": {
            "get": {
                "description": "Show transaction by network \u0026 hash",
                "tags": [
                    "transactions"
                ],
                "summary": "Show transaction",
                "parameters": [
                    {
                        "type": "string",
                        "description": "network",
                        "name": "network",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "transaction hash",
                        "name": "hash",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.DesiredTxResponseData"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "service.DesiredBlockResponseData": {
            "type": "object",
            "properties": {
                "block_no": {
                    "type": "integer"
                },
                "network": {
                    "type": "string"
                },
                "next_blockhash": {
                    "type": "string"
                },
                "previous_blockhash": {
                    "type": "string"
                },
                "size": {
                    "type": "integer"
                },
                "time": {
                    "description": "actual type is ` + "`" + `int` + "`" + `, but we need to display this time in string format in API response",
                    "type": "string"
                },
                "txs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/service.DesiredTxResponseData"
                    }
                }
            }
        },
        "service.DesiredTxResponseData": {
            "type": "object",
            "properties": {
                "fee": {
                    "type": "string"
                },
                "sent_value": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                },
                "txid": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8081",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "Sochain API Explorer",
	Description: "This is an example server using Sochain API at backend",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
