import * as user from "@m3o/services/user";

// Delete an account by id
async function DeleteUserAccount() {
  let userService = new user.UserService(process.env.MICRO_API_TOKEN);
  let rsp = await userService.delete({
    id: "fdf34f34f34-f34f34-f43f43f34-f4f34f",
  });
  console.log(rsp);
}

await DeleteUserAccount();
