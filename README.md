# How to run the server

Server is tested on Ubuntu 22.04

## Requirements

1. go 1.19x
2. The rest of the requirements are in the `go.mod` file, dependencies will be automatically be installed

## Steps

1. Create a new database named `govtech_assessment` using this command:

```
CREATE DATABASE govtech_assessment;
```

2. Create a new user in mysql named `govtech` with password `%00B5Rp8DV` using this command (not recommended to use the same password in this README, this is just to make it easier to set up the db and run the server locally):

```
CREATE USER 'govtech'@'localhost' IDENTIFIED BY '00B5Rp8DV';
```

3. Grant permissions to the database for the user: govtech

```
GRANT ALL ON govtech_assessment.* TO 'govtech'@'localhost';
```

4. Copy and paste the shared `.env` file into the project's root directory.

5. If any of the database details were changed, update the `.env` file accordingly.

6. Uncomment line 34 in `main.go` if you wish to seed the database with some dummy data. run `go run main.go` to seed the database. Comment back the line afterwards.

7. Use this command to run the server:

```
go run main.go
```

8. To access the protected APIs, the user needs to login first with an account. If you had seeded the database earlier, you can login with any of the teacher's credentials. Example:

```
{
  "email": "teacher1@gmail.com",
  "password": "Password1!"
}
```

9. You will receive a response containing the access and refresh tokens. The access token will be needed to gain access to the protected APIs.

10. Paste the access token under the Authorization header in the form of:

```
Bearer <access token>
```

11. Now you will be able to access the protected API endpoints until the access token expire.

# How to run the tests

1. Assuming that your current working directory is at the project's root directory, use this command to run the tests:

```
go test -v ./...;
```
