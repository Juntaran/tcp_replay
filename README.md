# TCP Replay
Small go application to record and replay TCP streams

These are both stand-alone applications, build  / run them separately.

Building:
```bash
go build record.go
go build play.go
```


Recording
```bash
./record --remote 10.10.10.10:1234 > tcp_dump
```

Playback
```bash
./play --bind :1234 --source tcp_dump
# or
cat tcp_dump | ./play --bind :1234
```

The TCP streams are dumped with the current unix nanosecond, and saved in hex strings - not efficient, but human readable.

