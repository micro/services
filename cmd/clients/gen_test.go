package main

import (
	"encoding/json"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

type tspec struct {
	openapi  string
	tsresult string
	key      string
}

var cases = []tspec{
	{
		openapi: `{
	"components": {
	  "schemas": {
		"QueryRequest": {
		  "description": "Query posts. Acts as a listing when no id or slug provided.\n Gets a single post by id or slug if any of them provided.",
		  "properties": {
			"id": {
			  "type": "string"
			},
			"limit": {
			  "format": "int64",
			  "type": "number"
			},
			"offset": {
			  "format": "int64",
			  "type": "number"
			},
			"slug": {
			  "type": "string"
			},
			"tag": {
			  "type": "string"
			}
		  },
		  "title": "QueryRequest",
		  "type": "object"
		}
	  }
	}
  }`,
		key: "QueryRequest",
		tsresult: `export interface QueryRequest {
  id?: number;
  limit?: number;
  offset?: number;
  slug?: number;
  tag?: number;
}`,
	},
	{
		openapi: `{"components": { "schemas": {
	"QueryResponse": {
	  "properties": {
		"posts": {
		  "items": {
			"properties": {
			  "author": {
				"type": "string"
			  },
			  "content": {
				"type": "string"
			  },
			  "created": {
				"format": "int64",
				"type": "number"
			  },
			  "id": {
				"type": "string"
			  },
			  "image": {
				"type": "string"
			  },
			  "metadata": {
				"items": {
				  "properties": {
					"key": {
					  "type": "string"
					},
					"value": {
					  "type": "string"
					}
				  },
				  "type": "object"
				},
				"type": "array"
			  },
			  "slug": {
				"type": "string"
			  },
			  "tags": {
				"items": {
				  "type": "string"
				},
				"type": "array"
			  },
			  "title": {
				"type": "string"
			  },
			  "updated": {
				"format": "int64",
				"type": "number"
			  }
			},
			"type": "object"
		  },
		  "type": "array"
		}
	  },
	  "title": "QueryResponse",
	  "type": "object"
}}}}`,
		key: "QueryResponse",
		tsresult: `
export interface QueryResponse {
	posts?: {
	  author?: string;
	  content?: string;
	  created?: number;
	  id?: string;
	  image?: string;
	  metadata?: {
		key?: string;
		value?: string;
	  }[];
	  slug?: string;
	  tags?: string[];
	  title?: string;
	  updated?: number;
	}[];
}`,
	},
}

func TestTsGen(t *testing.T) {
	for _, c := range cases {
		spec := &openapi3.Swagger{}
		err := json.Unmarshal([]byte(c.openapi), &spec)
		if err != nil {
			t.Fatal(err)
		}
		//spew.Dump(spec.Components.Schemas)
		res := schemaToTs(c.key, spec.Components.Schemas[c.key])
		if res != c.tsresult {
			t.Logf("Expected %v, got: %v", c.tsresult, res)
		}
	}
}
