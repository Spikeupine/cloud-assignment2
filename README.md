# Assignment_two

## Countries Dashboard Service

## Contributors
Solveig Langbakk

Trygve Sollund

Sofia Serine Mikkelsen
***

# Docker setup
First you need to make sure you have the firebase key saved in `firebase_privatekey.json` at the project root.

## Start the docker engine
### On linux 
```bash
 sudo systemctl start docker
 ```

### On Windows or MacOS
Start `docker desktop` from programs


Pull golang
```bash
docker pull golang
```
Build the container:
```bash
docker build -t dashboardservice .
```
Run the container:
```bash
docker run -d -p 8080:8080 --name dashboard-service dashboardservice
```


# Endpoints

```
/dashboard/v1/registrations/
/dashboard/v1/dashboards/
/dashboard/v1/notifications/
/dashboard/v1/status/
```

**Registrations endpoint**

```/dashboard/v1/registrations/```
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
