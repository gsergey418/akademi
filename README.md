# Akademi

Akademi is a [Kademlia](https://en.wikipedia.org/wiki/Kademlia) implementation written in Go. The nodes of Akademi communicate with each other via UDP messages in [Protocol Buffers](https://protobuf.dev/). The project employs a modular architechture with loosely-coupled modules, each with their respective responsibilities. 

## Quick Start

To get started with Akademi clone the repo and run make, a binary file will appear in the root of the projects

1. Clone the project 
```
$ git clone https://github.com/gseryey418alt/akademi
```
2. Install development dependecies
```
$ sudo apt install go protoc
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
3. Build it
```
$ make
```
4. Start the daemon. 
```
$ ./akademi daemon
```

## Usage

Akademi uses 40-bytes base32-encoded IDs that identify nodes and data on the network. They look like this: ```QSWTYJD3HPOE54DWURBPICK7FAWWMVD3```. To store data on the network, you have to run the store command and provide up to 4KiB of data you want to store:

```
$ ./akademi store https://github.com/gsergey418alt/akademi
Data stored on the DHT successfully. KeyID: NB2HI4DTHIXS6Z3JORUHKYROMNXW2L3H.
```

The data will be stored on the DHT. To fetch it later use the get command and provide the key ID from the previous command:

```
$ ./akademi get NB2HI4DTHIXS6Z3JORUHKYROMNXW2L3H
Data retreived from the DHT successfully.
Content:
https://github.com/gsergey418alt/akademi
```

The data is usually stored on the DHT for up to an hour, after that it needs to be republished.

## Docker

There's an option to run an akademi network simulation in docker with 3 bootstrap nodes and 100 regular nodes. To start it run make swarm. (You need to have docker installed on your system)
```
$ make swarm
```

Use the docker exec command to interact with the containers:

```
$ docker exec akademi_1 akademi store "Hello, World!"
Data stored on the DHT successfully. KeyID: JBSWY3DPFQQFO33SNRSCDWRZUPXF422L.
```

You can also run a regular akademi node from a docker image.

```
$ make docker
$ docker run -p 3865:3865 -p 3855:3855 akademi:latest
```

## Running Tests

To run tests on the projects launch the daemon in standalone mode and run make tests.

1. Run the daemon
```
$ ./akademi daemon --no-bootstrap
```
2. Run the tests
```
$ make test
```