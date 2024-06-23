# remote deployment svc
This is orchestrator that reliably deploys application using fly machine apis. Historically, flyctl has always been our orchestrator, but we've decided to move some of the orchestrating into a remote service with a backing datastore(volume).

# development

Install tooling for codegen 
```
go install github.com/bufbuild/buf/cmd/buf@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

To add a new migration
```
make migration name=<name_of_migration>
```

To run locally
```
make run
```