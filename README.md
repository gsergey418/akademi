# Akademi

Akademi is a [Kademlia](https://en.wikipedia.org/wiki/Kademlia) implementation written in Go. The nodes of Akademi communicate with each other via UDP messages in [Protocol Buffers](https://protobuf.dev/).

![Akademi](screenshot.png)

## Quick Start

To get started with Akademi, you can download the binary for your platform from the [releases](https://github.com/gsergey418/akademi/releases) page, or build from source by cloning the repository and running `make`, which will generate a binary file in the project's root directory.

1. Clone the Akademi repository:
```
$ git clone https://github.com/gsergey418/akademi
```
2. Install development dependencies:
```
$ sudo apt install go protoc
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
3. Build it:
```
$ make
```
4. Start the daemon and provide addresses of bootstrap nodes:
```
$ ./akademi daemon --bootstrap-nodes 133.146.89.223:3865
```

## Usage

Akademi uses 40-bytes base32-encoded IDs that identify nodes and data on the network. They look like this: `QSWTYJD3HPOE54DWURBPICK7FAWWMVD3`. To store data on the network, you have to run the `publish` command and provide up to 4KiB of data you want to store:

```
$ ./akademi publish https://github.com/gsergey418/akademi
Data published to the DHT successfully. KeyID: NB2HI4DTHIXS6Z3JORUHKYROMNXW2L3H.
```

The data will be stored on the DHT. To retrieve it later, use the `get` command and provide the key ID obtained from `publish`:

```
$ ./akademi get NB2HI4DTHIXS6Z3JORUHKYROMNXW2L3H
Data retreived from the DHT successfully.
Content:
https://github.com/gsergey418/akademi
```

The data is usually stored on the DHT for up to an hour, after that it needs to be republished.

## Docker

There's an option to run an Akademi network simulation in docker-compose with 2 bootstrap nodes and 6 regular nodes. To start it run `docker-compose up -d`. (You need to have docker and docker-compose installed on your system)
```
$ docker-compose up -d
```

Use the docker exec command to interact with the containers:

```
$ docker-compose exec akademi_1 akademi publish "Hello, World!"
Data published to the DHT successfully. KeyID: JBSWY3DPFQQFO33SNRSCDWRZUPXF422L.
```

You can also run a regular Akademi node from a docker image.

```
$ make docker
$ docker run --name akademi --rm -d -p 3865:3865/udp ghcr.io/gsergey418/akademi:latest
```

## Technical Details

When started, Akademi opens ports 3855 on the localhost and 3865 on all interfaces. The first one is used for communication with the CLI via RPC, and the second one for communicating with other nodes. Akademi has three main modules responsible for the core functionality: AkademiNode, Listener and Dispatcher. They communicate with each other via interfaces in the following way: AkademiNode holds the core Kademlia logic regarding routing and data storage, it is used by the Listener to react to incoming requests from other nodes, which is a one-way relationship from the Listener to AkademiNode. Meanwhile, AkademiNode holds an instance of the Dispatcher interface, that is used for dispatching UDP messages to other nodes, which is also a one-way relationship. These packages all depend on core types held in the "core" package and protocol buffer definitions in the package "pb". In regards to the core Kademlia logic, this software mostly adheres to the original whitepaper. Information on the network is stored as bytes with a maximum length of 4KiB. They are addressed via base32-encoded SHA1 hashes of their content, which is computed on write. Values in storage expire after one hour. This wasn't built with account for any persistence, so all the data and routing information is stored in memory.

```mermaid
graph TD;
    subgraph Daemon;
    UDPListener--Changes/Requests state-->AkademiNode;
    AkademiNodeRPCServer--Changes/Requests state-->AkademiNode;
    subgraph AkademiNode
    routingTable
    dataStore
    end
    AkademiNode--Dispatches requests to other nodes-->UDPDispatcher;
    end
    UDPDispatcher--Serialization-->OutgoingSocket["Outgoing UDP socket"]
    IncomingSocket["0.0.0.0:3865
    Incoming UDP socket"]--Deserialization-->UDPListener
    main["main()"]--./akademi daemon-->Daemon
    main--RPC Calls (publish, get, etc.)-->RPCSocket["127.0.0.1:3855
    RPC socket"]
    RPCSocket--Processes user's RPC calls-->AkademiNodeRPCServer
```

## Running Tests

To run tests on the project, run make test.

1. Run the tests, the daemon would be automatically started in the background on port 3865.
```
$ make test
```

## Contributing to Akademi

Thank you for considering contributing to Akademi! In order to submit your patch, please fork the repository and create a pull request on the main branch.
