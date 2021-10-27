# User

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/User/api](https://m3o.com/User/api).

Endpoints:

## Read

Read an account by id, username or email. Only one need to be specified.


[https://m3o.com/user/api#Read](https://m3o.com/user/api#Read)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Read an account by id, username or email. Only one need to be specified.
func ReadAnAccountById() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Read(&user.ReadRequest{
		Id: "usrid-1",

	})
	fmt.Println(rsp, err)
}
```
## Read

Read an account by id, username or email. Only one need to be specified.


[https://m3o.com/user/api#Read](https://m3o.com/user/api#Read)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Read an account by id, username or email. Only one need to be specified.
func ReadAccountByUsernameOrEmail() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Read(&user.ReadRequest{
		Username: "usrname-1",

	})
	fmt.Println(rsp, err)
}
```
## Read

Read an account by id, username or email. Only one need to be specified.


[https://m3o.com/user/api#Read](https://m3o.com/user/api#Read)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Read an account by id, username or email. Only one need to be specified.
func ReadAccountByEmail() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Read(&user.ReadRequest{
		Email: "joe@example.com",

	})
	fmt.Println(rsp, err)
}
```
## VerifyEmail

Verify the email address of an account from a token sent in an email to the user.


[https://m3o.com/user/api#VerifyEmail](https://m3o.com/user/api#VerifyEmail)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Verify the email address of an account from a token sent in an email to the user.
func VerifyEmail() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.VerifyEmail(&user.VerifyEmailRequest{
		Token: "t2323t232t",

	})
	fmt.Println(rsp, err)
}
```
## Login

Login using username or email. The response will return a new session for successful login,
401 in the case of login failure and 500 for any other error


[https://m3o.com/user/api#Login](https://m3o.com/user/api#Login)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Login using username or email. The response will return a new session for successful login,
// 401 in the case of login failure and 500 for any other error
func LogAuserIn() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Login(&user.LoginRequest{
		Email: "joe@example.com",
Password: "mySecretPass123",

	})
	fmt.Println(rsp, err)
}
```
## Logout

Logout a user account


[https://m3o.com/user/api#Logout](https://m3o.com/user/api#Logout)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Logout a user account
func LogAuserOut() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Logout(&user.LogoutRequest{
		SessionId: "sds34s34s34-s34s34-s43s43s34-s4s34s",

	})
	fmt.Println(rsp, err)
}
```
## ReadSession

Read a session by the session id. In the event it has expired or is not found and error is returned.


[https://m3o.com/user/api#ReadSession](https://m3o.com/user/api#ReadSession)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Read a session by the session id. In the event it has expired or is not found and error is returned.
func ReadAsessionByTheSessionId() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.ReadSession(&user.ReadSessionRequest{
		SessionId: "sds34s34s34-s34s34-s43s43s34-s4s34s",

	})
	fmt.Println(rsp, err)
}
```
## Create

Create a new user account. The email address and username for the account must be unique.


[https://m3o.com/user/api#Create](https://m3o.com/user/api#Create)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Create a new user account. The email address and username for the account must be unique.
func CreateAnAccount() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Create(&user.CreateRequest{
		Email: "joe@example.com",
Id: "usrid-1",
Password: "mySecretPass123",
Username: "usrname-1",

	})
	fmt.Println(rsp, err)
}
```
## Update

Update the account username or email


[https://m3o.com/user/api#Update](https://m3o.com/user/api#Update)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Update the account username or email
func UpdateAnAccount() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Update(&user.UpdateRequest{
		Email: "joeotheremail@example.com",
Id: "usrid-1",

	})
	fmt.Println(rsp, err)
}
```
## UpdatePassword

Update the account password


[https://m3o.com/user/api#UpdatePassword](https://m3o.com/user/api#UpdatePassword)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Update the account password
func UpdateTheAccountPassword() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.UpdatePassword(&user.UpdatePasswordRequest{
		ConfirmPassword: "myEvenMoreSecretPass123",
NewPassword: "myEvenMoreSecretPass123",
OldPassword: "mySecretPass123",

	})
	fmt.Println(rsp, err)
}
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

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Send a verification email
// to the user being signed up. Email from will be from 'support@m3o.com',
// but you can provide the title and contents.
// The verification link will be injected in to the email as a template variable, $micro_verification_link.
// Example: 'Hi there, welcome onboard! Use the link below to verify your email: $micro_verification_link'
// The variable will be replaced with an actual url that will look similar to this:
// 'https://user.m3o.com/user/verify?token=a-verification-token&redirectUrl=your-redir-url'
func SendVerificationEmail() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.SendVerificationEmail(&user.SendVerificationEmailRequest{
		Email: "joe@example.com",
FailureRedirectUrl: "https://m3o.com/verification-failed",
FromName: "Awesome Dot Com",
RedirectUrl: "https://m3o.com",
Subject: "Email verification",
TextContent: `Hi there,

Please verify your email by clicking this link: $micro_verification_link`,

	})
	fmt.Println(rsp, err)
}
```
## Delete

Delete an account by id


[https://m3o.com/user/api#Delete](https://m3o.com/user/api#Delete)

```go
package example

import(
	"fmt"
	"os"

	"github.com/micro/services/clients/go/user"
)

// Delete an account by id
func DeleteUserAccount() {
	userService := user.NewUserService(os.Getenv("MICRO_API_TOKEN"))
	rsp, err := userService.Delete(&user.DeleteRequest{
		Id: "fdf34f34f34-f34f34-f43f43f34-f4f34f",

	})
	fmt.Println(rsp, err)
}
```
