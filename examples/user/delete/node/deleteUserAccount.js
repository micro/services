const { UserService } = require("m3o/user");

// Delete an account by id
async function deleteUserAccount() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.delete({
    id: "fdf34f34f34-f34f34-f43f43f34-f4f34f",
  });
  console.log(rsp);
}

deleteUserAccount();
