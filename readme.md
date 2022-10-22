# KSOC Backend Tech Test

Note: please do not post your solution to a public git repo
Create a clock application that will print the following values at the following intervals to stdout:

"tick" every second
"tock" every minute
"bong" every hour

Only one value should be printed in a given second, i.e. when printing "bong" on the hour, the "tick" and "tock" values should not be printed.

It should run for three hours and then exit.

A mechanism should exist for the user to alter any of the printed values while the program is running, i.e. after the clock has run for 10 minutes I should, without stopping the program, be able to change it so that it stops printing "tick" every second and starts printing "quack" instead.

## HOW TO USE

### Build

```
$ go build -o ticker_server cmd/main.go
```

### Test 

```
$ go test -v ./...
```

### Run

```
$ ./ticker_server -configpath=config/config.json
tick
tick
tock
tick
tick
tock
tick
tick
tock
bong
..
```

### List all current tickers

```
$ curl -H 'content-type: application/json' http://localhost:9876/signal
[{"id":3,"frequency":10,"message":"bong"},{"id":2,"frequency":3,"message":"tock"},{"id":1,"frequency":1,"message":"tick"}]
```

### Update ticker message


```
$ curl -H 'content-type: application/json' -X POST http://localhost:9876/signal --data '{"signal_id": 1, "message": "boom"}'
{"status": "ok"}
```