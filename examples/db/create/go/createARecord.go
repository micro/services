package example

import (
	"fmt"
	"github.com/micro/services/clients/go/db"
	"os"
)

// Create a record in the database. Optionally include an "id" field otherwise it's set automatically.
func CreateArecord() {
	dbService := db.NewDbService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := dbService.Create(&db.CreateRequest{
		Record: map[string]interface{}{
			"name":     "Jane",
			"age":      42,
			"isActive": true,
			"id":       "1",
		},
		Table: "users",
	})
	fmt.Println(rsp, err)
}
