# Gengar - Document storage and retrieval using CRYSTALS Kyber encryption

> #### Information Security and Privacy - Assignment 3.2

### Project Structure

```
├── Makefile
├── README.md
├── cmd
│   ├── client
│   │   └── main.go
│   └── server
│       ├── main.go
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
    ├── document_service.proto
    └── generated
        └── proto
            ├── document_service.pb.go
            └── document_service_grpc.pb.go

10 directories, 13 files
```

#### TODO:
- [ ] Documentation
- [ ] Client and server implementations
- [ ] Unit tests? (maybe)
- [ ] Figure out proper decryption