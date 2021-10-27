# Function

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Function/api](https://m3o.com/Function/api).

Endpoints:

## Delete

Delete a function by name


[https://m3o.com/function/api#Delete](https://m3o.com/function/api#Delete)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/function"
)

// Delete a function by name
func DeleteAfunction() {
	functionService := function.NewFunctionService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := functionService.Delete(&function.DeleteRequest{
		Name: "my-first-func",
Project: "tests",

	})
	fmt.Println(rsp, err)
}
```
## Describe

Get the info for a deployed function


[https://m3o.com/function/api#Describe](https://m3o.com/function/api#Describe)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/function"
)

// Get the info for a deployed function
func DescribeFunctionStatus() {
	functionService := function.NewFunctionService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := functionService.Describe(&function.DescribeRequest{
		Name: "my-first-func",
Project: "tests",

	})
	fmt.Println(rsp, err)
}
```
## Deploy

Deploy a group of functions


[https://m3o.com/function/api#Deploy](https://m3o.com/function/api#Deploy)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/function"
)

// Deploy a group of functions
func DeployAfunction() {
	functionService := function.NewFunctionService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := functionService.Deploy(&function.DeployRequest{
		Entrypoint: "helloworld",
Name: "my-first-func",
Project: "tests",
Repo: "github.com/m3o/nodejs-function-example",
Runtime: "nodejs14",

	})
	fmt.Println(rsp, err)
}
```
## Call

Call a function by name


[https://m3o.com/function/api#Call](https://m3o.com/function/api#Call)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/function"
)

// Call a function by name
func CallAfunction() {
	functionService := function.NewFunctionService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := functionService.Call(&function.CallRequest{
		Name: "my-first-func",
Request: map[string]interface{}{
},

	})
	fmt.Println(rsp, err)
}
```
## List

List all the deployed functions


[https://m3o.com/function/api#List](https://m3o.com/function/api#List)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/function"
)

// List all the deployed functions
func ListFunctions() {
	functionService := function.NewFunctionService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := functionService.List(&function.ListRequest{
		
	})
	fmt.Println(rsp, err)
}
```
