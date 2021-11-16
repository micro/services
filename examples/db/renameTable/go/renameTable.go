package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/db"
)

// Rename a table
func RenameTable() {
	dbService := db.NewDbService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := dbService.RenameTable(&db.RenameTableRequest{
		From: "events",
		To:   "events_backup",
	})
	fmt.Println(rsp, err)
}
