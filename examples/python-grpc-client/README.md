# python-grpc-client example

Minimal Python client showing how behavox (or any non-Go consumer)
talks to `cmd/demanglegrpc`.

## Setup

```
python3 -m venv .venv
. .venv/bin/activate
pip install grpcio grpcio-tools

# Generate Python stubs from the proto.
python3 -m grpc_tools.protoc \
    -I ../../cmd/demanglegrpc/proto \
    --python_out=. \
    --grpc_python_out=. \
    ../../cmd/demanglegrpc/proto/demangle.proto
```

## Server (on the target host)

```
cd ../..
go run ./cmd/demanglegrpc --listen 127.0.0.1:50061
```

## Client

```python
# client.py
import grpc
import demangle_pb2
import demangle_pb2_grpc

channel = grpc.insecure_channel("127.0.0.1:50061")
stub = demangle_pb2_grpc.DemangleStub(channel)

resp = stub.Demangle(demangle_pb2.Request(input="_ZN4llvm5Value4dumpEv"))
print(resp.scheme, "→", resp.output)

# Batch streaming.
def requests():
    for i, s in enumerate([
        "_ZN4llvm5Value4dumpEv",
        "$s4main3fooyyF",
        "Java_com_example_Foo_bar",
    ]):
        yield demangle_pb2.Request(id=i, input=s)

for resp in stub.DemangleStream(requests()):
    print(resp.id, resp.scheme, "→", resp.output)
```

## Upload a ProGuard map + resolve

```python
with open("app.map", "rb") as f:
    blob = f.read()
up = stub.UploadContext(demangle_pb2.UploadContextRequest(
    kind="proguard_map", blob=blob))
print("uploaded:", up.sha256)

resp = stub.Demangle(demangle_pb2.Request(
    input="a.b",
    scheme="proguard-map",
    options=demangle_pb2.Options(context_sha256=up.sha256),
))
print(resp.output)  # com.example.Foo.bar
```

## Integrating with skynet behavox

Behavox on saskia calls the demangle library through skynet's
GraphQL (see `docs/architecture.md`) — that's the primary path and
doesn't need this gRPC client. This example is the alternative for
when a non-Go non-skynet consumer appears (Stage 6.5 trigger).
