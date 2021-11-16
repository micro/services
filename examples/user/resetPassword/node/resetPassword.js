const { UserService } = require("m3o/user");

// Reset password with the code sent by the "SendPasswordResetEmail" endoint.
async function resetPassword() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.resetPassword({
    code: "some-code-from-email",
    confirmPassword: "newpass123",
    newPassword: "newpass123",
  });
  console.log(rsp);
}

resetPassword();
