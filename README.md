# Gengar - Document storage and retrieval using CRYSTALS Kyber encryption

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

#### TODO:
- [ ] Documentation
- [x] Client and server implementations
- [ ] Unit tests? (maybe)
- [x] Figure out proper decryption