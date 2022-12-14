# go-web-app-boilerplate

A REST API Boilerplate written in Go.


[![Go Report Card](https://goreportcard.com/badge/github.com/ablancas22/messenger-backend)](https://goreportcard.com/report/github.com/ablancas22/messenger-backend)

* Author(s): John Connor Sanders
* Current Version: 0.0.1
* Release Date: 9/25/2022
* MIT License
___
## Getting Started

Follow the instructions below to get the Go REST API up and running on your local environment

### Prerequisites

* MongoDB 4+
* Go 1.18+

### Setup

1. Create a conf.json file from the example file and configure the following settings:

* MongoDB URI Connection String
* Secret Encryption String
* Mongo Database Name
* Master Admin Username
* Master Admin Email
* Master Admin Initial Password
* Whether to run App with HTTPS
* If HTTPS is on, the cert.pem file
* If HTTPS is on, the path to the key.pem file
* Whether you want new users to be able to sign themselves up for accounts
* Run ENV

2. Use the provided install.sh script to build a background service

```bash
$ cp conf.json.example conf.json
$ vi conf.json
```
___
## Running the API

### Production

One way to set up and start a production build is to run the following:

* Build executable and install the SystemD service:
```bash
$ go build github.com/ablancas22/messenger-backend
$ sh ./scripts/setup_service.sh
```

* To start the API:
```bash
$ sh ./scripts/start.sh
```

* To stop the API:
```bash
$ sh ./scripts/stop.sh
```

### Development

To start the API in development:
```bash
$ go run github.com/ablancas22/messenger-backend
```

To stop the development API, enter 'ctrl + c'

### Testing

1. Integration Test
```bash
$ go test github.com/ablancas22/messenger-backend/cmd
```

2. Unit Tests
* Test Auth Module:
```bash
$ go test github.com/ablancas22/messenger-backend/auth
```

* Test Database Module:
```bash
$ go test github.com/ablancas22/messenger-backend/database
```

* Test Models Module:
```bash
$ go test github.com/ablancas22/messenger-backend/models
```

___
## API Route Guide
### I) Authentication Routes

___
#### 1. Signin
* POST - /auth

##### Request

***
* Headers

```
{
  Content-Type: application/json
}
```

* Body
```
{
  "email": "email@example.com",
  "password": "userpass"
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Auth-Token: "",
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH
}
```

* Body
```
{
  "id": "000000000000000000000011",
  "username": "userName",
  "firstname": "john",
  "lastname": "smith",
  "email": "user@example.com",
  "role": "member",
  "group_id": 000000000000000000000001",
  "last_modified": 2019-06-07 20:17:14.630917778 +0000 UTC,
  "created_at": 2019-06-07 20:17:14.630917778 +0000 UTC
}
```

#### 2. Signup
* POST - /auth/register
* This route will return a 404 if the "Registration" setting is set to "off" in the conf.json file.

##### Request

***
* Headers

```
{
  Content-Type: application/json
}
```

* Body
```
{
  "firstname": "john",
  "lastname": "smith",
  "email": "user@example.com"",
  "username": "userName",
  "password": "userpass"
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Auth-Token: "",
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
{
  "id": "000000000000000000000012",
  "username": "userName",
  "firstname": "john",
  "lastname": "smith",
  "email": "user@example.com",
  "role": "member",
  "groupuuid": "000000000000000000000002",
  "last_modified": 2019-06-07 20:17:14.630917778 +0000 UTC,
  "created_at": 2019-06-07 20:17:14.630917778 +0000 UTC
}
```

#### 3. Refresh Token
* GET - /auth

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Auth-Token: "",
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

#### 4. Signout
* DELETE - /auth

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH
}
```

#### 5. API Key - Expires after 6 Months
* GET - /auth/api-key

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Auth-Token: "",
  API-Key: "",
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

#### 6. Update Password
* POST - /auth/password

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

***
* Body

```
{
  "current_password": "current_password",
  "new_password": "new_password"
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

### II) Task Routes

___
#### 1. List Task(s)
* GET - /tasks/{taskId}
* todoId parameter is optional, if used the request will only return an object for that item.

##### Request

***
* Headers
```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers
```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
[
  {
    "id":  "000000000000000000000021",
    "name": "todo_name",
    "due": 2019-08-01 12:04:01 -0000 UTC,
    "completed": false,
    "description": "Task to complete",
    "user_id": "000000000000000000000011",
    "group_id": "000000000000000000000001",
    "last_modified": 2019-06-07 20:28:09.400248747 +0000 UTC,
    "created_at": 2019-06-07 20:28:09.400248747 +0000 UTC
  }
]
```

#### 2. Create Task
* POST - /tasks

##### Request

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Auth-Token: ""
}
```

* Body
```
{
    "name": "todo_name",
    "due": 2019-08-01 12:04:01 -0000 UTC,
    "description": "Task to complete"
}
```


##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
{
   "id":  "000000000000000000000021",
   "name": "todo_name",
   "due": 2019-08-01 12:04:01 -0000 UTC,
   "completed": false,
   "description": "Task to complete",
   "user_id": "000000000000000000000011",
   "group_id": "000000000000000000000001",
   "last_modified": 2019-06-07 20:28:09.400248747 +0000 UTC,
   "created_at": 2019-06-07 20:28:09.400248747 +0000 UTC
}
```

#### 3. Modify Todo
* PATCH - /todos/{todoId}

##### Request

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Auth-Token: ""
}
```

* Body
```
{
    "name": "new_todo_name",
    "due": 2019-08-06 12:04:01 -0000 UTC,
    "description": "Updated Task to complete",
    "completed": true,
    "user_id": "000000000000000000000011"
}
```


##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
{
   "id":  "000000000000000000000022",
   "name": "new_todo_name",
   "due": 2019-08-01 12:04:01 -0000 UTC,
   "completed": true,
   "description": "Task to complete",
   "user_id": "000000000000000000000011",
   "group_id": "000000000000000000000001",
   "last_modified": 2019-06-07 20:28:09.400248747 +0000 UTC,
   "created_at": 2019-06-07 20:28:09.400248747 +0000 UTC
}
```

#### 4. Delete Todo
* DELETE - /todos/{todoId}

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

### III) Users Routes (Admins Only)

___
#### 1. List User(s)
* GET - /users/{userId}
* userId parameter is optional, if used the request will only return an object for that item.

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
[
  {
    "id": "000000000000000000000011",
    "username": "userName",
    "firstname": "jane",
    "lastname": "smith",
    "email": "user@example.com",
    "role": "member",
    "group_id": "000000000000000000000001",
    "last_modified": 2019-06-07 20:17:14.630917778 +0000 UTC,
    "created_at": 2019-06-07 20:17:14.630917778 +0000 UTC
  }
]
```


#### 2. Create User
* POST - /users

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

* Body
```
{
  "username": "userName",
  "firstname": "jane",
  "lastname": "smith",
  "email": "user@example.com",
  "password": "xyz789",
  "role": "member",
  "last_modified": 2019-06-07 20:17:14.630917778 +0000 UTC,
  "created_at": 2019-06-07 20:17:14.630917778 +0000 UTC
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
{
  "id": "000000000000000000000012",
  "username": "userName",
  "firstname": "jane",
  "lastname": "smith",
  "email": "user@example.com",
  "role": "member",
  "group_id": "000000000000000000000001",
  "last_modified": 2019-06-07 20:17:14.630917778 +0000 UTC,
  "created_at": 2019-06-07 20:17:14.630917778 +0000 UTC
}
```

#### 3. Modify User
* PATCH - /users/{userId}

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

* Body
```
{
  "id": "000000000000000000000012",
  "username": "newUserName",
  "firstname": "jane",
  "lastname": "smith",
  "email": "new_test@email.com",
  "password": "newUserpass",
  "role": "member",
  "last_modified": 2019-06-07 20:17:14.630917778 +0000 UTC,
  "created_at": 2019-06-07 20:17:14.630917778 +0000 UTC
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
{
  "id": "000000000000000000000012",
  "username": "newUserName",
  "firstname": "jane",
  "lastname": "smith",
  "email": "new_test@email.com",
  "role": "member",
  "group_id": "000000000000000000000001",
  "last_modified": 2019-06-07 20:17:14.630917778 +0000 UTC,
  "created_at": 2019-06-07 20:17:14.630917778 +0000 UTC
}
```


#### 4. Delete User
* DELETE - /users/{userId}

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

### IV) User Group Routes (Admins Only)

___
#### 1. List User Group(s)
* GET - /groups/{groupId}
* groupId parameter is optional, if used the request will only return an object for that item.

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
[
  {
    "id": 000000000000000000000001,
    "name": "groupName",        
    "last_modified": 2019-06-07 20:17:14.358617998 +0000 UTC,
    "creation_datetime": 2019-06-07 20:17:14.358617998 +0000 UTC
  }
]
```


#### 2. Create User Group
* POST - /groups

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

* Body
```
{
  "Name": "newGroup"
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
{
  "id": "000000000000000000000002",
  "name": "newGroup",    
  "last_modified": 2019-06-07 20:18:15.145971952 +0000 UTC,
  "creation_datetime": 2019-06-07 20:18:15.145971952 +0000 UTC
}
```


#### 3. Modify User Group
* PATCH - /groups/{groupId}

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

* Body
```
{
  "Name": "newGroupName"
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```

* Body
```
{
  "id": "000000000000000000000002",
  "name": "newGroup",    
  "last_modified": 2019-06-07 20:18:15.145971952 +0000 UTC,
  "creation_datetime": 2019-06-07 20:18:15.145971952 +0000 UTC
}
```

#### 4. Delete User Group
* DELETE - /groups/{groupId}

##### Request

***
* Headers

```
{
  Content-Type: application/json,
  Auth-Token: ""
}
```

##### Response

***
* Headers

```
{
  Content-Type: application/json; charset=UTF-8,
  Date: DoW, DD MMM YYYY HH:mm:SS GMT,
  Content-Length: 0,
  Access-Control-Allow-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Expose-Headers: Content-Type, Auth-Token, API-Key,
  Access-Control-Allow-Origin: *,
  Access-Control-Allow-Methods: GET,DELETE,POST,PATCH  
}
```