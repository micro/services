const { UserService } = require("m3o/user");

// Verify the email address of an account from a token sent in an email to the user.
async function verifyEmail() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.verifyEmail({
    token: "t2323t232t",
  });
  console.log(rsp);
}

verifyEmail();
