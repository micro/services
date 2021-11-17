const { UserService } = require("m3o/user");

// Read a session by the session id. In the event it has expired or is not found and error is returned.
async function readAsessionByTheSessionId() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.readSession({
    sessionId: "df91a612-5b24-4634-99ff-240220ab8f55",
  });
  console.log(rsp);
}

readAsessionByTheSessionId();
