const { UserService } = require("m3o/user");

// Update the account username or email
async function updateAnAccount() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.update({
    email: "joeotheremail@example.com",
    id: "usrid-1",
  });
  console.log(rsp);
}

updateAnAccount();
