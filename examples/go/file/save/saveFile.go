package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/file"
)

// Save a file
func SaveFile() {
	fileService := file.NewFileService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := fileService.Save(&file.SaveRequest{
		File: &file.Record{
			Content: "file content example",
			Path:    "/document/text-files/file.txt",
			Project: "examples",
		},
	})
	fmt.Println(rsp, err)
}
