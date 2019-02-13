# kamailio-exporter

Simple prometheus exporter to query your kamailio instance for statistics and present it in prometheus format

Usage:

```
  $ go build
  $ kamailio-exporter --listen-address :8080 --instance http://kamailio:5060 --every 1m
```

All command line flags are read from ENV as well: `KAM_LISTEN_ADDRESS`, `KAM_INSTANCE` and `KAM_EVERY`.
