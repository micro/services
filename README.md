# Micro Services

This repo includes reusable micro services.

## Overview

Services provides a home for real world examples for using Micro v3.

- [blog](blog) - A blog app composed as micro services
- [helloworld](helloworld) - A simple helloworld service
- [test](test) - A set of sample test services for Micro

## Usage

Pull the service directly from github

```
# install micro
go get github.com/micro/micro/v3

# run the server
micro server

# login with user: admin pass: micro
micro login

# run the service
micro run github.com/micro/services/helloworld
```

## Contributing

Feel free to contribute by PR and signoff.

## License

[Polyform Strict](https://polyformproject.org/licenses/strict/1.0.0/)

