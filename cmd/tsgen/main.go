package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stoewer/go-strcase"
)

const (
	postContentPath = "docs/hugo-tania/site/content/post"
	docsURL         = "services.m3o.com"
)

func main() {
	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	workDir, _ := os.Getwd()

	tsPath := filepath.Join(workDir, "clients", "ts")

	for _, f := range files {
		if f.IsDir() && !strings.HasPrefix(f.Name(), ".") {

			serviceDir := filepath.Join(workDir, f.Name())
			serviceFiles, err := ioutil.ReadDir(serviceDir)
			if err != nil {
				fmt.Println("Failed to read service dir", err)
				os.Exit(1)
			}
			skip := false

			// detect openapi json file
			apiJSON := ""
			for _, serviceFile := range serviceFiles {
				if strings.Contains(serviceFile.Name(), "api") && strings.HasSuffix(serviceFile.Name(), ".json") {
					apiJSON = serviceFile.Name()
				}
				if serviceFile.Name() == "skip" {
					skip = true
				}
			}
			if skip {
				continue
			}
			fmt.Println(apiJSON)

			fmt.Println("Processing folder", serviceDir)

			// generate typescript files from openapi json
			//gents := exec.Command("npx", "openapi-typescript", apiJSON, "--output", "schema.ts")
			//gents.Dir = serviceDir
			//fmt.Println(serviceDir)
			//outp, err := gents.CombinedOutput()
			//if err != nil {
			//	fmt.Println("Failed to make docs", string(outp))
			//	os.Exit(1)
			//}

			// get latest version from github
			getVersions := exec.Command("npm", "show", "@micro/services", "time", "--json")
			getVersions.Dir = serviceDir
			fmt.Println(serviceDir)
			outp, err := getVersions.CombinedOutput()
			if err != nil {
				fmt.Println("Failed to make docs", string(outp))
				os.Exit(1)
			}
			versions := map[string]interface{}{}
			json.Unmarshal(outp, versions)
			var latest *semver.Version
			for version, _ := range versions {
				v, err := semver.NewVersion(version)
				if err != nil {
					fmt.Println("Failed to parse semver", err)
					os.Exit(1)
				}
				if latest == nil {
					latest = v
				}
				if v.GreaterThan(latest) {
					latest = v
				}
			}
			latest.IncPatch()

			// bump package to latest version
			repl := exec.Command("sed", "-i", "-e", "'s/0.0.1/"+latest.String()+"/g'", "package.json")
			repl.Dir = tsPath
			outp, err = repl.CombinedOutput()
			if err != nil {
				fmt.Println("Failed to make docs", string(outp))
				os.Exit(1)
			}
		}
	}
}

type specType struct {
	name          string
	tag           string
	includeReadme bool
	filePostFix   string
	titlePostFix  string
	template      string
}

var specTypes = []specType{
	{
		name:          "default markdown",
		tag:           "Readme",
		filePostFix:   ".md",
		template:      defTempl,
		includeReadme: true,
	},
	{
		name:          "microjs markdown",
		tag:           "Micro.js",
		filePostFix:   "-microjs.md",
		titlePostFix:  " Micro.js",
		template:      microJSTempl,
		includeReadme: false,
	},
}

var servicesToTags = map[string][]string{
	"users":      []string{"Backend"},
	"helloworld": []string{"Backend"},
	"emails":     []string{"Communications"},
	"sms":        []string{"Communications"},
	"posts":      []string{"Headless CMS"},
	"tags":       []string{"Headless CMS"},
	"feeds":      []string{"Headless CMS"},
	"datastore":  []string{"Backend"},
	"geocoding":  []string{"Logistics"},
	"places":     []string{"Logistics"},
	"routing":    []string{"Logistics"},
	"etas":       []string{"Logistics"},
	"notes":      []string{"Misc"},
	"messages":   []string{"Misc"},
}

