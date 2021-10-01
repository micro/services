package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/crufter/nested"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stoewer/go-strcase"
)

func schemaToGoExample(serviceName, typeName string, schemas map[string]*openapi3.SchemaRef, values map[string]interface{}) string {
	var recurse func(props map[string]*openapi3.SchemaRef, path []string) string

	var spec *openapi3.SchemaRef = schemas[typeName]
	if spec == nil {
		existing := ""
		for k, _ := range schemas {
			existing += k + " "
		}
		panic("can't find schema " + typeName + " but found " + existing)
	}
	detectType := func(currentType string, properties map[string]*openapi3.SchemaRef) (string, bool) {
		index := map[string]bool{}
		for key, prop := range properties {
			index[key+prop.Value.Title] = true
		}
		for k, schema := range schemas {
			// we don't want to return the type matching itself
			if strings.ToLower(k) == currentType {
				continue
			}
			if strings.HasSuffix(k, "Request") || strings.HasSuffix(k, "Response") {
				continue
			}
			if len(schema.Value.Properties) != len(properties) {
				continue
			}
			found := false
			for key, prop := range schema.Value.Properties {

				_, ok := index[key+prop.Value.Title]
				found = ok
				if !ok {
					break
				}
			}
			if found {
				return schema.Value.Title, true
			}
		}
		return "", false
	}
	var fieldSeparator, objectOpen, objectClose, arrayPrefix, arrayPostfix, fieldDelimiter, stringType, boolType string
	var int32Type, int64Type, floatType, doubleType, mapType, anyType, typeInstancePrefix string
	var fieldUpperCase bool
	language := "go"
	switch language {
	case "go":
		fieldUpperCase = true
		fieldSeparator = ": "
		arrayPrefix = "[]"
		arrayPostfix = ""
		objectOpen = "{\n"
		objectClose = "}"
		fieldDelimiter = ","
		stringType = "string"
		boolType = "bool"
		int32Type = "int32"
		int64Type = "int64"
		floatType = "float32"
		doubleType = "float64"
		mapType = "map[string]%v"
		anyType = "interface{}"
		typeInstancePrefix = "&"
	}

	valueToType := func(v *openapi3.SchemaRef) string {
		switch v.Value.Type {
		case "string":
			return stringType
		case "boolean":
			return boolType
		case "number":
			switch v.Value.Format {
			case "int32":
				return int32Type
			case "int64":
				return int64Type
			case "float":
				return floatType
			case "double":
				return doubleType
			}
		default:
			return "unrecognized: " + v.Value.Type
		}
		return ""
	}

	printMap := func(m map[string]interface{}, level int) string {
		ret := ""
		for k, v := range m {
			marsh, _ := json.Marshal(v)
			ret += strings.Repeat("\t", level) + fmt.Sprintf("\"%v\": %v,\n", k, string(marsh))
		}
		return ret
	}

	recurse = func(props map[string]*openapi3.SchemaRef, path []string) string {
		ret := ""

		i := 0
		var keys []string
		for k := range props {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, v := range path {
			path[i] = strcase.LowerCamelCase(v)
		}
		for _, k := range keys {
			v := props[k]
			ret += strings.Repeat("\t", len(path))
			if fieldUpperCase {
				k = strcase.UpperCamelCase(k)
			}

			var val interface{}
			p := strings.Replace(strings.Join(append(path, strcase.LowerCamelCase(k)), "."), ".[", "[", -1)
			val, ok := nested.Get(values, p)
			if !ok {
				continue
			}
			// hack
			if str, ok := val.(string); ok {
				if str == "<nil>" {
					continue
				}
			}
			switch v.Value.Type {
			case "object":
				typ, found := detectType(k, v.Value.Properties)
				if found {
					ret += k + fieldSeparator + typeInstancePrefix + serviceName + "." + strings.Title(typ) + objectOpen + recurse(v.Value.Properties, append(path, k)) + objectClose + fieldDelimiter
				} else {
					// type is a dynamic map
					// if additional properties is present, then it's a map string string or other typed map
					if v.Value.AdditionalProperties != nil {
						ret += k + fieldSeparator + fmt.Sprintf(mapType, valueToType(v.Value.AdditionalProperties)) + objectOpen + printMap(val.(map[string]interface{}), len(path)+1) + objectClose + fieldDelimiter
					} else {
						// if additional properties is not present, it's an any type,
						// like the proto struct type
						ret += k + fieldSeparator + fmt.Sprintf(mapType, anyType) + objectOpen + printMap(val.(map[string]interface{}), len(path)+1) + objectClose + fieldDelimiter
					}
				}
			case "array":
				typ, found := detectType(k, v.Value.Items.Value.Properties)
				if found {
					ret += k + fieldSeparator + arrayPrefix + serviceName + "." + strings.Title(typ) + objectOpen + serviceName + "." + strings.Title(typ) + objectOpen + recurse(v.Value.Items.Value.Properties, append(append(path, k), "[0]")) + objectClose + objectClose + arrayPostfix + fieldDelimiter
				} else {
					arrint := val.([]interface{})
					switch v.Value.Items.Value.Type {
					case "string":
						arrstr := make([]string, len(arrint))
						for i, v := range arrint {
							arrstr[i] = fmt.Sprintf("%v", v)
						}

						ret += k + fieldSeparator + fmt.Sprintf("%#v", arrstr) + fieldDelimiter
					case "number", "boolean":
						ret += k + fieldSeparator + arrayPrefix + fmt.Sprintf("%v", val) + arrayPostfix + fieldDelimiter
					case "object":
						ret += k + fieldSeparator + arrayPrefix + fmt.Sprintf(mapType, valueToType(v.Value.AdditionalProperties)) + objectOpen + fmt.Sprintf(mapType, valueToType(v.Value.AdditionalProperties)) + objectOpen + recurse(v.Value.Items.Value.Properties, append(append(path, k), "[0]")) + strings.Repeat("\t", len(path)) + objectClose + objectClose + arrayPostfix + fieldDelimiter
					}
				}
			case "string":
				if strings.Contains(val.(string), "\n") {
					ret += k + fieldSeparator + fmt.Sprintf("`%v`", val) + fieldDelimiter
				} else {
					ret += k + fieldSeparator + fmt.Sprintf("\"%v\"", val) + fieldDelimiter
				}
			case "number", "boolean":
				ret += k + fieldSeparator + fmt.Sprintf("%v", val) + fieldDelimiter
			}

			if i < len(props) {
				ret += "\n"
			}
			i++

		}
		return ret
	}
	return recurse(spec.Value.Properties, []string{})
}
