Store and search for points of interest

# Places Service

The places API stores points of interest and enables you to search for places nearby or last visited.


## Usage

Places makes use of postgres. Set the config for the database

```
micro user config set places.database "postgresql://postgres@localhost:5432/locations?sslmode=disable"
```

Run the service

```
micro run .
```
