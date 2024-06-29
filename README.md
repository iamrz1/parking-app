# Parking App

Parking App manges parking vehicle in lots. It has APIs for mangers and users.

- A manager can add parking lots with slots, and change maintenance mode of parking slots.
- A manager is created when the APP runs for the first time. You can login using the public login endpoint.
- Default manager credentials are username: `admin` and password: `admin`
- User can be registered with public registration endpoint 
- Users can access parking lots, and book/unbook parking space (one parking space per user).

We assume that this App will run behind a gateway and will not require stricter authorization
beyond JWT verification.



See bellow for API details.

## How to run
Parking App doesn't have any external dependency and can run on most machines out of the box.
This server has been tested on MacOS 12.1 Monterey, using go version go1.19 darwin/amd64.

Build and run this server:
```shell
$ make build
$ make run
```
This will create a binary and run it on local machine.

The server can be also run without creating an explicit binary:
```shell
$ make test_run
```

In any case, a config.yaml is used as the configuration file. It contains the port numbers, and database config.

Test cases has also been added for the APIs that can be run with the following command:
```shell
$ make test 
```

## API Doc:

This Signature server includes swagger for api documentation. It is exposed at `/doc/` endpoint (ie. http://localhost:8080/doc/index.html#/ for a local instance).

Here's a small rundown of the APIs:

## Public APIs

### Login:
```shell
curl --location 'http://localhost:8080/api/v1/public/login' \
--header 'Content-Type: application/json' \
--data '{
    "Username": "admin",
    "Password": "admin"
}'
```

Response:
```json
{
    "status": "OK",
    "message": "Logged In",
    "success": true,
    "data": {
        "AccessToken": "access-token",
        "RefreshToken": "refresh-token"
    },
    "timestamp": "2024-06-30T07:04:03.282Z"
}
```

NB: Default manager credentials are username: `admin` and password: `admin`

You will receive Access Token by logging in successfully.
Manager's token can be used to access manger endpoints only
User's token can be used to access user endpoints only
### Register as user

```shell
curl --location 'http://localhost:8080/api/v1/public/register' \
--header 'Content-Type: application/json' \
--data '{
    "Name": "Rezoan",
    "Username": "rezoan",
    "Password": "1234"
}'
```

Response:

```json
{
    "status": "Created",
    "message": "Registered User",
    "success": true,
    "data": {
        "ID": 3,
        "Username": "rezoan",
        "Name": "Rezoan"
    },
    "timestamp": "2024-06-30T07:05:28.013Z"
}
```


## Manager APIs:

You need to obtain manger's access token by logging in as manager using the login endpoint

### Create a parking lot with n slots:

```shell
curl --location 'http://localhost:8080/api/v1/manager/parking-lots' \
--header 'Authorization: <manager-token>' \
--header 'Content-Type: application/json' \
--data '{
    "Name": "Lot A",
    "NumberOfSlots": 2
}'
```

Response:
```json
{
  "status": "Created",
  "message": "Created Parking Lot",
  "success": true,
  "data": {
    "ID": 1,
    "Name": "Lot A",
    "Slots": [
      {
        "ID": 1,
        "UnderMaintenance": false,
        "Booked": false,
        "BookedAt": "0001-01-01T00:00:00Z"
      },
      {
        "ID": 2,
        "UnderMaintenance": false,
        "Booked": false,
        "BookedAt": "0001-01-01T00:00:00Z"
      }
    ]
  },
  "timestamp": "2024-06-30T05:36:55.085Z"
}

```

### Fetch All Parking Lots:

```shell
curl --location 'http://localhost:8080/api/v1/manager/parking-lots' \
--header 'Authorization: <manager-token>'
```

Response:
```json
{
    "status": "OK",
    "message": "Fetched Parking Lots",
    "success": true,
    "data": [
        {
            "ID": 1,
            "Name": "Lot A"
        }
    ],
    "timestamp": "2024-06-30T07:07:40.245Z"
}
```

### Fetch A Parking Lot:

```shell
curl --location 'http://localhost:8080/api/v1/manager/parking-lots/{lot_id}' \
--header 'Authorization: <manager-token>'
```

Response:
```json
{
    "status": "OK",
    "message": "Fetched Parking Lot",
    "success": true,
    "data": {
        "ID": 1,
        "Name": "Lot A",
        "Slots": [
            {
                "ID": 1,
                "UnderMaintenance": true,
                "Booked": false,
                "BookedAt": "0001-01-01T00:00:00Z"
            },
            {
                "ID": 2,
                "UnderMaintenance": false,
                "Booked": false,
                "BookedAt": "0001-01-01T00:00:00Z"
            }
        ]
    },
    "timestamp": "2024-06-30T07:08:03.971Z"
}
```

### Switch Maintenance Mode:
```shell
curl --location 'http://localhost:8080/api/v1/manager/parking-slot-status' \
--header 'Authorization: <manager-token>' \
--header 'Content-Type: application/json' \
--data '{
    "SlotID": 1,
    "MaintenanceMode": true
}'
```

Response:
```json
{
    "status": "OK",
    "message": "Maintenance Status Updated",
    "success": true,
    "data": {},
    "timestamp": "2024-06-30T06:16:58.596Z"
}
```

Slot that is under maintenance mode will not be assigned to users.

## User APIs:

You need to obtain user's access token by logging in as user using the login endpoint

### Fetch All Parking Lots:

```shell
curl --location 'http://localhost:8080/api/v1/user/parking-lots' \
--header 'Authorization: <user-token>'
```

Response:
```json
{
    "status": "OK",
    "message": "Fetched Parking Lots",
    "success": true,
    "data": [
        {
            "ID": 1,
            "Name": "Lot A"
        }
    ],
    "timestamp": "2024-06-30T05:38:11.147Z"
}
```

### Book a parking spot:

```shell
curl --location 'http://localhost:8080/api/v1/user/park' \
--header 'Authorization: <user-token>' \
--header 'Content-Type: application/json' \
--data '{
    "LotID": 1
}'
```

Response:
```json
{
    "status": "OK",
    "message": "Booked Parking Slot",
    "success": true,
    "data": {
        "LotID": 1,
        "SlotID": 2,
        "ParkedAt": "2024-06-30T00:59:31.253948Z"
    },
    "timestamp": "2024-06-30T06:59:31.255Z"
}
```

You will be assigned the nearest parking spot in the parking lot you select. You can book only one spot.

### Unparking from parking spot:
```shell
curl --location --request POST 'http://localhost:8080/api/v1/user/unpark' \
--header 'Authorization: <user-token>'
```

Response:
```json
{
    "status": "OK",
    "message": "Unbooked Parking Slot",
    "success": true,
    "data": {
        "Charge": "Rs. 10"
    },
    "timestamp": "2024-06-30T06:59:35.584Z"
}
```

User will be charged Rs. 10 for each one hour ceiling (ie 1-60 Minutes = 1 Hour).