import * as otp from "m3o/otp";

// Generate an OTP (one time pass) code
async function GenerateOtp() {
  let otpService = new otp.OtpService(process.env.MICRO_API_TOKEN);
  let rsp = await otpService.generate({
    id: "asim@example.com",
  });
  console.log(rsp);
}

await GenerateOtp();
