package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/db"
)

//
func ListTables() {
	dbService := db.NewDbService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := dbService.ListTables(&db.ListTablesRequest{})
	fmt.Println(rsp, err)
}
