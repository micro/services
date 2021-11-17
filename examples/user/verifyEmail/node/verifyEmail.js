const { UserService } = require("micro-js-client/user");

// Verify the email address of an account from a token sent in an email to the user.
async function verifyEmail() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.verifyEmail({
    email: "joe@example.com",
    token: "012345",
  });
  console.log(rsp);
}

verifyEmail();
