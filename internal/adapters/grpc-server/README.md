# gRPC Server Adapter

This adapter exposes the application's use cases (`ports.ApiPort`) over gRPC, the
same way `internal/adapters/chi-server` exposes them over REST. It does not
contain any business logic — every RPC handler just validates input, delegates
to `ports.ApiPort`, and maps the result/error to a protobuf message / gRPC
status code.

## Layout

```
proto/users/v1/users.proto                              # RPC contracts (source of truth)
internal/adapters/grpc-server/
  gen/users/v1/                                          # generated code — never hand-edit
    users.pb.go                                          # message types
    users_grpc.pb.go                                     # service/client interfaces
  handlers/
    mappers.go                                            # shared: application error -> grpc status
    users/
      handler.go                                          # UsersService RPC implementations
      mappers.go                                           # domain <-> proto message mapping
  interceptors/
    logging.go                                             # request logging (applies to every service)
  server.go                                                # grpc.Server wiring, service registration
cmd/cli/cobra/grpc-server/                                 # `grpc-server` cobra command
```

Requests flow: `gRPC client -> generated *ServiceServer interface -> handler
(validation + error mapping) -> ports.ApiPort -> application/api -> repository`.

## Prerequisites

Code generation needs the `protoc` compiler itself, which isn't Go-installable.
Install it once via your package manager or a
[release binary](https://github.com/protocolbuffers/protobuf/releases)
(`protoc --version` should work afterwards). The `protoc-gen-go` /
`protoc-gen-go-grpc` plugins are installed automatically by `make
generate/proto`.

## Adding a new RPC to an existing service (e.g. `UsersService.GetUser`)

1. **Extend the `.proto` file** (`proto/users/v1/users.proto`):

   ```protobuf
   message GetUserRequest {
     int64 id = 1;
   }

   message GetUserResponse {
     User user = 1;
   }

   service UsersService {
     rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
     rpc GetUser(GetUserRequest) returns (GetUserResponse);
   }
   ```

2. **Regenerate the Go code**:

   ```bash
   make generate/proto
   ```

   This overwrites everything in `internal/adapters/grpc-server/gen/`. It also
   regenerates `UsersServiceServer`, adding `GetUser(...)` to the interface.
   Because the handler embeds `UnimplementedUsersServiceServer`, this compiles
   fine even before you implement it — it just responds with a
   `codes.Unimplemented` error at runtime until you add the method below, so
   don't skip step 3.

3. **Implement the method** on the existing handler
   (`handlers/users/handler.go`), reusing `ports.ApiPort` and the shared error
   mapper:

   ```go
   func (s *Server) GetUser(_ context.Context, req *usersv1.GetUserRequest) (*usersv1.GetUserResponse, error) {
       user, err := s.api.GetUserByID(strconv.FormatInt(req.GetId(), 10))
       if err != nil {
           return nil, handlers.MapErrorIntoStatusError(err)
       }

       return &usersv1.GetUserResponse{User: UserDomainToProto(user)}, nil
   }
   ```

   No changes are needed in `server.go` — the method becomes available the
   moment it's registered on the (already-registered) `UsersService`.

## Adding a brand new service (e.g. `OrdersService`)

1. **Create the proto package**: `proto/orders/v1/orders.proto`, with its own
   `go_package` option pointing at a new generated package, e.g.:

   ```protobuf
   syntax = "proto3";

   package orders.v1;

   option go_package = "github.com/diogorodriguesc/boilerplate-go/internal/adapters/grpc-server/gen/orders/v1;ordersv1";

   message Order {
     int64 id = 1;
   }

   message CreateOrderRequest {
     int64 user_id = 1;
   }

   message CreateOrderResponse {
     Order order = 1;
   }

   service OrdersService {
     rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
   }
   ```

2. **Regenerate**: `make generate/proto` (it globs every `.proto` under
   `proto/`, so new files are picked up automatically — nothing to configure).
   This produces `internal/adapters/grpc-server/gen/orders/v1/`.

3. **If the use case doesn't exist yet, add it to the application layer
   first**, the same way you would for a new REST endpoint:
   - `internal/application/domain` — domain type, if needed.
   - `internal/application/ports/ports.go` — add the method to `ApiPort`
     (and to `UserRepository`/a new repository interface, if it needs
     persistence).
   - `internal/application/api/api.go` — implement it.
   - The repository adapter under `internal/adapters/postgres-service-repository/...`.

   The gRPC handler should never talk to a repository directly — always go
   through `ports.ApiPort`, exactly like the REST handlers do.

4. **Add the handler package**:
   `internal/adapters/grpc-server/handlers/orders/handler.go`:

   ```go
   package orders

   import (
       "context"

       ordersv1 "github.com/diogorodriguesc/boilerplate-go/internal/adapters/grpc-server/gen/orders/v1"
       "github.com/diogorodriguesc/boilerplate-go/internal/adapters/grpc-server/handlers"
       "github.com/diogorodriguesc/boilerplate-go/internal/application/ports"
   )

   type Server struct {
       ordersv1.UnimplementedOrdersServiceServer
       api ports.ApiPort
   }

   func NewServer(api ports.ApiPort) *Server {
       return &Server{api: api}
   }

   func (s *Server) CreateOrder(_ context.Context, req *ordersv1.CreateOrderRequest) (*ordersv1.CreateOrderResponse, error) {
       // validate req, call s.api.CreateOrder(...), map the result/error
   }
   ```

   Add a `mappers.go` next to it for domain <-> proto conversions, following
   `handlers/users/mappers.go`.

5. **Register the service** in `server.go`'s `registerServices()`:

   ```go
   ordersv1.RegisterOrdersServiceServer(server, ordersHandler.NewServer(s.api))
   ```

   Reflection and the logging interceptor are registered once on the shared
   `grpc.Server`, so every service you add is automatically covered by both —
   no per-service wiring needed.

6. **Test it** (reflection is on, so no `.proto` needed client-side):

   ```bash
   grpcurl -plaintext localhost:9090 list
   grpcurl -plaintext -d '{"userId": 1}' localhost:9090 orders.v1.OrdersService/CreateOrder
   ```

## Conventions to follow

- **Validation**: use `validator.Var(...)` per field in the handler (proto
  messages don't carry struct tags), mirroring the equivalent REST request's
  `validate` tags — see `handlers/users/handler.go`.
- **Errors**: never return a raw `error` from a handler. Map it with
  `handlers.MapErrorIntoStatusError`, extending that switch if a new
  `applicationerrors` sentinel needs a dedicated gRPC status code.
- **No business logic here**: if a handler needs more than
  validate -> call `ports.ApiPort` -> map response, that logic belongs in
  `internal/application`, not in this adapter.
- **Don't hand-edit anything under `gen/`** — it's regenerated wholesale by
  `make generate/proto` and any manual change will be silently lost.
