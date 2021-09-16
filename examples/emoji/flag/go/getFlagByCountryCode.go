package example

import (
	"fmt"
	"github.com/micro/services/clients/go/emoji"
	"os"
)

// Get the flag for a country. Requires country code e.g GB for great britain
func GetFlagByCountryCode() {
	emojiService := emoji.NewEmojiService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := emojiService.Flag(&emoji.FlagRequest{})
	fmt.Println(rsp, err)
}
