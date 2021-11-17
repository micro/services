const { UserService } = require("micro-js-client/user");

// Send an email with a verification code to reset password.
// Call "ResetPassword" endpoint once user provides the code.
async function sendPasswordResetEmail() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.sendPasswordResetEmail({
    email: "joe@example.com",
    fromName: "Awesome Dot Com",
    subject: "Password reset",
    textContent:
      "Hi there,\n click here to reset your password: myapp.com/reset/code?=$code",
  });
  console.log(rsp);
}

sendPasswordResetEmail();
