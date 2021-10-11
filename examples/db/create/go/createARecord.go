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
			"id":       "1",
			"name":     "Jane",
			"age":      42,
			"isActive": true,
		},
		Table: "users",
	})
	fmt.Println(rsp, err)
}
