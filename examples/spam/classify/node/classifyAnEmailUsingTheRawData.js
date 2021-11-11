const { SpamService } = require("m3o/spam");

// Check whether an email is likely to be spam based on its attributes
async function classifyAnEmailUsingTheRawData() {
  let spamService = new SpamService(process.env.MICRO_API_TOKEN);
  let rsp = await spamService.classify({
    email_body:
      "Subject: Welcome\r\nTo: hello@emaple.com\r\nFrom: noreply@m3o.com\r\n\r\nHi there,\n\nWelcome to M3O.\n\nThanks\nM3O team",
  });
  console.log(rsp);
}

classifyAnEmailUsingTheRawData();
