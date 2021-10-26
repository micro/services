package example

import (
	"fmt"
	"os"

	"github.com/micro/services/clients/go/otp"
)

// Generate an OTP (one time pass) code
func GenerateOtp() {
	otpService := otp.NewOtpService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := otpService.Generate(&otp.GenerateRequest{
		Id: "asim@example.com",
	})
	fmt.Println(rsp, err)
}
