# Email

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Email/api](https://m3o.com/Email/api).

Endpoints:

## Send

Send an email by passing in from, to, subject, and a text or html body


[https://m3o.com/email/api#Send](https://m3o.com/email/api#Send)

```js
const { EmailService } = require('m3o/email');

// Send an email by passing in from, to, subject, and a text or html body
async function sendEmail() {
	let emailService = new EmailService(process.env.MICRO_API_TOKEN)
	let rsp = await emailService.send({
  "from": "Awesome Dot Com",
  "subject": "Email verification",
  "textBody": "Hi there,\n\nPlease verify your email by clicking this link: $micro_verification_link"
})
	console.log(rsp)
}

sendEmail()
```
