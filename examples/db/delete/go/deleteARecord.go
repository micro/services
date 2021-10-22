package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/db"
)

// Delete a record in the database by id.
func DeleteArecord() {
	dbService := db.NewDbService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := dbService.Delete(&db.DeleteRequest{
		Id:    "1",
		Table: "users",
	})
	fmt.Println(rsp, err)
}
