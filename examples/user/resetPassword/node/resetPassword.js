const { UserService } = require("m3o/user");

// Reset password with the code sent by the "SendPasswordResetEmail" endoint.
async function resetPassword() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.resetPassword({
    code: "012345",
    confirmPassword: "NewPassword1",
    email: "joe@example.com",
    newPassword: "NewPassword1",
  });
  console.log(rsp);
}

resetPassword();
