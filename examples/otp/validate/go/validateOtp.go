package example

import (
	"fmt"
	"github.com/micro/services/clients/go/otp"
	"os"
)

// Validate the OTP code
func ValidateOtp() {
	otpService := otp.NewOtpService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := otpService.Validate(&otp.ValidateRequest{
		Code: "656211",
		Id:   "asim@example.com",
	})
	fmt.Println(rsp, err)
}
