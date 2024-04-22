# Assignment_two

## Countries Dashboard Service

## Contributors
Solveig Langbakk

Trygve Sollund

Sofia Serine Mikkelsen
***

# General Setup
Following, examples for terminal commands are given below. Make sure to follow each step!

## Environment setup 
First copy the .env.example file to .env
```bash
 cp .env.example .env
```

Then create a directory for storing secrets called .secrets
```bash
mkdir .secrets
```


Then copy the .env file to the .secrets folder. 
The .secrets/.env file is used in the docker compose runtime environment.
This will therefore let you, for example, run different firebase credentials in the container and locally.
```bash
 cp .env .secrets/.env
```

## Firebase 
This project uses firebase and firestore.
This then requires you to have a firebase project with a corresponding service account key. 
[Here is a guide from Google on this](https://firebase.google.com/docs/app-distribution/authenticate-service-account?platform=ios) (you can ignore the `GOOGLE_APPLICATION_CREDENTIALS` step since we don't use it). 

The service account key should be stored in the `.secrets` folder,
and the path to it needs to be specified in the `.env` file. 
This also needs to be specified in the `.secrets/.env` file as well if you want to run it in **Docker**.  
## Docker setup
### First you need to start the docker engine. Start the docker engine:
#### On linux 
```bash
 sudo systemctl start docker
 ```

#### On Windows or MacOS
Start the `docker desktop` app 


Pull golang
```bash
docker pull golang
```
Run the container:
```bash
docker compose -f compose.yml up -d
```
To see running containers use
```bash
docker ps 
```


# Endpoints

```
/dashboard/v1/registrations/
/dashboard/v1/dashboards/
/dashboard/v1/notifications/
/dashboard/v1/status/
```

**Registrations endpoint**

GET and POST

```/dashboard/v1/registrations/```

GET, PUT, and DELETE

```/dashboard/v1/registrations/{id}```

**Notifications endpoint**

```/dashboard/v1/notifications/```

```/dashboard/v1/notifications/{id}```

**Status endpoint**

```dashboard/v1/status/```

## Services used in this project

*REST Countries API (instance hosted for this course)*

Endpoint: http://129.241.150.113:8080/v3.1

Documentation: http://129.241.150.113:8080/

*Open-Meteo APIs (hosted externally, hence please be responsible)*

Documentation: https://open-meteo.com/en/features#available-apis

*Currency API*

Endpoint: http://129.241.150.113:9090/currency/

Documentation: http://129.241.150.113:9090/

## Task
In this group assignment, we have to developed a REST web application in Golang 
that provides the client with the ability to configure information dashboards that 
are dynamically populated when requested. The dashboard configurations are saved in 
the service in a persistent way (they should survive service restart),
and populated based on external services. It also includes a simple notification 
service that can listen to specific events. The application will be dockerized and 
deployed using an IaaS system.
