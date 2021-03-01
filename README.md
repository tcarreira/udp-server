# UDP Server (and client)

Listens to UDP traffic and dumps it to output.

Example output:
``` 
##### packet-received: bytes=12 from=172.17.0.1:48241
##### binary/hex dump (first 32 bytes):
    0:     48 65 6c 6c 6f 2c 20 77 
    8:     6f 72 6c 64 00 00 00 00 
   16:     00 00 00 00 00 00 00 00 
   24:     00 00 00 00 00 00 00 00 
########################################
Hello, world
``` 



## Run with docker

Run server:

```
docker run --rm -p 1337:1337/udp tcarreira/udp-server 
```

Run client on another terminal (or send UDP with netcat or whatever):

```
docker run --rm --network=host tcarreira/udp-server --client Hello world
``` 

check other parameters (flags):
```
docker run --rm tcarreira/udp-server -h
``` 

## Run locally

Build (depends on Go):

```
make build
```


Run server:
```
./udp-server
``` 

Run client on another terminal (or send UDP with netcat or whatever):

```
./udp-server --client Hello world
``` 

check other parameters (flags):
```
./udp-server -h
``` 

---

Kudos to https://github.com/cirocosta/go-sample-udp
