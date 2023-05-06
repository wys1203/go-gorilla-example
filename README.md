# go-gorilla-example

This is a simple project use gorillatoolkit build a API server

## Development key pair

In the future, we plan to use asymmetric encryption for encrypting and validating JSON Web Tokens (JWTs). Here is step by step to generate a key pair by `openssl`

1. Generate a private key (RSA 2048 bit):

```
openssl genrsa -out private_key.pem 2048
```

2. Extract the public key from the private key:

```
openssl rsa -in private_key.pem -pubout -out public_key.pem
```
