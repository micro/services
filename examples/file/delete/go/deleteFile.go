package example

import (
	"fmt"
	"github.com/micro/services/clients/go/file"
	"os"
)

// Delete a file by project name/path
func DeleteFile() {
	fileService := file.NewFileService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := fileService.Delete(&file.DeleteRequest{
		Path:    "/document/text-files/file.txt",
		Project: "examples",
	})
	fmt.Println(rsp, err)
}
