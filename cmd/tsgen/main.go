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
			serviceName := f.Name()
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
					apiJSON = filepath.Join(serviceDir, serviceFile.Name())
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
			gents := exec.Command("npx", "openapi-typescript", apiJSON, "--output", serviceName+".ts")
			gents.Dir = serviceDir
			fmt.Println(serviceDir)
			outp, err := gents.CombinedOutput()
			if err != nil {
				fmt.Println("Failed to make docs", string(outp))
				os.Exit(1)
			}

			// copy generated file to folder
			copyFileContents(filepath.Join(serviceDir, serviceName+".ts"), filepath.Join(tsPath, serviceName+"_schema.ts"))

			f, err := os.OpenFile(filepath.Join(serviceDir, "index.ts"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				fmt.Println("Failed to open npmrc", err)
				os.Exit(1)
			}
			_, err = f.Write([]byte("export { default as " + strings.Title(serviceName) + " } from './" + serviceName + "';\n"))
			if err != nil {
				fmt.Println("Failed to append to index file", err)
				os.Exit(1)
			}

			js, err := ioutil.ReadFile(apiJSON)

			if err != nil {
				fmt.Println("Failed to read json spec", err)
				os.Exit(1)
			}
			spec := &openapi3.Swagger{}
			err = json.Unmarshal(js, &spec)
			if err != nil {
				fmt.Println("Failed to unmarshal", err)
				os.Exit(1)
			}
			err = saveFile(tsPath, serviceName, spec)
			if err != nil {
				fmt.Println("Failed to generate app", err)
				os.Exit(1)
			}
		}
	}
	// login to NPM
	f, err := os.OpenFile(filepath.Join(tsPath, ".npmrc"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Failed to open npmrc", err)
		os.Exit(1)
	}

	defer f.Close()
	if len(os.Getenv("NPM_TOKEN")) == 0 {
		fmt.Println("No NPM_TOKEN env found")
		os.Exit(1)
	}
	if _, err = f.WriteString("\n//npm.pkg.github.com/:_authToken=" + os.Getenv("NPM_TOKEN")); err != nil {
		fmt.Println("Failed to open npmrc", err)
		os.Exit(1)
	}

	// get latest version from github
	getVersions := exec.Command("npm", "show", "@micro/services", "time", "--json")
	getVersions.Dir = tsPath

	outp, err := getVersions.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to get versions of NPM package", string(outp))
		os.Exit(1)
	}
	versions := map[string]interface{}{}
	err = json.Unmarshal(outp, &versions)
	if err != nil {
		fmt.Println("Failed to unmarshal versions", string(outp))
		os.Exit(1)
	}

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
	newV := latest.IncPatch()

	// bump package to latest version
	fmt.Println("Bumping to ", newV.String())
	repl := exec.Command("sed", "-i", "-e", "s/1.0.1/"+newV.String()+"/g", "package.json")
	repl.Dir = tsPath
	outp, err = repl.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to make docs", string(outp))
		os.Exit(1)
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
}

func saveFile(tsDir string, serviceName string, spec *openapi3.Swagger) error {
	for _, v := range specTypes {
		fmt.Println("Processing ", v.name)
		contentFile := filepath.Join(tsDir, serviceName+".ts")
		fi, err := os.OpenFile(contentFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			return err
		}
		tmpl, err := template.New("test").Funcs(template.FuncMap{
			"toLower": func(s string) string {
				return strings.ToLower(s)
			},
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
import { components } from './{{ .Info.Title | toLower }}_schema';

export interface types extends components {};
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
