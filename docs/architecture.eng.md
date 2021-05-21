
`
Please excuse me, this text has been translated into English using the automatic translation service.
`
# Application architecture
When developing the application, the principles of DDD, SOLID, and the API paradigm of the Go-kit framework were taken into account

## Domains
The domain that implements the business logic of the application is wallet (internal/domain/wallet/).
The domain entities are Account and Payment (internal / domain/wallet/entity)

For storing and manipulating domain data, the domain repository (internal/domain/wallet/repository) is used,
in which interaction with the DBMS is implemented via drivers (internal/domain/wallet/repository/driver).

A repository for accessing drivers based on the dependency inversion principle. Thanks to this approach, you can easily change the
DBMS for storing data by writing a driver and specifying it when calling the object factory (AccountFactory, PaymentFactory)

## Services
In services (internal/services), the business logic of the API is implemented in accordance with the Go kit paradigm
Services are an intermediate layer between the domain business logic and the transport layer.

## Endpoints (internal/endpoints)
In the Go kit paradigm, these are virtual RPC methods. They make calls to the corresponding services.
Also, thanks to the provision of the interface type, endpoints allow you to completely abstract from the transport layer.

## Transport layer (internal / transport)
Provides interaction with users via HTTP requests, thus implementing the REST API.
In the Go Kit paradigm, the transport layer can implement other transport protocols (HTTP, gRPC, NATS, etc.).)

## Auxiliary Packages (pkg)
# # # Object caching (pkg / memcache)
Used by domain entity drivers for data caching. Allows you to save the number of queries to the database

### Gorutin Manager (pkg/subprocmgr)
A package for working with goroutines. Provides synchronization of goroutine completion at the end of the program.