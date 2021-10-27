# Sunnah

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Sunnah/api](https://m3o.com/Sunnah/api).

Endpoints:

## Books

Get a list of books from within a collection. A book can contain many chapters
each with its own hadiths.


[https://m3o.com/sunnah/api#Books](https://m3o.com/sunnah/api#Books)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/sunnah"
)

// Get a list of books from within a collection. A book can contain many chapters
// each with its own hadiths.
func GetTheBooksWithinAcollection() {
	sunnahService := sunnah.NewSunnahService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := sunnahService.Books(&sunnah.BooksRequest{
		Collection: "bukhari",

	})
	fmt.Println(rsp, err)
}
```
## Chapters

Get all the chapters of a given book within a collection.


[https://m3o.com/sunnah/api#Chapters](https://m3o.com/sunnah/api#Chapters)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/sunnah"
)

// Get all the chapters of a given book within a collection.
func ListTheChaptersInAbook() {
	sunnahService := sunnah.NewSunnahService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := sunnahService.Chapters(&sunnah.ChaptersRequest{
		Book: 1,
Collection: "bukhari",

	})
	fmt.Println(rsp, err)
}
```
## Hadiths

Hadiths returns a list of hadiths and their corresponding text for a
given book within a collection.


[https://m3o.com/sunnah/api#Hadiths](https://m3o.com/sunnah/api#Hadiths)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/sunnah"
)

// Hadiths returns a list of hadiths and their corresponding text for a
// given book within a collection.
func ListTheHadithsInAbook() {
	sunnahService := sunnah.NewSunnahService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := sunnahService.Hadiths(&sunnah.HadithsRequest{
		Book: 1,
Collection: "bukhari",

	})
	fmt.Println(rsp, err)
}
```
## Collections

Get a list of available collections. A collection is
a compilation of hadiths collected and written by an author.


[https://m3o.com/sunnah/api#Collections](https://m3o.com/sunnah/api#Collections)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/sunnah"
)

// Get a list of available collections. A collection is
// a compilation of hadiths collected and written by an author.
func ListAvailableCollections() {
	sunnahService := sunnah.NewSunnahService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := sunnahService.Collections(&sunnah.CollectionsRequest{
		
	})
	fmt.Println(rsp, err)
}
```
