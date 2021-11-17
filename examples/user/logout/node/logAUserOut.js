const { UserService } = require("micro-js-client/user");

// Logout a user account
async function logAuserOut() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.logout({
    sessionId: "df91a612-5b24-4634-99ff-240220ab8f55",
  });
  console.log(rsp);
}

logAuserOut();
