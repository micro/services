package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/db"
)

// Count records in a table
func CountEntriesInAtable() {
	dbService := db.NewDbService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := dbService.Count(&db.CountRequest{
		Table: "users",
	})
	fmt.Println(rsp, err)
}
