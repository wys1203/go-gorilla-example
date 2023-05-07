# go-gorilla-example

This repository contains an example of a backend API server built using the Gorilla toolkit and the gorilla/mux package in Go. The server demonstrates basic CRUD operations and can be easily extended to accommodate additional functionality.

## Requirements

- Go 1.16+
- openssl

### Development key pair

In the future, we plan to use asymmetric encryption for encrypting and validating JSON Web Tokens (JWTs). Here is step by step to generate a key pair by `openssl`

1. Generate a private key (RSA 2048 bit):

```
openssl genrsa -out private_key.pem 2048
```

2. Extract the public key from the private key:

```
openssl rsa -in private_key.pem -pubout -out public_key.pem
```

## Getting Started

To get started with this project, follow these steps:

1. Clone the repository:

```
git clone https://github.com/wys1203/go-gorilla-example.git
cd go-gorilla-example
```

2. Run the server with docker-compose:

```
make docker-up
```

By default, the server will start on port 8080. You can access the API endpoints using any HTTP client, like curl or Postman.

## API Endpoints

Here's a list of the example API endpoints provided by this server:

```
POST /signup: Create a new user.
POST /signin: User login and retrieve JWT token

GET /users: Retrieve a list of users. (JWT Auth need)
GET /users/search: Retrieve a specific user by given fullname (JWT Auth need)
GET /users/{acct}: Retrieve a specific user by Account. (JWT Auth need)
DELETE /users/{acct}: Delete a specific user by Account. (JWT Auth need)
PATCH /users/{acct}: Update a specific user by Account. (JWT Auth need)
PUT /users/{acct}/fullname: Update a specific user fullname by Account. (JWT Auth need)
```

## Enable security

- Enable HTTPS
- CSRF
- XSS

### Enable HTTPS

To enable HTTPS, you need create a self-signed certificate using the openssl command.

```
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

Please use command-line tools such as curl may require the -k or --insecure flag to allow connections to a server using a self-signed certificate.

### CSRF

```
go get -u github.com/gorilla/csrf
```

## License

This project is licensed under the MIT License. For more information