func saveSpec(originalMarkDown []byte, contentDir, serviceName string, spec *openapi3.Swagger) error {
	for _, v := range specTypes {
		fmt.Println("Processing ", v.name)
		contentFile := filepath.Join(contentDir, serviceName+v.filePostFix)
		var app []byte
		if v.includeReadme {
			app = originalMarkDown
		}
		tags := []string{v.tag}
		serviceTags, ok := servicesToTags[serviceName]
		if ok {
			tags = append(tags, serviceTags...)
		}
		tagsString := "\n- " + strings.Join(tags, "\n- ")

		err := ioutil.WriteFile(contentFile, append([]byte("---\ntitle: "+serviceName+v.titlePostFix+"\nservicename: "+serviceName+"\nlabels: "+tagsString+"\n---\n"), app...), 0777)
		if err != nil {
			fmt.Printf("Failed to write post content to %v:\n%v\n", err)
			os.Exit(1)
		}
		fi, err := os.OpenFile(contentFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return err
		}
		tmpl, err := template.New("test").Funcs(template.FuncMap{
			"params": func(p openapi3.Parameters) string {
				ls := ""
				for _, v := range p {
					//if v.Value.In == "body" {
					bs, _ := v.MarshalJSON()
					ls += string(bs) + ", "
					//}
				}
				return ls
			},
			// @todo should take SpecRef here not RequestBodyRef
			"schemaJSON": func(prepend int, ref string) string {
				for k, v := range spec.Components.Schemas {
					// ie. #/components/requestBodies/PostsSaveRequest contains
					// SaveRequest, can't see any other way to correlate
					if strings.HasSuffix(ref, k) {
						bs, _ := json.MarshalIndent(schemaToMap(v, spec.Components.Schemas), "", strings.Repeat(" ", prepend)+"  ")
						// last line wont get prepended so we fix that here
						parts := strings.Split(string(bs), "\n")
						// skip if it's only 1 line, ie it's '{}'
						if len(parts) <= 1 {
							return string(bs)
						}
						parts[len(parts)-1] = parts[len(parts)-1]
						return strings.Join(parts, "\n")
					}
				}

				return "Schema related to " + ref + " not found"

			},
			"schemaDescription": func(ref string) string {
				for k, v := range spec.Components.Schemas {
					// ie. #/components/requestBodies/PostsSaveRequest contains
					// SaveRequest, can't see any other way to correlate
					if strings.HasSuffix(ref, k) {
						return v.Value.Description
					}
				}

				return "Schema related to " + ref + " not found"
			},
			// turn chat/Chat/History
			// to Chat History
			"titleize": func(s string) string {
				parts := strings.Split(s, "/")
				if len(parts) > 2 {
					return strings.Join(parts[2:], " ")
				}
				return strings.Join(parts, " ")
			},
			"firstResponseRef": func(rs openapi3.Responses) string {
				return rs.Get(200).Ref
			},
		}).Parse(v.template)
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(fi, spec)
		if err != nil {
			return err
		}
	}
	return nil
}

func schemaToMap(spec *openapi3.SchemaRef, schemas map[string]*openapi3.SchemaRef) map[string]interface{} {
	var recurse func(props map[string]*openapi3.SchemaRef) map[string]interface{}
	recurse = func(props map[string]*openapi3.SchemaRef) map[string]interface{} {
		ret := map[string]interface{}{}
		for k, v := range props {
			k = strcase.SnakeCase(k)
			//v.Value.
			if v.Value.Type == "object" {
				// @todo identify what is a slice and what is not!
				// currently the openapi converter messes this up
				// see redoc html output
				ret[k] = recurse(v.Value.Properties)
				continue
			}
			if v.Value.Type == "array" {
				// @todo identify what is a slice and what is not!
				// currently the openapi converter messes this up
				// see redoc html output
				ret[k] = []interface{}{recurse(v.Value.Properties)}
				continue
			}
			switch v.Value.Type {
			case "string":
				if len(v.Value.Description) > 0 {
					ret[k] = strings.Replace(v.Value.Description, "\n", ".", -1)
				} else {
					ret[k] = v.Value.Type
				}
			case "number":
				ret[k] = 1
			case "boolean":
				ret[k] = true
			}

		}
		return ret
	}
	return recurse(spec.Value.Properties)
}

const defTempl = `
## cURL

{{ range $key, $value := .Paths }}
### {{ $key | titleize }}
<!-- We use the request body description here as endpoint descriptions are not
being lifted correctly from the proto by the openapi spec generator -->
{{ $value.Post.RequestBody.Ref | schemaDescription }}
` + "```" + `shell
> curl 'https://api.m3o.com{{ $key }}' \
  -H 'micro-namespace: $yourNamespace' \
  -H 'authorization: Bearer $yourToken' \
  -d {{ $value.Post.RequestBody.Ref | schemaJSON 0 }};
# Response
{{ $value.Post.Responses | firstResponseRef | schemaJSON 0 }}
` + "```" + `

{{ end }}
`

const microJSTempl = `
## Micro.js

{{ range $key, $value := .Paths }}
### {{ $key | titleize }}
<!-- We use the request body description here as endpoint descriptions are not
being lifted correctly from the proto by the openapi spec generator -->
{{ $value.Post.RequestBody.Ref | schemaDescription }}
` + "```" + `html
<script src="https://web.m3o.com/assets/micro.js"></script>
<script type="text/javascript">
  document.addEventListener("DOMContentLoaded", function (event) {
    // Login is only required for endpoints doing authorization
    Micro.requireLogin(function () {
      Micro.post(
        "{{ $key }}",
        "micro",
        {{ $value.Post.RequestBody.Ref | schemaJSON 8 }},
        function (data) {
          console.log("Success.");
        }
      );
    });
  });
</script>
` + "```" + `

{{ end }}
`

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
// from https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
