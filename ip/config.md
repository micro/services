# IP Config

The ip service depends on the maxmind geolite2 dataset. You must configure its on disk location.

```
micro config set ip.city.database /tmp/GeoLite2-City.mmdb
micro config set ip.asn.database /tmp/GeoLite2-ASN.mmdb
```

In the event the config is not found it will attempt to read these two files from the local directory. 
If the config value is prefixed with `blob://` it will attempt to read it from the blob store and store on disk.
