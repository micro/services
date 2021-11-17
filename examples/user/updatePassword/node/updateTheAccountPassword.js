const { UserService } = require("m3o/user");

// Update the account password
async function updateTheAccountPassword() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.updatePassword({
    confirmPassword: "Password2",
    id: "user-1",
    newPassword: "Password2",
    oldPassword: "Password1",
  });
  console.log(rsp);
}

updateTheAccountPassword();
