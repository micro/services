package example

import (
	"fmt"
	"github.com/micro/services/clients/go/emoji"
	"os"
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
