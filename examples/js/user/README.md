# User

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/User/api](https://m3o.com/User/api).

Endpoints:

## Read

Read an account by id, username or email. Only one need to be specified.


[https://m3o.com/user/api#Read](https://m3o.com/user/api#Read)

```js
const { UserService } = require('m3o/user');

// Read an account by id, username or email. Only one need to be specified.
async function readAnAccountById() {
	let userService = new UserService(process.env.MICRO_API_TOKEN)
	let rsp = await userService.read({
  "id": "usrid-1"
})
	console.log(rsp)
}

readAnAccountById()
```
## Read

Read an account by id, username or email. Only one need to be specified.


[https://m3o.com/user/api#Read](https://m3o.com/user/api#Read)

```js
const { UserService } = require('m3o/user');

// Read an account by id, username or email. Only one need to be specified.
async function readAccountByUsernameOrEmail() {
	let userService = new UserService(process.env.MICRO_API_TOKEN)
	let rsp = await userService.read({
  "username": "usrname-1"
})
	console.log(rsp)
}

readAccountByUsernameOrEmail()
```
## Read

Read an account by id, username or email. Only one need to be specified.


[https://m3o.com/user/api#Read](https://m3o.com/user/api#Read)

```js
const { UserService } = require('m3o/user');

// Read an account by id, username or email. Only one need to be specified.
async function readAccountByEmail() {
	let userService = new UserService(process.env.MICRO_API_TOKEN)
	let rsp = await userService.read({
  "email": "joe@example.com"
})
	console.log(rsp)
}

readAccountByEmail()
```
## VerifyEmail

Verify the email address of an account from a token sent in an email to the user.


[https://m3o.com/user/api#VerifyEmail](https://m3o.com/user/api#VerifyEmail)

```js
const { UserService } = require('m3o/user');

// Verify the email address of an account from a token sent in an email to the user.
async function verifyEmail() {
	let userService = new UserService(process.env.MICRO_API_TOKEN)
	let rsp = await userService.verifyEmail({
  "token": "t2323t232t"
})
	console.log(rsp)
}

verifyEmail()
```
## Login

Login using username or email. The response will return a new session for successful login,
401 in the case of login failure and 500 for any other error


[https://m3o.com/user/api#Login](https://m3o.com/user/api#Login)

```js
const { UserService } = require('m3o/user');

// Login using username or email. The response will return a new session for successful login,
// 401 in the case of login failure and 500 for any other error
async function logAuserIn() {
	let userService = new UserService(process.env.MICRO_API_TOKEN)
	let rsp = await userService.login({
  "email": "joe@example.com",
  "password": "mySecretPass123"
})
	console.log(rsp)
}

logAuserIn()
```
## Logout

Logout a user account


[https://m3o.com/user/api#Logout](https://m3o.com/user/api#Logout)

```js
const { UserService } = require('m3o/user');

// Logout a user account
async function logAuserOut() {
	let userService = new UserService(process.env.MICRO_API_TOKEN)
	let rsp = await userService.logout({
  "sessionId": "sds34s34s34-s34s34-s43s43s34-s4s34s"
})
	console.log(rsp)
}

logAuserOut()
```
## ReadSession

Read a session by the session id. In the event it has expired or is not found and error is returned.


[https://m3o.com/user/api#ReadSession](https://m3o.com/user/api#ReadSession)

```js
const { UserService } = require('m3o/user');

// Read a session by the session id. In the event it has expired or is not found and error is returned.
async function readAsessionByTheSessionId() {
	let userService = new UserService(process.env.MICRO_API_TOKEN)
	let rsp = await userService.readSession({
  "sessionId": "sds34s34s34-s34s34-s43s43s34-s4s34s"
})
	console.log(rsp)
}

readAsessionByTheSessionId()
```
## Create

Create a new user account. The email address and username for the account must be unique.


[https://m3o.com/user/api#Create](https://m3o.com/user/api#Create)

```js
const { UserService } = require('m3o/user');

// Create a new user account. The email address and username for the account must be unique.
async function createAnAccount() {
	let userService = new UserService(process.env.MICRO_API_TOKEN)
	let rsp = await userService.create({
  "email": "joe@example.com",
  "id": "usrid-1",
  "password": "mySecretPass123",
  "username": "usrname-1"
})
	console.log(rsp)
}

createAnAccount()
```
## Update

Update the account username or email


[https://m3o.com/user/api#Update](https://m3o.com/user/api#Update)

```js
const { UserService } = require('m3o/user');

// Update the account username or email
async function updateAnAccount() {
	let userService = new UserService(process.env.MICRO_API_TOKEN)
	let rsp = await userService.update({
  "email": "joeotheremail@example.com",
  "id": "usrid-1"
})
	console.log(rsp)
}

updateAnAccount()
```
## UpdatePassword

Update the account password


[https://m3o.com/user/api#UpdatePassword](https://m3o.com/user/api#UpdatePassword)

```js
const { UserService } = require('m3o/user');

// Update the account password
async function updateTheAccountPassword() {
	let userService = new UserService(process.env.MICRO_API_TOKEN)
	let rsp = await userService.updatePassword({
  "confirmPassword": "myEvenMoreSecretPass123",
  "id": "usrid-1",
  "newPassword": "myEvenMoreSecretPass123",
  "oldPassword": "mySecretPass123"
})
	console.log(rsp)
}

updateTheAccountPassword()
```
## SendVerificationEmail

Send a verification email
to the user being signed up. Email from will be from 'support@m3o.com',
but you can provide the title and contents.
The verification link will be injected in to the email as a template variable, $micro_verification_link.
Example: 'Hi there, welcome onboard! Use the link below to verify your email: $micro_verification_link'
The variable will be replaced with an actual url that will look similar to this:
'https://user.m3o.com/user/verify?token=a-verification-token&redirectUrl=your-redir-url'


[https://m3o.com/user/api#SendVerificationEmail](https://m3o.com/user/api#SendVerificationEmail)

```js
const { UserService } = require('m3o/user');

// Send a verification email
// to the user being signed up. Email from will be from 'support@m3o.com',
// but you can provide the title and contents.
// The verification link will be injected in to the email as a template variable, $micro_verification_link.
// Example: 'Hi there, welcome onboard! Use the link below to verify your email: $micro_verification_link'
// The variable will be replaced with an actual url that will look similar to this:
// 'https://user.m3o.com/user/verify?token=a-verification-token&redirectUrl=your-redir-url'
async function sendVerificationEmail() {
	let userService = new UserService(process.env.MICRO_API_TOKEN)
	let rsp = await userService.sendVerificationEmail({
  "email": "joe@example.com",
  "failureRedirectUrl": "https://m3o.com/verification-failed",
  "fromName": "Awesome Dot Com",
  "redirectUrl": "https://m3o.com",
  "subject": "Email verification",
  "textContent": "Hi there,\n\nPlease verify your email by clicking this link: $micro_verification_link"
})
	console.log(rsp)
}

sendVerificationEmail()
```
## Delete

Delete an account by id


[https://m3o.com/user/api#Delete](https://m3o.com/user/api#Delete)

```js
const { UserService } = require('m3o/user');

// Delete an account by id
async function deleteUserAccount() {
	let userService = new UserService(process.env.MICRO_API_TOKEN)
	let rsp = await userService.delete({
  "id": "fdf34f34f34-f34f34-f43f43f34-f4f34f"
})
	console.log(rsp)
}

deleteUserAccount()
```
