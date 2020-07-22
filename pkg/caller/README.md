# pkg/caller

Package : `github.com/quadroops/goplugin/pkg/caller`

How should we handle a communication between application's host and our plugins ?

Each of plugins must implement a same communication channels.

There are two available actions:

1. Ping
2. Exec

**Ping**

Used to check plugin's process health.  It's simple just a request with payloand message of "ping" 
and should've got a response of "pong"

**Exec**

Payload:

- Command: A string of command action's name
- Content: An encoded byte data 

## Usages

A host and plugin should be agreed to define command's names and also their content payload (byte encoded) you can use 
protocol buffer or bson or messagepack to encode the data into binary, and you have to convert the binary data into string
using `hex` mechanism, example `hex` in golang:

```go

// payload is []byte
resp, err := client.Exec(context.Background(), &pbPlugin.ExecRequest{
    Command: cmdName,
    Payload: hex.EncodeToString(payload),
})

```

A plugin may have multiple actions, but only has a single entrypoint, an `exec` path.  When a plugin have multiple actions,
they can manage it by routing the action based on command's name.

For an example, a plugin of analyze log's data, it may have processes like:

- Store management, command name: `plugin.log.store`
- Filtering management, command name: `plugin.log.filter`
- Retention management, command name: `plugin.log.retention`

A command name should be used like routing key in RabbitMQ.

---

### Protocol: REST

There are two endpoints need to implement :

Endpoint: `/ping`

```
Method: GET
Request headers:
- X-Host: <hostname>

Response status header: 200 (OK)
Response body:
{
    "status": "success",
    "data": {
        "response": "pong"
    }
}
```

Endpoint: `/exec`

```
Method: POST
Request header:
- X-Host: <hostname>

Request body:
{
    "command": "<defined.command.name>",
    "payload": "<hexadecimal_encoded_bytes_in_string>" 
}

Response status header: 202 (accepted)
Response body:
{
    "status": "success",
    "data": {
        "response": "<hexadecimal_encoded_bytes_in_string>"
    }
}
```

An example of hex encoded byte (in golang):

```go
import "encoding/hex"

payload := "test"
b, _ := msgpack.Marshal(payload)
str := hex.EncodeToString(b) 
```

Unmarshalling

```go
// will return a bytes, you doesn't need to decode the hex
b, err := rest.Exec("rest.testing", []byte("test"))
assert.NoError(t, err)

var s string
err = msgpack.Unmarshal(b, &s)
assert.NoError(t, err)
```

Why hex ? We cannot take any assumptions here what kind of data type you will used, the most safety way to manage the request and response data is by sending
the data in the byte's form and encode the data (byte) using hex.  By using this method, you can manage you request and response payload based on your needs.
Please remember, even your request or response data is on "string" type, you have to convert it into byte and encode it using hex.

Data flows:

```
<your_data_type> -> convert into byte -> encode with hex
```

---

### Protocol: GRPC

For grpc's channel, a plugin must use and implement available proto's file

```protobuf
syntax = "proto3";
package plugin;

import "google/protobuf/empty.proto";

service Plugin {
    rpc Ping (google.protobuf.Empty) returns (PingResponse);
    rpc Exec (ExecRequest) returns (ExecResponse);
}

message Data {
    string Response = 1;
}

message PingResponse {
    string Status = 1;
    Data Data = 2;
}

message ExecRequest {
    string Command = 1;
    string Payload = 2;
}

message ExecResponse {
    string Status = 1;
    Data Data = 2;
}
```

---

We should using a same interface for our communication protocols :

```go
type Caller interface {
    Ping() (string, error)
    Exec(cmdName string, payload []byte) ([]byte, error)
}
```

Each of communication protocol, must implement `Caller` interface, usages:

```go
// a communication caller's helper used to construct a communicaton instance
// between rest or grpc
buildCaller := func(commType string, port int) caller.Caller {
    switch commType {
        case "rest":
            return driver.NewREST(commPort)
        case "grpc":
            return driver.NewGRPC(func() (pbPlugin.PluginClient, error) {
                conn, err := grpc.Dial(commPort)
                if err != nil {
                    return nil, err
                }

                return pbPlugin.NewPluginClient(conn), nil
            })
    }

    return &caller.Plugin{}
}

// get plugin's caller from container by passing plugin's name and the buildCaller helper
// so we doesnt need to kno what kind of communication's types used
plugin, err := container.Get("plugin_1", buildCaller)

// ping request
resp, err := plugin.Ping()

// exec request
resp, err := plugin.Exec("plugin.command.name", []byte("payload example"))
```