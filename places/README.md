# Places Service

The places service stores places of interest by geolocation

Generated with

```
micro new places
```

## Usage

Places makes use of postgres. Set the config for the database

```
micro user config set places.database "postgresql://postgres@localhost:5432/locations?sslmode=disable"
```

Run the service

```
micro run .
```
