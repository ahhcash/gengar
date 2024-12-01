# Gengar &nbsp; ![GitHub Actions](https://github.com/ahhcash/gengar/actions/workflows/build.yml/badge.svg)

### Document storage and retrieval using Kyber post quantum encryption

> #### Information Security and Privacy - Assignment 3.2

### Project Structure
```
.
├── Makefile
├── README.md
├── cmd
│   ├── client
│   │   ├── client.go
│   │   └── main.go
│   └── server
│       ├── main.go
│       ├── server.go
│       └── store.go
├── go.mod
├── go.sum
├── internal
│   ├── crypto
│   │   ├── encryptor.go
│   │   └── kyber.go
│   └── types
│       └── document.go
└── proto
    └── document_service.proto

8 directories, 13 files
```
### Overview.
This project is a proof of concept for Kyber 768 encryption in a (simulated) cloud storage system that performs document storage and retrieval. The project is built using [OQS' Go port of Kyber](https://github.com/open-quantum-safe/liboqs-go) library and the [gRPC](https://grpc.io/) framework.

##### What is Kyber?

Kyber is a post-quantum cryptographic algorithm that is designed to be secure against quantum computers. The OQS implementation of Kyber performs key encapsulation. This means that the client and server can share a secret key, but the key itself encrypted using the public key of the client. We use Kyber768 as the key encapsulation mechanism in conjunction with AES-GCM for actual document encryption.

There are two main components to this project:
#### 1. Client
The client is a simple command line REPL (Read-Eval-Print Loop) that connects to the server and can perform the following operations:
- Upload a document
- Download a document
- List all documents
- View an encrypted document
- Exit

Multiple client instances can connect to server, **but** they will have their own encryption keys. This means that if a client uploads a document to the server, only that client will be able to download it. However, other clients will be able to only view the document's encrypted content.

#### 2. Server
The server is a gRPC service that implements the same operations:
- Upload a document
- Download a document
- List all documents
- View an encrypted document
- Exit

### Instructions.

#### 1. Setting up `liboqs`
 The [liboqs-go README](https://github.com/open-quantum-safe/liboqs-go#installation) has instructions on how to install `liboqs` on your system. Make sure the `.config` directory in `liboqs-go` has the `LIBOQS_INCLUDE_DIR` and `LIBOQS_LIB_DIR` variables set to the location of `liboqs` headers and libraries on your system, since they're dynamically linked (`liboqs-go` internally uses CGO to link to `liboqs`).

#### 2. Generating protobuf stubs
Since we're using gRPC, we need to generate protobuf stubs for our service. A `Makefile` with a `proto` target is provided to do exactly that.
Simply run 
```bash
make proto
```
to generate the stubs. They'll be located in `proto/generated/proto`.

#### 3. Starting the server
The `Makefile` is super handy for this. Simply run
```bash
make server 
```
and the server will start listening on port `50051`.

#### 4. Spinning up clients
Use the `client` target in the `Makefile` to start a client.
```bash
make client
```
This will start a REPL client that connects to the server on port `50051`. Typing `help` will give a list of valid REPL commands.
Since there's no key persistence (everything's in memory), exiting a client will mean the client (if spun up again) will no longer be able to download documents uploaded previously.

### Okay, and? 
A cloud storage system is simulated by the `server` and `clients` are able to store and fetch documents stored in this cloud. With this, we successfully delivered a proof of concept for Kyber 768 encryption in a cloud storage system.

### What's next?
Potential improvements/features:
- add a key persistence mechanism
- the ability to share documents with other clients
- add a web interface
- attach a persistent storage layer (rather than in-memory)

### References
#### here are the main ones:
- [OQS](https://openquantumsafe.org/)
- [Kyber](https://pq-crystals.org/kyber/)
- [gRPC](https://grpc.io/)
- [liboqs-go](https://github.com/open-quantum-safe/liboqs-go)
- [liboqs](https://github.com/open-quantum-safe/liboqs)
- [protobuf](https://developers.google.com/protocol-buffers)

#### and here are some that helped along the way:
- [protobuf with GitHub Actions](https://www.andreamedda.com/posts/go-buf-github-actions/)
- [add to PATH on GitHub Actions](https://stackoverflow.com/questions/60169752/how-to-update-the-path-in-a-github-action-workflow-file-for-a-windows-latest-hos)
- [add to ENV on GitHub Actions](https://stackoverflow.com/questions/57968497/how-do-i-set-an-env-var-with-a-bash-expression-in-github-actions)
- [rebuild dynamic libraries on linux](https://unix.stackexchange.com/questions/694156/how-do-you-reload-so-files-dynamic-libraries-in-linux)