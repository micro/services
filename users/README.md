# Users Service

A user service for storing accounts and simple auth.

## Getting started

```
micro run github.com/micro/services/users
```

## Usage

User server implements the following RPC Methods

Users
- Create
- Read
- Update
- Delete
- Search
- UpdatePassword
- Login
- Logout
- ReadSession


### Create

```shell
micro call users Users.Create '{"user":{"id": "ff3c06de-9e43-41c7-9bab-578f6b4ad32b", "username": "asim", "email": "asim@example.com"}, "password": "password1"}'
```

### Read

```shell
micro call users Users.Read '{"id": "ff3c06de-9e43-41c7-9bab-578f6b4ad32b"}'
```

### Update

```shell
micro call users Users.Update '{"user":{"id": "ff3c06de-9e43-41c7-9bab-578f6b4ad32b", "username": "asim", "email": "asim+update@example.com"}}'
```

### Update Password

```shell
micro call users Users.UpdatePassword '{"userId": "ff3c06de-9e43-41c7-9bab-578f6b4ad32b", "oldPassword": "password1", "newPassword": "newpassword1", "confirmPassword": "newpassword1" }'
```

### Delete

```shell
micro call users Users.Delete '{"id": "ff3c06de-9e43-41c7-9bab-578f6b4ad32b"}'
```

### Login

```shell
micro call users Users.Login '{"username": "asim", "password": "password1"}'
```

### Read Session

```shell
micro call users Users.ReadSession '{"sessionId": "sr7UEBmIMg5hYOgiljnhrd4XLsnalNewBV9KzpZ9aD8w37b3jRmEujGtKGcGlXPg1yYoSHR3RLy66ugglw0tofTNGm57NrNYUHsFxfwuGC6pvCn8BecB7aEF6UxTyVFq"}'
```

### Logout

```shell
micro call users Users.Logout '{"sessionId": "sr7UEBmIMg5hYOgiljnhrd4XLsnalNewBV9KzpZ9aD8w37b3jRmEujGtKGcGlXPg1yYoSHR3RLy66ugglw0tofTNGm57NrNYUHsFxfwuGC6pvCn8BecB7aEF6UxTyVFq"}'
```
