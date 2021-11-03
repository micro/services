package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
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

type PublicAPI struct {
	Name         string           `json:"name"`
	Category     string           `json:"category,omitempty"`
	Description  string           `json:"description"`
	Icon         string           `json:"icon,omitempty"`
	OpenAPIJson  string           `json:"open_api_json"`
	Pricing      map[string]int64 `json:"pricing,omitempty"`
	ExamplesJson string           `json:"examples_json,omitempty"`
	PostmanJson  string           `json:"postman_json,omitempty"`
	DisplayName  string           `json:"display_name,omitempty"`
	Quotas       map[string]int64 `json:"quotas,omitempty"`
}

const (
	prodAPIDomain = "api.m3o.com"
)

func publishAPI(apiSpec *PublicAPI, domain string) error {
	client := &http.Client{}

	//Encode the data
	postBody, _ := json.Marshal(map[string]interface{}{
		"api": apiSpec,
	})

	rbody := bytes.NewBuffer(postBody)

	//Leverage Go's HTTP Post function to make request
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/publicapi/Publish", domain), rbody)

	// Add auth headers here if needed
	req.Header.Add("Authorization", `Bearer `+os.Getenv("MICRO_ADMIN_TOKEN"))
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(b))
	}

	io.Copy(ioutil.Discard, resp.Body)

	return nil
}

func main() {
	workDir, _ := os.Getwd()
	domainFlag := flag.String("domain", prodAPIDomain, "domain to publish to e.g. api.m3o.com")
	serviceFlag := flag.String("service", "", "individual service to publish e.g. helloworld")
	flag.Parse()

	files, err := ioutil.ReadDir(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.Contains(f.Name(), "clients") || strings.Contains(f.Name(), "examples") {
			continue
		}
		if len(*serviceFlag) > 0 && f.Name() != *serviceFlag {
			continue
		}
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
			makeProto := exec.Command("make", "api")
			makeProto.Dir = serviceDir
			fmt.Println(serviceDir)
			outp, err := makeProto.CombinedOutput()
			if err != nil {
				fmt.Println("Failed to make api", string(outp))
				os.Exit(1)
			}

			serviceName := f.Name()

			// generate the Postman collection
			postman := exec.Command("openapi2postmanv2", "-s", fmt.Sprintf("api-%s.json", serviceName), "-o", "postman.json")
			postman.Dir = serviceDir
			outp, err = postman.CombinedOutput()
			if err != nil {
				fmt.Printf("Failed to generate postman collection %s %s\n", string(outp), err)
				os.Exit(1)
			}

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

			// we have to read an openapi spec otherwise we can't publish
			if err := json.Unmarshal(js, &spec); err != nil {
				fmt.Println("Failed to unmarshal", err)
				os.Exit(1)
			}

			// define the default public api values
			publicApi := new(PublicAPI)

			// if we find a public api definition we load it
			if b, err := ioutil.ReadFile(filepath.Join(serviceDir, "publicapi.json")); err == nil {
				// unpack the info if we read the file
				json.Unmarshal(b, &publicApi)
			}

			// If we didn't get the default info from a file, populate it
			if publicApi.Name == "" {
				publicApi.Name = serviceName
			}
			if publicApi.Description == "" {
				publicApi.Description = string(dat)
			}
			if publicApi.OpenAPIJson == "" {
				publicApi.OpenAPIJson = string(js)
			}

			// load the examples if they exist
			if examples, err := ioutil.ReadFile(filepath.Join(serviceDir, "examples.json")); err == nil {
				if len(examples) > 0 {
					publicApi.ExamplesJson = string(examples)
				}
			}

			// load the separate pricing if it exists
			if pricingRaw, err := ioutil.ReadFile(filepath.Join(serviceDir, "pricing.json")); err == nil {
				pricing := map[string]int64{}
				// unmarshal the pricing info
				if len(pricingRaw) > 0 {
					json.Unmarshal(pricingRaw, &pricing)
					publicApi.Pricing = pricing
				}
			}

			// load the postman json
			if postman, err := ioutil.ReadFile(filepath.Join(serviceDir, "postman.json")); err == nil {
				if len(postman) > 0 {
					publicApi.PostmanJson = string(postman)
				}
			}

			// publish the api
			if err := publishAPI(publicApi, *domainFlag); err != nil {
				fmt.Println("Failed to save data to publicapi service", err)
				os.Exit(1)
			}
		}
	}

}
