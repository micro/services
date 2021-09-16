import * as sms from "@m3o/services/sms";

// Send an SMS.
async function SendSms() {
  let smsService = new sms.SmsService(process.env.MICRO_API_TOKEN);
  let rsp = await smsService.send({
    from: "Alice",
    message: "Hi there!",
    to: "+447681129",
  });
  console.log(rsp);
}

await SendSms();
