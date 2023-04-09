InnSecure --- Totally Secure Room Booking
=========================================

Innsecure provides a JSON API to manage hotel bookings.


## Prerequisites
- Install [Docker](https://docs.docker.com/get-docker/)
- Install [kubectl](https://kubernetes.io/docs/tasks/tools/)
- Install [Go 1.18](https://golang.org/) (you can use [gvm](https://github.com/moovweb/gvm) for using a particular version)

## Getting Started
- `make test` will run unit and end-to-end tests
- `make build` builds the app and token generator (see 'Generating JWTs, below')
- `make start-kind` creates a local k8s cluster + builds and spins up the service on port 8080

Then run
```
kubectl port-forward service/innsecure 8080:8080 -n innsecure
```

- `make stop-kind` tears down the local k8s cluster

`k8s/` folder contains all the necessary `k8s` manifests to run this application in a production environment.

## Generating JWTs
Authorisation is handled using [JSON Web Tokens](jwt.io). The default service started via Kind uses a simple signing key of `SigningString`. Acquisition of a JWT is outside the scope of this task, but a tool is included to help with generating tokens for testing. `make build` produces a binary at `./bin/token`, which will generate a working JWT for a standard user in hotel `123`. The signing key, hotel ID, and admin status can be overridden using the `-key`, `-hotel`, and `-admin` flags, respectively.

Using the generated token, set the `Authorization` header of your HTTP requests to `Bearer $YOUR_TOKEN_HERE`, e.g.

```
curl --location --request GET 'http://localhost:8080/hotels/123/bookings' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwibmFtZSI6IkguQS4gS2VyciIsIm9yZyI6MTIzLCJzdWIiOiI4NzkwYzUxNC03M2I2LTQwMGYtOGYyOC1hY2M3NGQzNDJhMjIifQ.jEVthaAkJZ2mQ0jNsXH1oSMGcYX1-_mcXcdqkLYZz8Q'
```

## Endpoints
### Create a Booking
`POST /hotels/:hotelID/bookings`

Using a body like this one:
```
{

"hotel_id":123,
"type": "Booking",
"version": 0,
"arrive": " 2021-08-13",
"leave": " 2021-08-14",
"name": "Geoff Capes"
}
```

Bookings can only be created by an admin user. See _Generating JWTs_, above for info on how to create an admin token.
Bookings can only be created in the hotel of the user.

### Fetch a Booking
`GET /hotels/:hotelID/bookings/:bookingID`

Bookings can only be fetched from a user's own hotel.

### List Bookings
`GET /hotels/:hotelID/bookings`

A user can only list bookings from their own hotel.