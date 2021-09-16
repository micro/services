import * as otp from "m3o/otp";

// Validate the OTP code
async function ValidateOtp() {
  let otpService = new otp.OtpService(process.env.MICRO_API_TOKEN);
  let rsp = await otpService.validate({
    code: "656211",
    id: "asim@example.com",
  });
  console.log(rsp);
}

await ValidateOtp();
