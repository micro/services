const { SpamService } = require("micro-js-client/spam");

// Check whether an email is likely to be spam based on its attributes
async function classifyAnEmail() {
  let spamService = new SpamService(process.env.MICRO_API_TOKEN);
  let rsp = await spamService.classify({
    from: "noreply@m3o.com",
    subject: "Welcome",
    text_body: "Hi there,\n\nWelcome to M3O.\n\nThanks\nM3O team",
    to: "hello@example.com",
  });
  console.log(rsp);
}

classifyAnEmail();
