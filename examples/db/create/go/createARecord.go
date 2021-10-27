package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/db"
)

// Create a record in the database. Optionally include an "id" field otherwise it's set automatically.
func CreateArecord() {
	dbService := db.NewDbService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := dbService.Create(&db.CreateRequest{
		Record: map[string]interface{}{
			"isActive": true,
			"id":       "1",
			"name":     "Jane",
			"age":      42,
		},
		Table: "users",
	})
	fmt.Println(rsp, err)
}
