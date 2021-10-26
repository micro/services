const { UserService } = require("m3o/user");

// Update the account password
async function updateTheAccountPassword() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.updatePassword({
    confirmPassword: "myEvenMoreSecretPass123",
    id: "usrid-1",
    newPassword: "myEvenMoreSecretPass123",
    oldPassword: "mySecretPass123",
  });
  console.log(rsp);
}

updateTheAccountPassword();
