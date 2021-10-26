package example

import (
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
