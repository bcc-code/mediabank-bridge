# Mediabank bridge

This internal tool enables authenticated gRPC based endpoint for securely communicating with systems
like:

* Telestream Vantage Workflow (http://www.telestream.net/vantage/overview_workflow.htm)
* Vidispine (https://www.vidispine.com/)


## Building

### Requirements

You will need a Go 1.17+ compiler installed.
If you wish to make changes to the protobuf files you will also need `protoc` and the
Go plugin installed.

### Proto

```
make proto
```

### Binary

```
make release
```
The generated proto files are checked in, so you should be able do a build with no `protoc` installed
if no changes are made to the `.proto` files

### Config

The config is a `.json` file. The file is either located as specified by the `CONF_FILE`
environment variable, or (if not present) the following locations are checked in order:

* `~/.config/mediabank-bridge/config.json`
* `.config.json`
* `config.json`

An example is present in `src/config.example.json`.

Note: `DryRun` is not in use yet
