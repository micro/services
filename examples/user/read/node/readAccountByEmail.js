import * as user from "m3o/user";

// Read an account by id, username or email. Only one need to be specified.
async function ReadAccountByEmail() {
  let userService = new user.UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.read({
    email: "joe@example.com",
  });
  console.log(rsp);
}

await ReadAccountByEmail();
