# Sms

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Sms/api](https://m3o.com/Sms/api).

Endpoints:

## Send

Send an SMS.


[https://m3o.com/sms/api#Send](https://m3o.com/sms/api#Send)

```js
const { SmsService } = require('m3o/sms');

// Send an SMS.
async function sendSms() {
	let smsService = new SmsService(process.env.MICRO_API_TOKEN)
	let rsp = await smsService.send({
  "from": "Alice",
  "message": "Hi there!",
  "to": "+447681129"
})
	console.log(rsp)
}

sendSms()
```
