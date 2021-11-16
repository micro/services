package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/db"
)

// Drop a table in the DB
func DropTable() {
	dbService := db.NewDbService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := dbService.DropTable(&db.DropTableRequest{
		Table: "users",
	})
	fmt.Println(rsp, err)
}
