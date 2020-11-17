# Location Service

A realtime GPS location tracking and search service

Generated with

```
micro new location
```

## Usage

Generate the proto code

```
make proto
```

Run the service

```
micro run .
```

### Test Service

```
go run examples/client.go
```

Output

```
Saved entity: id:"id123" type:"runner" location:<latitude:51.516509 longitude:0.124615 timestamp:1425757925 > 
Read entity: id:"id123" type:"runner" location:<latitude:51.516509 longitude:0.124615 timestamp:1425757925 > 
Search results: [id:"id123" type:"runner" location:<latitude:51.516509 longitude:0.124615 timestamp:1425757925 > ]
```
