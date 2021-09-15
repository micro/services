import * as user from "@m3o/services/user";

// Read an account by id, username or email. Only one need to be specified.
async function ReadAnAccountById() {
  let userService = new user.UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.read({
    id: "usrid-1",
  });
  console.log(rsp);
}

await ReadAnAccountById();
