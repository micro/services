const { UserService } = require("micro-js-client/user");

// Update the account username or email
async function updateAnAccount() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.update({
    email: "joe+2@example.com",
    id: "user-1",
  });
  console.log(rsp);
}

updateAnAccount();
