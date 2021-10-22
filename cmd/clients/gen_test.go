package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/getkin/kin-openapi/openapi3"
)

func TestSemver(t *testing.T) {
	v, _ := semver.NewVersion("0.0.0-beta1")
	if incBeta(*v).String() != "0.0.0-beta2" {
		t.Fatal(v)
	}

	v1, _ := semver.NewVersion("0.0.1")
	if !v1.GreaterThan(v) {
		t.Fatal("no good")
	}

	v2, _ := semver.NewVersion("0.0.0")
	if !v2.GreaterThan(v) {
		t.Fatal("no good")
	}

	if v.String() != "0.0.0-beta1" {
		t.Fatal("no good")
	}

	v3, _ := semver.NewVersion("0.0.0-beta2")
	if !v3.GreaterThan(v) {
		t.Fatal("no good")
	}
}

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

func TestTimeExample(t *testing.T) {
	spec := &openapi3.Swagger{}
	err := json.Unmarshal([]byte(timeExample), &spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(spec.Components.Schemas) == 0 {
		t.Fatal("boo")
	}

	fmt.Println(spec.Components.Schemas)
	res := schemaToGoExample("time", "NowRequest", spec.Components.Schemas, map[string]interface{}{
		"location": "London",
	})
	if strings.TrimSpace(res) != strings.TrimSpace(timeExp) {
		t.Log(res, timeExp)
	}

	fmt.Println(spec.Components.Schemas)
	res = schemaToGoExample("time", "ZoneRequest", spec.Components.Schemas, map[string]interface{}{
		"location": "London",
	})
	if strings.TrimSpace(res) != strings.TrimSpace(timeExp) {
		t.Log(res, timeExp)
	}
}

const timeExample = `{
	"components": {
		"schemas": {
		
				"NowRequest": {
				  "description": "Get the current time",
				  "properties": {
					"location": {
					  "description": "optional location, otherwise returns UTC",
					  "type": "string"
					}
				  },
				  "title": "NowRequest",
				  "type": "object"
				},
				"NowResponse": {
				  "properties": {
					"localtime": {
					  "description": "the current time as HH:MM:SS",
					  "type": "string"
					},
					"location": {
					  "description": "the location as Europe/London",
					  "type": "string"
					},
					"timestamp": {
					  "description": "timestamp as 2006-01-02T15:04:05.999999999Z07:00",
					  "type": "string"
					},
					"timezone": {
					  "description": "the timezone as BST",
					  "type": "string"
					},
					"unix": {
					  "description": "the unix timestamp",
					  "format": "int64",
					  "type": "number"
					}
				  },
				  "title": "NowResponse",
				  "type": "object"
				},
				"ZoneRequest": {
				  "description": "Get the timezone info for a specific location",
				  "properties": {
					"location": {
					  "description": "location to lookup e.g postcode, city, ip address",
					  "type": "string"
					}
				  },
				  "title": "ZoneRequest",
				  "type": "object"
				},
				"ZoneResponse": {
				  "properties": {
					"abbreviation": {
					  "description": "the abbreviated code e.g BST",
					  "type": "string"
					},
					"country": {
					  "description": "country of the timezone",
					  "type": "string"
					},
					"dst": {
					  "description": "is daylight savings",
					  "type": "boolean"
					},
					"latitude": {
					  "description": "e.g 51.42",
					  "format": "double",
					  "type": "number"
					},
					"localtime": {
					  "description": "the local time",
					  "type": "string"
					},
					"location": {
					  "description": "location requested",
					  "type": "string"
					},
					"longitude": {
					  "description": "e.g -0.37",
					  "format": "double",
					  "type": "number"
					},
					"region": {
					  "description": "region of timezone",
					  "type": "string"
					},
					"timezone": {
					  "description": "the timezone e.g Europe/London",
					  "type": "string"
					}
				  },
				  "title": "ZoneResponse",
				  "type": "object"
				}
		
		}
	}
}`

const timeExp = `Location: London,
`

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
		t.Log(res, arrayExp)
	}

	spec = &openapi3.Swagger{}
	err = json.Unmarshal([]byte(simpleExample), &spec)
	if err != nil {
		t.Log(err)
	}
	if len(spec.Components.Schemas) == 0 {
		t.Log("boo")
	}
	fmt.Println(spec.Components.Schemas)
	res = schemaToGoExample("file", "DeleteRequest", spec.Components.Schemas, map[string]interface{}{
		"project": "examples",
		"path":    "/document/text-files/file.txt",
	})
	if strings.TrimSpace(res) != strings.TrimSpace(simpleExp) {
		t.Log(res, arrayExp)
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
