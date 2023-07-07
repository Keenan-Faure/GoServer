# GoServer

A Web Server built with Golang

This project serves both local files on the browser and has an API (documentation below) and was created from a guide
found on the [Boot.Dev](https://boot.dev) course.

Please visit the Prerequisites below before running the program.

## Motivation

This project serves as the baby steps in my learning of Golang, so please do check it out if you need to determine my
current (at this time) knowledge of Golang. However, since being introduced to Golang by Boot.Dev it has been great
and very interesting, so I wish to go from baby steps into giant strides from here.

### Goal

The main goal in this project is very simple: it serves as a start in my learning of this amazing language.

## Main Features

- Uses [Chi routing](https://github.com/go-chi/chi) for the API
- Reading and writing from a json file (acts as the database) using mutex locks
- Environment variables to store JWT Tokens and API Keys
- Reading and Responding with JSON over HTTP with typical HTTP Requests (GET, POST, PUT etc.)
- A command line flag `--debug` more information below:
- Using JWT package to create, parse, and validate a JWT Token
  - Use of [jwt.io](https://jwt.io) to examine JWT data
  - Use of Refresh Tokens to get new JWT (access) Tokens
  - Revoking Refresh Tokens through API requests
- Authentication and Request Authorization
- Documented API endpoints

## How to install the program

You can download the `.zip` file on the repository page

Or if you prefer:

```
gh repo clone Keenan-Faure/Integration-web-app
```

Then download the dependencies of the program:

```
go mod download
```

## How to run the program

Once in the directory of the project in the command line enter:

```
go build -o out && ./out
```

## How is my data I create over the API Saved

Note all the data is written and read from a `database.json` file found at the root of the product and will be created
once the program starts

## Prerequisites: Creating a environment variables:

Create a file called `.env` in the directory of the project and populate it with the format below:

```
JWT_SECRET=your-secret-key
WEBHOOK_TOKEN=your-token
```

Replace `your-secret-key` with a long, random string. Below shows some steps on how to generate random tokens
Replace `your-token` with a long, random string. You can generate one on the command line like this:

### How to generate a Token

```
openssl rand -base64 64
```

# API Documentation

**Note that Authorization needs to be added as headers in the respective requests**

## Chirp Resource

```json
{
  "id": 1,
  "body": "test string"
}
```

**Max length is 140 characters**

### GET /api/chirps

Returns an array of users

Params:

```json
{
  "key": "author_id",
  "value": 1
}
```

```json
{
  "key": "sort",
  "value": "asc"
}
```

Where `author_id` is the id of the user who created the chirp
The value of `sort` is `asc` by default but can be adjusted to `desc`

### GET /api/chirps/{id}

Returns a specific chirp

Response Example:

```json
{
  "id": 1,
  "body": "example text",
  "author_id": 1
}
```

### DELETE /api/chirps

Removes a specific chirp

Authorization:

```
Authorization: Bearer <token>
```

## User Resource

```json
{
  "id": 1,
  "email": "test@gmail.com",
  "password": "abc123"
}
```

### POST /api/users

Creates a new user

Request Body:

```json
{
  "email": "test@gmail.com",
  "password": "abc123"
}
```

Response Example:

```json
{
  "id": 1,
  "email": "test23@gmail.com"
}
```

### POST /api/login

Login as a user

Request Body:

```json
{
  "email": "test@gmail.com",
  "password": "abc123"
}
```

Response Body:

```json
{
  "id": 1,
  "email": "test23@gmail.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiMSIsImV4cCI6MTY4ODQ2NTYwMCwiaWF0IjoxNjg4NDYyMDAwfQ.hx7bVdtYldScCQl_kDxw9NOF2MbE_tEkPcd-vZ76xaY",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktcmVmcmVzaCIsInN1YiI6IjEiLCJleHAiOjE2ODg1NDg0MDAsImlhdCI6MTY4ODQ2MjAwMH0.tlZ-CG2APGbO-h7RH02IexFYG8OP9b3UekITFlQLUnI"
}
```

## PUT /api/users

Updates a user

Request Body:

```json
{
  "email": "testNew@gmail.com",
  "password": "newPassword"
}
```

Response Body:

```json
{
  "id": 1,
  "email": "testNew@gmail.com"
}
```

Authorization:

```
Authorization: Bearer <token>
```

### GET /api/healthz

Returns the status of the API

### POST /api/refresh

Requests new access tokens

Response Body:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiMSIsImV4cCI6MTY4ODQ2NTYwMCwiaWF0IjoxNjg4NDYyMDAwfQ.hx7bVdtYldScCQl_kDxw9NOF2MbE_tEkPcd-vZ76xaY"
}
```

Authorization:

```
Authorization: Bearer <refreshToken>
```

### POST /api/revoke

Revokes the access token found in the header

Authorization:

```
Authorization: Bearer <token>
```

### POST /polka/webhooks

Accepts webhook notifications from Polka

Request Body:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": 3
  }
}
```

Authorization:

```
Authorization: ApiKey <key>
```
