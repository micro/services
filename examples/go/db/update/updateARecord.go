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
			"age": 43,
			"id":  "1",
		},
		Table: "users",
	})
	fmt.Println(rsp, err)
}
