const { UserService } = require("m3o/user");

// Read a session by the session id. In the event it has expired or is not found and error is returned.
async function readAsessionByTheSessionId() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.readSession({
    sessionId: "sds34s34s34-s34s34-s43s43s34-s4s34s",
  });
  console.log(rsp);
}

readAsessionByTheSessionId();
