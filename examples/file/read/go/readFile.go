package example

import (
	"fmt"
	"github.com/micro/services/clients/go/file"
	"os"
)

// Read a file by path
func ReadFile() {
	fileService := file.NewFileService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := fileService.Read(&file.ReadRequest{
		Path:    "/document/text-files/file.txt",
		Project: "examples",
	})
	fmt.Println(rsp, err)
}
