import * as user from "@m3o/services/user";

// Read an account by id, username or email. Only one need to be specified.
async function ReadAccountByUsernameOrEmail() {
  let userService = new user.UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.read({
    username: "usrname-1",
  });
  console.log(rsp);
}

await ReadAccountByUsernameOrEmail();
