# Emoji

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Emoji/api](https://m3o.com/Emoji/api).

Endpoints:

## Find

Find an emoji by its alias e.g :beer:


[https://m3o.com/emoji/api#Find](https://m3o.com/emoji/api#Find)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/emoji"
)

// Find an emoji by its alias e.g :beer:
func FindEmoji() {
	emojiService := emoji.NewEmojiService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := emojiService.Find(&emoji.FindRequest{
		Alias: ":beer:",

	})
	fmt.Println(rsp, err)
}
```
## Flag

Get the flag for a country. Requires country code e.g GB for great britain


[https://m3o.com/emoji/api#Flag](https://m3o.com/emoji/api#Flag)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/emoji"
)

// Get the flag for a country. Requires country code e.g GB for great britain
func GetFlagByCountryCode() {
	emojiService := emoji.NewEmojiService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := emojiService.Flag(&emoji.FlagRequest{
		
	})
	fmt.Println(rsp, err)
}
```
## Print

Print text and renders the emojis with aliases e.g
let's grab a :beer: becomes let's grab a 🍺


[https://m3o.com/emoji/api#Print](https://m3o.com/emoji/api#Print)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/emoji"
)

// Print text and renders the emojis with aliases e.g
// let's grab a :beer: becomes let's grab a 🍺
func PrintTextIncludingEmoji() {
	emojiService := emoji.NewEmojiService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := emojiService.Print(&emoji.PrintRequest{
		Text: "let's grab a :beer:",

	})
	fmt.Println(rsp, err)
}
```
## Send

Send an emoji to anyone via SMS. Messages are sent in the form '<message> Sent from <from>'


[https://m3o.com/emoji/api#Send](https://m3o.com/emoji/api#Send)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/emoji"
)

// Send an emoji to anyone via SMS. Messages are sent in the form '<message> Sent from <from>'
func SendAtextContainingAnEmojiToAnyoneViaSms() {
	emojiService := emoji.NewEmojiService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := emojiService.Send(&emoji.SendRequest{
		From: "Alice",
Message: "let's grab a :beer:",
To: "+44782669123",

	})
	fmt.Println(rsp, err)
}
```
