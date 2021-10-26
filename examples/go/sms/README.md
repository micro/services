# Sms

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Sms/api](https://m3o.com/Sms/api).

Endpoints:

## Send

Send an SMS.


[https://m3o.com/sms/api#Send](https://m3o.com/sms/api#Send)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/sms"
)

// Send an SMS.
func SendSms() {
	smsService := sms.NewSmsService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := smsService.Send(&sms.SendRequest{
		From: "Alice",
Message: "Hi there!",
To: "+447681129",

	})
	fmt.Println(rsp, err)
}
```
