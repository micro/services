import * as user from "m3o/user";

// Read a session by the session id. In the event it has expired or is not found and error is returned.
async function ReadAsessionByTheSessionId() {
  let userService = new user.UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.readSession({
    sessionId: "sds34s34s34-s34s34-s43s43s34-s4s34s",
  });
  console.log(rsp);
}

await ReadAsessionByTheSessionId();
