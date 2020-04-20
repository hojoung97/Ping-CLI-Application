# Ping-CLI-Application
Ping CLI Application in Go

## Build
Execute the line from project root directory to create an executable called ping

```
go build -o ping main.go
```

## Run
```
./ping [-t timeout] [-i interval] [-c count] [-s packetSize] [-root root] [-h help] <hostname or IP address>
```

Square brackets([]) are flags and angle brackets(<>) are arguments to the executable

** if running with 'sudo' keyword, MUST set the -root flag when executing

## Notes
- Support for both IPv4 and IPv6 available
- Make sure to check if local machine has IPv6 configuration before running with IPv6 mode
