const { UserService } = require("micro-js-client/user");

// Delete an account by id
async function deleteUserAccount() {
  let userService = new UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.delete({
    id: "8b98acbe-0b6a-4d66-a414-5ffbf666786f",
  });
  console.log(rsp);
}

deleteUserAccount();
