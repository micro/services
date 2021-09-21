import * as user from "m3o/user";

// Logout a user account
async function LogAuserOut() {
  let userService = new user.UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.logout({
    sessionId: "sds34s34s34-s34s34-s43s43s34-s4s34s",
  });
  console.log(rsp);
}

await LogAuserOut();
