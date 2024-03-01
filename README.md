# SolarEdge Reporting

This application receives stores SolarEdge statistics of a single site on a local timescaledb. 

## Environment variables

| Variable                    | Description                                         | Default |
|-----------------------------|-----------------------------------------------------| -------- |
| SEM_SOLAR_EDGE_API_KEY      | SolarEdge Api Key                                   | - |
| SEM_SOLAR_EDGE_SITE_ID      | SolarEdge site id to grab data for                  | - |
| SEM_SOLAR_EDGE_START_DATE   | Inception time of the SolarEdge site                | - |
| SEM_DB_NAME                 | DB name                                             | - |
| SEM_DB_HOST                 | DB host                                             | - |
| SEM_DB_PORT                 | DB port                                             | - |
| SEM_DB_USER                 | DB user                                             | - |
| SEM_DB_PASSWORD             | DB password                                         | - |
| SEM_AVOID_DOT_ENV           | Whether or not to load a .env configuration file    | `false` |
| SEM_HOST                    | The server host                                     | `localhost` |
| PORT                        | The server port                                     | - |
| SEM_WRITE_TIMEOUT           | The server write timeout                            | `10s` |
| SEM_READ_TIMEOUT            | The server read timeout                             | `10s` |
| SEM_IDLE_TIMEOUT            | The server idle timeout                             | `60s` |
| SEM_SERVER_SHUTDOWN_TIMEOUT | The timeout for gracefully shutting down the server | `5s` |

## Test
### Unit tests
```
$ go test ./...
```

### Unit and Integration tests
```
$ go test -tags=integration ./...
```
