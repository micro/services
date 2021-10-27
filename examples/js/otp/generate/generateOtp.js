const { OtpService } = require("m3o/otp");

// Generate an OTP (one time pass) code
async function generateOtp() {
  let otpService = new OtpService(process.env.MICRO_API_TOKEN);
  let rsp = await otpService.generate({
    id: "asim@example.com",
  });
  console.log(rsp);
}

generateOtp();
