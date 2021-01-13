# Locations Service

This is the Locations service

Generated with

```
micro new locations
```

## Usage

Locations makes use of postgres. Set the config for the database

```
micro user config set locations.database "postgresql://postgres@localhost:5432/locations?sslmode=disable"
```

Run the service

```
micro run .
```
