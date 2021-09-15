import * as user from "@m3o/services/user";

// Login using username or email. The response will return a new session for successful login,
// 401 in the case of login failure and 500 for any other error
async function LogAuserIn() {
  let userService = new user.UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.login({
    email: "joe@example.com",
    password: "mySecretPass123",
  });
  console.log(rsp);
}

await LogAuserIn();
