# This is the case project given by Acquisition.ai

## How to run locally?
To run this locally, you will need to change a few things:
- Go into backend directory, and change the env.example to .env, the credentials inside will expire in a week (15/09/2024)
- Run ```make build``` ```make run-auth``` and ```make run-weather``` in backend diroctory to build and start the backend services
- The auth and weather services are exposed to ports 8080, and 8081 respectively
- Switch to frontend directory, and run ```make build``` ```make run``` to start the react project locally
- This project makes use of cookies, make sure to run project from ```127.0.0.1:5173/```

## What is this project?
- Login, Logout, Register, GetWeather, UpdateWeather functionality.
- Simple APIs using both query parameters and json body.
- Clean and simple UI for pages Login, Register, and Dashboard

## Improvements:
- Haven't made the front-end responsive.
- Having better go structs for bson.M data

