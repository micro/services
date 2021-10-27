# File

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/File/api](https://m3o.com/File/api).

Endpoints:

## Save

Save a file


[https://m3o.com/file/api#Save](https://m3o.com/file/api#Save)

```go
package example

import(
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
			Path: "/document/text-files/file.txt",
	Project: "examples",
	},

	})
	fmt.Println(rsp, err)
}
```
## List

List files by their project and optionally a path.


[https://m3o.com/file/api#List](https://m3o.com/file/api#List)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/file"
)

// List files by their project and optionally a path.
func ListFiles() {
	fileService := file.NewFileService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := fileService.List(&file.ListRequest{
		Project: "examples",

	})
	fmt.Println(rsp, err)
}
```
## Delete

Delete a file by project name/path


[https://m3o.com/file/api#Delete](https://m3o.com/file/api#Delete)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/file"
)

// Delete a file by project name/path
func DeleteFile() {
	fileService := file.NewFileService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := fileService.Delete(&file.DeleteRequest{
		Path: "/document/text-files/file.txt",
Project: "examples",

	})
	fmt.Println(rsp, err)
}
```
## Read

Read a file by path


[https://m3o.com/file/api#Read](https://m3o.com/file/api#Read)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/file"
)

// Read a file by path
func ReadFile() {
	fileService := file.NewFileService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := fileService.Read(&file.ReadRequest{
		Path: "/document/text-files/file.txt",
Project: "examples",

	})
	fmt.Println(rsp, err)
}
```
