# micronomicon

The backend for the micro-learning app _**micronomicon**_. 

Usage:

There are many available tags. User subscribes to new tags and gets a _micron_ to read/watch on those topics. 
A Micron is a random resource on the given topics

## Endpoints

| method  | endpoint  | does |
|:---|:---| :--- |
| POST   | /register | Register a user |                   
| POST   | /login | Log a user in|                   
| GET    | /users/me | Get the user |          
| GET    | /users/me/tags | Get the user's tags |   
| POST   | /users/me/tags | Add new tags to the user tags |   
| DELETE | /users/me/tags | Delete tags from the user tags |  
| GET    | /users/me/microns | Get a micron for the user |
| GET    | /tags | Get all the available tags |
