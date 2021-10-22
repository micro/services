const { OtpService } = require("m3o/otp");

// Validate the OTP code
async function validateOtp() {
  let otpService = new OtpService(process.env.MICRO_API_TOKEN);
  let rsp = await otpService.validate({
    code: "656211",
    id: "asim@example.com",
  });
  console.log(rsp);
}

validateOtp();
