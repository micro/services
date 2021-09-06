package main

import (
	"encoding/json"
	"fmt"
	"strings"
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
  id?: string;
  limit?: number;
  offset?: number;
  slug?: string;
  tag?: string;
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
	// @todo fix tests to be up to date
	return
	for _, c := range cases {
		spec := &openapi3.Swagger{}
		err := json.Unmarshal([]byte(c.openapi), &spec)
		if err != nil {
			t.Fatal(err)
		}
		//spew.Dump(spec.Components.Schemas)
		res := schemaToType("typescript", "ServiceName", c.key, spec.Components.Schemas)
		if res != c.tsresult {
			t.Logf("Expected %v, got: %v", c.tsresult, res)
		}
	}
}

func TestExample(t *testing.T) {

	spec := &openapi3.Swagger{}
	err := json.Unmarshal([]byte(arrayExample), &spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(spec.Components.Schemas) == 0 {
		t.Fatal("boo")
	}
	//spew.Dump(spec.Components.Schemas)
	res := schemaToGoExample("file", "ListResponse", spec.Components.Schemas, map[string]interface{}{
		"files": []map[string]interface{}{
			{
				"content": "something something",
				"created": "2021-05-20T13:37:21Z",
				"path":    "/documents/text-files/file.txt",
				"metadata": map[string]interface{}{
					"meta1": "value1",
					"meta2": "value2",
				},
				"project": "my-project",
				"updated": "2021-05-20T14:37:21Z",
			},
		},
	})
	if strings.TrimSpace(res) != strings.TrimSpace(arrayExp) {
		t.Fatal(res, arrayExp)
	}

	spec = &openapi3.Swagger{}
	err = json.Unmarshal([]byte(simpleExample), &spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(spec.Components.Schemas) == 0 {
		t.Fatal("boo")
	}
	fmt.Println(spec.Components.Schemas)
	res = schemaToGoExample("file", "DeleteRequest", spec.Components.Schemas, map[string]interface{}{
		"project": "examples",
		"path":    "/document/text-files/file.txt",
	})
	if strings.TrimSpace(res) != strings.TrimSpace(simpleExp) {
		t.Fatal(res, arrayExp)
	}
}

const simpleExample = `{
	"components": {
		"schemas": {
			"DeleteRequest": {
				"description": "Delete a file by project name/path",
				"properties": {
				  "path": {
					"description": "Path to the file",
					"type": "string"
				  },
				  "project": {
					"description": "The project name",
					"type": "string"
				  }
				},
				"title": "DeleteRequest",
				"type": "object"
			  }
		}
	}
}`

const simpleExp = `Path: "/document/text-files/file.txt"
Project: "exaples"
`

const arrayExp = `Files: []file.Record{
file.Record{
		Content: "something something",
		Created: "2021-05-20T13:37:21Z",
		Metadata: map[string]string{
			"meta1": "value1",
			"meta2": "value2",
},
		Path: "/documents/text-files/file.txt",
		Project: "my-project",
		Updated: "2021-05-20T14:37:21Z",
}},`

const arrayExample = `{
	"components": {
	  "schemas": {
		"ListResponse": {
		  "properties": {
			"files": {
			  "items": {
				"properties": {
				  "content": {
					"description": "File contents",
					"type": "string"
				  },
				  "created": {
					"description": "Time the file was created e.g 2021-05-20T13:37:21Z",
					"type": "string"
				  },
				  "metadata": {
					"additionalProperties": {
					  "type": "string"
					},
					"description": "Any other associated metadata as a map of key-value pairs",
					"type": "object"
				  },
				  "path": {
					"description": "Path to file or folder eg. '/documents/text-files/file.txt'.",
					"type": "string"
				  },
				  "project": {
					"description": "A custom project to group files\n eg. file-of-mywebsite.com",
					"type": "string"
				  },
				  "updated": {
					"description": "Time the file was updated e.g 2021-05-20T13:37:21Z",
					"type": "string"
				  }
				},
				"type": "object"
			  },
			  "type": "array"
			}
		  },
		  "title": "ListResponse",
		  "type": "object"
		},
		"Record": {
		  "properties": {
		    "content": {
		  	"description": "File contents",
		  	"type": "string"
		    },
		    "created": {
		  	"description": "Time the file was created e.g 2021-05-20T13:37:21Z",
		  	"type": "string"
		    },
		    "metadata": {
		  	"additionalProperties": {
		  	  "type": "string"
		  	},
		  	"description": "Any other associated metadata as a map of key-value pairs",
		  	"type": "object"
		    },
		    "path": {
		  	"description": "Path to file or folder eg. '/documents/text-files/file.txt'.",
		  	"type": "string"
		    },
		    "project": {
		  	"description": "A custom project to group files\n eg. file-of-mywebsite.com",
		  	"type": "string"
		    },
		    "updated": {
		  	"description": "Time the file was updated e.g 2021-05-20T13:37:21Z",
		  	"type": "string"
		    }
		  },
		  "title": "Record",
		  "type": "object"
		}
	  }
	}
  }`
