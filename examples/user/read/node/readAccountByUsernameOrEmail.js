const { UserService } = require("micro-js-client/user");

// Read an account by id, username or email. Only one need to be specified.
async function readAccountByUsernameOrEmail() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.read({
    username: "joe",
  });
  console.log(rsp);
}

readAccountByUsernameOrEmail();
