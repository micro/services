import * as user from "m3o/user";

// Update the account password
async function UpdateTheAccountPassword() {
  let userService = new user.UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.updatePassword({
    confirmPassword: "myEvenMoreSecretPass123",
    id: "usrid-1",
    newPassword: "myEvenMoreSecretPass123",
    oldPassword: "mySecretPass123",
  });
  console.log(rsp);
}

await UpdateTheAccountPassword();
