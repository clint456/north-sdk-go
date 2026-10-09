# north-sdk-go Example

A minimal, runnable example of a northbound service built on `north-sdk-go`.

## Layout

- `driver/simplenorth.go` – a `ProtocolDriver` stub that seeds a demo north
  profile/channel and logs everything it would send over the wire.
- `cmd/north-simple/main.go` – wires the driver into `startup.Bootstrap`.
- `cmd/north-simple/res/configuration.yaml` – service configuration.

## Run

Requires a running EdgeX **core-metadata** (default `localhost:59881`), because
metadata persistence is mandatory (see below).

```bash
cd example/cmd/north-simple
go run . -cd res
```

You should see the service boot, the demo `north-mqtt` channel start, and the
command handler register. Press `Ctrl+C` to stop gracefully.

## Using the SDK in your own driver

```go
type MyDriver struct{ sdk interfaces.NorthServiceSDK }

func (d *MyDriver) Initialize(sdk interfaces.NorthServiceSDK) error {
    d.sdk = sdk
    return nil
}

func (d *MyDriver) Start(northName string, cfg models.ProtocolProperties) error {
    // establish the outbound connection for this channel
    return d.sdk.RegisterCommandHandler(northName, d.handleCommand)
}

func (d *MyDriver) HandleSend(northName string, event models.Event) error {
    // encode event.Readings and push to the upstream platform
    return nil
}

// ... implement the remaining ProtocolDriver methods ...
```

Then:

```go
func main() {
    startup.Bootstrap("my-north-service", "0.0.1", &MyDriver{})
}
```

### Sending data northbound

- Synchronous: `sdk.PublishEvent(event)` where `event.DeviceName` is the target
  north channel name.
- Asynchronous: write `*sdkmodels.UpstreamEvent` values to
  `sdk.AsyncEventsChannel()`; the SDK dispatches them to the driver.

### Receiving commands

Register a handler with `sdk.RegisterCommandHandler(northName, handler)`. An
inbound transport (message bus / REST) delivers commands to the SDK by calling
`sdk.DispatchCommand(northName, cmd)`.

### Metadata persistence

Like `device-sdk-go`, metadata is **mandatory**: north channels and profiles are
persisted in EdgeX **core-metadata**, reached through the bootstrap `Clients`
configuration. `configuration.yaml` must declare a `core-metadata` client, and
core-metadata must be running when the service starts:

```yaml
Clients:
  core-metadata:
    Protocol: http
    Host: localhost
    Port: 59881
    UseMessageBus: false
```

> Note: core-metadata must be reachable at startup. As in `device-sdk-go`, the
> service refuses to start if the `core-metadata` client is missing or the
> metadata service is unavailable.
