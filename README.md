# TCP Replay
Small go application to record and replay TCP streams

These are both stand-alone applications, build  / run them separately.

Only supports one way streaming from the server, dial in, the server starts talking. That's it.

Building:
```bash
go build record.go
go build play.go
```

Recording
```bash
./record --remote 10.10.10.10:1234 > tcp_dump
```

The TCP streams are dumped with the current unix nanosecond, and saved in hex strings - not efficient, but human readable.

Playback
```bash
./play --bind :1234 --source tcp_dump
# or
cat tcp_dump | ./play --bind :1234
```

Use together to intercept and log a stream
```bash
./record --remote 10.10.10.10:1234 | ./play --bind :1234
```


