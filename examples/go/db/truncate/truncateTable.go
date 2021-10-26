package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/db"
)

// Truncate the records in a table
func TruncateTable() {
	dbService := db.NewDbService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := dbService.Truncate(&db.TruncateRequest{
		Table: "users",
	})
	fmt.Println(rsp, err)
}
