package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	indexTSFile = `
export * from "./{{.ServiceName}}_grpc_web_pb";
export * from "./{{.ServiceName}}_pb";
`

	indexJSFile = `
module.exports = Object.assign(
  {},
  require("./{{.ServiceName}}_pb"),
  require("./{{.ServiceName}}_grpc_web_pb")
);
`
)

func main() {
	var protos []string

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".proto") {
			protos = append(protos, path)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range protos {
		// resolves "blog/comments/proto/comments.proto" to "comments"
		serviceName := strings.TrimSuffix(filepath.Base(p), ".proto")
		dir := filepath.Dir(p)
		jsDir := filepath.Join(wd, "clients", "js", serviceName)

		if _, err := os.Stat(jsDir); os.IsNotExist(err) {
			err = os.MkdirAll(jsDir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Println("Generating Go, Micro, TS and JS for " + serviceName)
		c := exec.Command("protoc",
			serviceName+".proto",
			"--go_out=plugins=grpc,paths=source_relative:.",
			"--micro_out=paths=source_relative:.",
			"--js_out=import_style=commonjs,binary:"+jsDir,
			"--grpc-web_out=import_style=commonjs+dts,mode=grpcweb:"+jsDir,
		)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Dir = filepath.Join(wd, dir)
		c.Env = os.Environ()
		err := c.Run()
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(filepath.Join(jsDir, "index.d.ts"))
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}

		t := template.Must(template.New("indexTSFile").Parse(indexTSFile))
		err = t.Execute(f, struct{ ServiceName string }{ServiceName: serviceName})
		if err != nil {
			panic(err)
		}

		ft, err := os.Create(filepath.Join(jsDir, "index.js"))
		defer ft.Close()
		if err != nil {
			log.Fatal(err)
		}

		t = template.Must(template.New("indexJSFile").Parse(indexJSFile))
		err = t.Execute(ft, struct{ ServiceName string }{ServiceName: serviceName})
		if err != nil {
			panic(err)
		}
	}

	log.Println("Done!")
}
