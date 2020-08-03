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
| GET       | /users/me         | Get the user                      |          
| GET       | /users/me/tags    | Get the user's tags               |   
| POST      | /users/me/tags    | Add new tags to the user tags     |   
| DELETE    | /users/me/tags    | Delete tags from the user tags    |  
| GET       | /users/me/microns | Get a micron for the user         |
| GET       | /tags             | Get all the available tags        |


### Running tests

To run all the test, issue the following command:

```shell script
docker-compose -f docker-compose-test.yml up
```
