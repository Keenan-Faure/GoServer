# GoServer

A web server built with Golang

incorporates:

-   Chi routing ()
-   API Endpoint generation
-   Database as a file (.json)
-   Environment variables
-   Reading a Responding with JSON over the API
-   GET, POST, PUT Requests
-   Cmd flags in go (--debug)
-   API Authentication with JWT
-   Using JWT package to create/parse/validate a JWT

Creating an environment variable:

-   create a file called `.env` in the directory of the project
-   populate it with the format below:

```
JWT_SECRET=your-secret-key
```

-   Replace your-secret-key with a long, random string. You can generate one on the command line like this:

```
openssl rand -base64 64
```
