import * as user from "m3o/user";

// Create a new user account. The email address and username for the account must be unique.
async function CreateAnAccount() {
  let userService = new user.UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.create({
    email: "joe@example.com",
    id: "usrid-1",
    password: "mySecretPass123",
    username: "usrname-1",
  });
  console.log(rsp);
}

await CreateAnAccount();
