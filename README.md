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

## Registrations endpoint


### Register new dashboard configuration

```http request
Method: POST
Path: /dashboard/v1/registrations/
Content type: application/json
```
Body
```json
{
   "country": "Norway",                     
   "isoCode": "NO",                          
   "features": {
      "temperature": true,               
      "precipitation": true,             
      "capital": true,                   
      "coordinates": true,               
      "population": true,                
      "area": true,                      
      "targetCurrencies": ["DKK", "SEK"] 
   }
}
```
Response
```json
{
    "id": "5a602bab-a921-4d38-9173-c1981ff02822",
    "lastChange": "2024-02-29 12:31"
}

```


### View a specific registered dashboard configuration


```
Method: GET
Path: /dashboard/v1/registrations/{id}
```
Response
```json
{
   "id": "5a602bab-a921-4d38-9173-c1981ff02822",
   "country": "Norway",
   "isoCode": "NO",
   "features": {
      "temperature": true,
      "precipitation": true,
      "capital": true,
      "coordinates": true,
      "population": true,
      "area": false,
      "targetCurrencies": ["EUR", "USD", "SEK"]
   },
    "lastChange": "20240229 14:07"
}
```



### View all registered dashboard configurations
```
Method: GET
Path: /dashboard/v1/registrations/
```
Response
```json
[
   {
      "id":"5a602bab-a921-4d38-9173-c1981ff02822",
      "country": "Norway",
      "isoCode": "NO",
      "features": {
                     "temperature": true,
                     "precipitation": true,
                     "capital": true,
                     "coordinates": true,
                     "population": true,
                     "area": false,
                     "targetCurrencies": ["EUR", "USD", "SEK"]
                  }, 
      "lastChange": "20240229 14:07"
   },
   {
      "id": "56bf8e8e-330b-4f96-ac59-35cea3db897d",
      "country": "Denmark",
      "isoCode": "DK",
      "features": {
                     "temperature": false,
                     "precipitation": true,
                     "capital": true,
                     "coordinates": true,
                     "population": false,
                     "area": true,
                     "targetCurrencies": ["NOK", "MYR", "JPY", "EUR"]
                  },
       "lastChange": "20240224 08:27"
   },
   ...
]

```


### Replace a specific registered dashboard configuration

```
Method: PUT
Path: /dashboard/v1/registrations/{id}
```
Body
```json
{
   "country": "Norway",
   "isoCode": "NO",
   "features": {
                  "temperature": false, 
                  "precipitation": true,
                  "capital": true,
                  "coordinates": true, 
                  "population": true,
                  "area": false,
                  "targetCurrencies": ["EUR", "SEK"] 
               }
}

```

### Delete a specific registered dashboard configuration

```
Method: DELETE
Path: /dashboard/v1/registrations/{id}
```
Response 
```
Status code: 200 
```

## Dashboard endpoint

```http request
Method: GET
Path: /dashboard/v1/dashboards/{id}
```
Response 
```json
{
   "id": "02073924-0015-11ef-9560-00059a3c7a00"
   "country": "Norway",
   "isoCode": "NO",
   "features": {
      "temperature": -1.2,                       
      "precipitation": 0.80,                     
      "capital": "Oslo",                         
      "coordinates": {
        "latitude": 62.0,
        "longitude": 10.0
      },
      "population": 5379475,
      "area": 323802.0,
      "targetCurrencies": {
        "EUR": 0.087701435,  
        "USD": 0.095184741, 
        "SEK": 0.97827275
      }
   },
  "lastRetrieval": "20240229 18:15" 
}

```



## Notifications endpoint

### Registration of Webhook
```
Method: POST
Path: /dashboard/v1/notifications/
Content type: application/json
```
Events are:

`REGISTER`

`CHANGE`

`DELETE`

`INVOKE`


Body
```json
{
   "url": "https://localhost:8080/client/",  // URL to be invoked when event occurs
   "country": "NO",                          // Country that is registered, or empty if all countries
   "event": "INVOKE"                         // Event on which it is invoked
}
```
response
```js
{
    "id": "OIdksUDwveiwe"
}
```


### Deletion of Webhook
Method: DELETE
Path: /dashboard/v1/notifications/{id}


### View specific registered webhook
```
Method: GET
Path: /dashboard/v1/notifications/{id}

```
```json
{
"id": "OIdksUDwveiwe",
"url": "https://localhost:8080/client/",
"country": "NO",
"event": "INVOKE"
}
```

### View all registered webhooks
```
Method: GET
Path: /dashboard/v1/notifications/
```
Body
```json
[
   {
      "id": "OIdksUDwveiwe",
      "url": "https://localhost:8080/client/",
      "country": "NO",
      "event": "INVOKE"
   },
   {
      "webhook_id": "DiSoisivucios",
      "url": "https://localhost:8081/anotherClient/",
      "country": "",                                 // field can also be omitted if registered for all countries
      "event": "REGISTER"
   },
   ...
]
```

## Status endpoint

```
Method: GET
Path: dashboard/v1/status/
```
Response
```json
{
   "countries_api": <http status code for *REST Countries API*>,
   "meteo_api": <http status code for *Meteo API*>, 
   "currency_api": <http status code for *Currency API*>,
   "notification_db": <http status code for *Notification database*>,
   "webhooks": <number of registered webhooks>,
   "version": "v1",
   "uptime": <time in seconds from the last service restart>
}

```

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
