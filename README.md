# micronomicon

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/5a491d39cd6a4456954bd79dd883b7da)](https://www.codacy.com/manual/MensurOwary/micronomicon?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=MensurOwary/micronomicon&amp;utm_campaign=Badge_Grade)

The backend for the micro-learning app _**micronomicon**_. 

Usage:

There are many available tags. User subscribes to new tags and gets a _micron_ to read/watch on those topics. 
A Micron is a random resource on the given topics.

## Endpoints

| method    | endpoint          | does                              |
|:---       |:---               | :---                              |
| POST      | /register         | Register a user                   |                   
| POST      | /login            | Log a user in                     |                   
| POST      | /logout           | Logs the user out                 |                   
| GET       | /users/me         | Get the user                      |          
| GET       | /users/me/tags    | Get the user's tags               |   
| POST      | /users/me/tags    | Add new tags to the user tags     |   
| DELETE    | /users/me/tags    | Delete tags from the user tags    |  
| GET       | /users/me/microns | Get a micron for the user         |
| GET       | /tags             | Get all the available tags        |

### `/register`

<details>
  <summary>Expand</summary>

*Method* : **POST**

*Body*

```json
{
    "email": "jane@doe.com",
    "username": "jane",
    "name": "Jane Doe",
    "password": "abc-12345"
}
```

*Response*

```json
{
    "message": "Created the user"
}
```
    
</details>

### `/login`
<details>
  <summary>Expand</summary>

*Method* : **POST**

*Body*
```json
{
    "username" : "jane",
    "password" : "abc-12345"
}
```

*Response*
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY1NzE2MjIsImlhdCI6MTU5NjU2NDQyMiwidXNlcm5hbWUiOiJvd2FyeSJ9.CS7osQC8bQihpG7nQKWqUvGQysiB9Y0lkV7tdXSoS-o"
}
```
</details>

### `/logout`
<details>
  <summary>Expand</summary>

*Method* : **POST**

*Headers*
```text
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY1NzE2MjIsImlhdCI6MTU5NjU2NDQyMiwidXNlcm5hbWUiOiJvd2FyeSJ9.CS7osQC8bQihpG7nQKWqUvGQysiB9Y0lkV7tdXSoS-o
```

*Response*

```json
{
    "message": "Successfully logged out"
}
```
</details>

### `/users/me`

<details>
  <summary>Expand</summary>

*Method* : **GET**

*Headers*
```text
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY1NzE2MjIsImlhdCI6MTU5NjU2NDQyMiwidXNlcm5hbWUiOiJvd2FyeSJ9.CS7osQC8bQihpG7nQKWqUvGQysiB9Y0lkV7tdXSoS-o
```

*Response*

```json
{
    "username": "jane",
    "name": "Jane Doe",
    "email": "jane@hello.com",
    "tags": {
        "tags": [
            {
                "name": "c"
            },
            {
                "name": "react"
            },
            {
                "name": "python"
            }
        ],
        "size": 3
    }
}
```

</details>

### `/users/me/tags` : GET
<details>
  <summary>Expand</summary>

*Method* : **GET**

*Headers*
```text
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY1NzE2MjIsImlhdCI6MTU5NjU2NDQyMiwidXNlcm5hbWUiOiJvd2FyeSJ9.CS7osQC8bQihpG7nQKWqUvGQysiB9Y0lkV7tdXSoS-o
```

*Response*

```json
{
    "tags": [
        {
            "name": "c"
        },
        {
            "name": "react"
        },
        {
            "name": "python"
        }
    ],
    "size": 3
}
```

</details>

### `/users/me/tags` : POST
<details>
  <summary>Expand</summary>

*Method* : **POST**

*Headers*
```text
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY1NzE2MjIsImlhdCI6MTU5NjU2NDQyMiwidXNlcm5hbWUiOiJvd2FyeSJ9.CS7osQC8bQihpG7nQKWqUvGQysiB9Y0lkV7tdXSoS-o
```

*Body*
```json
{
    "ids": ["c", "angular", "react"]
}
```

*Response*

```json
{
    "message": "Successfully added the tags"
}
```

</details>

### `/users/me/tags` : DELETE
<details>
  <summary>Expand</summary>

*Method* : **DELETE**

*Headers*
```text
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY1NzE2MjIsImlhdCI6MTU5NjU2NDQyMiwidXNlcm5hbWUiOiJvd2FyeSJ9.CS7osQC8bQihpG7nQKWqUvGQysiB9Y0lkV7tdXSoS-o
```

*Body*
```json
{
    "ids": ["c", "angular", "react"]
}
```

*Response*

```json
{
    "message": "Successfully removed the tags"
}
```

</details>

### `/tags`

<details>
  <summary>Expand</summary>

*Method* : **GET**

*Headers*
```text
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY1NzE2MjIsImlhdCI6MTU5NjU2NDQyMiwidXNlcm5hbWUiOiJvd2FyeSJ9.CS7osQC8bQihpG7nQKWqUvGQysiB9Y0lkV7tdXSoS-o
```

*Response*

```json
{
    "tags": [
        {
            "name": "c"
        },
        {
            "name": "react"
        },
        {
            "name": "groovy"
        }
    ]
}
```
</details>

### `/users/me/microns`

<details>
  <summary>Expand</summary>

*Method* : **GET**

*Headers*
```text
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY1NzE2MjIsImlhdCI6MTU5NjU2NDQyMiwidXNlcm5hbWUiOiJvd2FyeSJ9.CS7osQC8bQihpG7nQKWqUvGQysiB9Y0lkV7tdXSoS-o
```

*Response*

```json
{
    "URL": "https://learnxinyminutes.com/docs/python/",
    "Title": "Python"
}
```

</details>

## Running tests

To run all the test, issue the following command:

```shell script
docker-compose -f docker-compose-test.yml up
```
