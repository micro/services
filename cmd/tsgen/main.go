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
			//gents := exec.Command("npx", "openapi-typescript", apiJSON, "--output", serviceName+".ts")
			//gents.Dir = serviceDir
			//fmt.Println(serviceDir)
			//outp, err := gents.CombinedOutput()
			//if err != nil {
			//	fmt.Println("Failed to make docs", string(outp))
			//	os.Exit(1)
			//}
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

			tsContent := ""
			typeNames := []string{}
			for k, v := range spec.Components.Schemas {
				tsContent += schemaToTs(k, v) + "\n\n"
				typeNames = append(typeNames, k)
			}
			os.MkdirAll(filepath.Join(tsPath, serviceName), 0777)
			f, err := os.OpenFile(filepath.Join(tsPath, serviceName, "index.ts"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				fmt.Println("Failed to open schema file", err)
				os.Exit(1)
			}
			_, err = f.Write([]byte(tsContent))
			if err != nil {
				fmt.Println("Failed to append to schema file", err)
				os.Exit(1)
			}

			f, err = os.OpenFile(filepath.Join(tsPath, "index.ts"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				fmt.Println("Failed to open index.ts", err)
				os.Exit(1)
			}

			_, err = f.Write([]byte(""))
			if err != nil {
				fmt.Println("Failed to append to index file", err)
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
	getVersions := exec.Command("npm", "show", "@micro/services", "--time", "--json")
	getVersions.Dir = tsPath

	outp, err := getVersions.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to get versions of NPM package", string(outp))
		os.Exit(1)
	}
	type npmVers struct {
		Versions []string `json:"versions"`
	}
	npmOutput := &npmVers{}
	var latest *semver.Version
	if len(outp) > 0 {
		err = json.Unmarshal(outp, npmOutput)
		if err != nil {
			fmt.Println("Failed to unmarshal versions", string(outp))
			os.Exit(1)
		}
	}

	for _, version := range npmOutput.Versions {
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
	if latest == nil {
		latest, _ = semver.NewVersion("0.0.0")
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

func schemaToTs(title string, spec *openapi3.SchemaRef) string {
	var recurse func(props map[string]*openapi3.SchemaRef, level int) string

	recurse = func(props map[string]*openapi3.SchemaRef, level int) string {
		ret := ""

		i := 0
		for k, v := range props {
			ret += strings.Repeat("  ", level)
			k = strcase.SnakeCase(k)
			//v.Value.
			switch v.Value.Type {
			case "object":
				// @todo identify what is a slice and what is not!
				// currently the openapi converter messes this up
				// see redoc html output
				ret += k + "?: {\n" + recurse(v.Value.Properties, level+1) + strings.Repeat("  ", level) + "};"

			case "array":
				if len(v.Value.Items.Value.Properties) == 0 {
					ret += k + "?: " + v.Value.Items.Value.Type + "[];"
				} else {
					// @todo identify what is a slice and what is not!
					// currently the openapi converter messes this up
					// see redoc html output
					ret += k + "?: {\n" + recurse(v.Value.Items.Value.Properties, level+1) + strings.Repeat("  ", level) + "}[];"
				}
			case "string":
				ret += k + "?: " + "string;"
			case "number":
				ret += k + "?: " + "number;"
			case "boolean":
				ret += k + "?: " + "boolean;"
			}

			if i < len(props) {
				ret += "\n"
			}
			i++

		}
		return ret
	}
	return "export interface " + title + " {\n" + recurse(spec.Value.Properties, 1) + "}"
}

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
