# Traefik Auth Manager

Go app to manage credentials for traefik forward authentication

## App components

### Go-app

Go-app framework is used to develop the PWA frontend which can be installed.

### Go Echo framework

Go Echo framework is used to develop the backend API server.

## Debugging

When developing PWA a useful command to run to enable live reload

```bash
wgo -xdir tmp -file .go -file .css -file .js make run
```

## Build

Use following command to build the docker image

```sh
docker build . --network=host --tag apogee-dev/traefik-auth-manager:local
```
