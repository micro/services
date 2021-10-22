const { UserService } = require("m3o/user");

// Logout a user account
async function logAuserOut() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.logout({
    sessionId: "sds34s34s34-s34s34-s43s43s34-s4s34s",
  });
  console.log(rsp);
}

logAuserOut();
