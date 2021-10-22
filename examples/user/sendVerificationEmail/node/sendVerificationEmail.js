const { UserService } = require("m3o/user");

// Send a verification email
// to the user being signed up. Email from will be from 'support@m3o.com',
// but you can provide the title and contents.
// The verification link will be injected in to the email as a template variable, $micro_verification_link.
// Example: 'Hi there, welcome onboard! Use the link below to verify your email: $micro_verification_link'
// The variable will be replaced with an actual url that will look similar to this:
// 'https://user.m3o.com/user/verify?token=a-verification-token&redirectUrl=your-redir-url'
async function sendVerificationEmail() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.sendVerificationEmail({
    email: "joe@example.com",
    failureRedirectUrl: "https://m3o.com/verification-failed",
    fromName: "Awesome Dot Com",
    redirectUrl: "https://m3o.com",
    subject: "Email verification",
    textContent:
      "Hi there,\n\nPlease verify your email by clicking this link: $micro_verification_link",
  });
  console.log(rsp);
}

sendVerificationEmail();
