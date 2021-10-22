package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/db"
)

// Update a record in the database. Include an "id" in the record to update.
func UpdateArecord() {
	dbService := db.NewDbService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := dbService.Update(&db.UpdateRequest{
		Record: map[string]interface{}{
			"id":  "1",
			"age": 43,
		},
		Table: "users",
	})
	fmt.Println(rsp, err)
}
