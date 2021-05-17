package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func publishAPI(service, readme, openapiJSON, examplesJSON string, pricing map[string]int64) error {
	client := &http.Client{}

	apiSpec := map[string]interface{}{
		"name":          service,
		"description":   readme,
		"open_api_json": openapiJSON,
		"pricing":       pricing,
		"examples_json": examplesJSON,
	}

	//Encode the data
	postBody, _ := json.Marshal(map[string]interface{}{
		"api": apiSpec,
	})
	rbody := bytes.NewBuffer(postBody)

	//Leverage Go's HTTP Post function to make request
	req, err := http.NewRequest("POST", "https://api.m3o.com/publicapi/Publish", rbody)

	// Add auth headers here if needed
	req.Header.Add("Authorization", `Bearer `+os.Getenv("MICRO_ADMIN_TOKEN"))
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)

	return nil
}

func main() {
	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	workDir, _ := os.Getwd()

	for _, f := range files {
		if f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
			serviceDir := filepath.Join(workDir, f.Name())
			serviceFiles, err := ioutil.ReadDir(serviceDir)
			if err != nil {
				fmt.Println("Failed to read service dir", err)
				os.Exit(1)
			}
			skip := false
			for _, serviceFile := range serviceFiles {
				if serviceFile.Name() == "skip" {
					skip = true
				}
			}
			if skip {
				continue
			}

			fmt.Println("Processing folder", serviceDir)
			makeProto := exec.Command("make", "docs")
			makeProto.Dir = serviceDir
			fmt.Println(serviceDir)
			outp, err := makeProto.CombinedOutput()
			if err != nil {
				fmt.Println("Failed to make docs", string(outp))
				os.Exit(1)
			}
			serviceName := f.Name()
			dat, err := ioutil.ReadFile(filepath.Join(serviceDir, "README.md"))
			if err != nil {
				fmt.Println("Failed to read readme", string(outp))
				os.Exit(1)
			}

			apiJSON := filepath.Join(serviceDir, "api-"+serviceName+".json")
			js, err := ioutil.ReadFile(apiJSON)
			if err != nil {
				apiJSON := filepath.Join(serviceDir, "api-protobuf.json")
				js, err = ioutil.ReadFile(apiJSON)
				if err != nil {
					fmt.Println("Failed to read json spec", err)
					os.Exit(1)
				}
			}
			spec := &openapi3.Swagger{}
			err = json.Unmarshal(js, &spec)
			if err != nil {
				fmt.Println("Failed to unmarshal", err)
				os.Exit(1)
			}

			// not every service has examples
			examples, _ := ioutil.ReadFile(filepath.Join(serviceDir, "examples.json"))

			pricingRaw, _ := ioutil.ReadFile(filepath.Join(serviceDir, "pricing.json"))
			pricing := map[string]int64{}
			if len(pricingRaw) > 0 {
				json.Unmarshal(pricingRaw, &pricing)
			}

			err = publishAPI(serviceName, string(dat), string(js), string(examples), pricing)
			if err != nil {
				fmt.Println("Failed to save data to publicapi service", err)
				os.Exit(1)
			}
		}
	}

}
