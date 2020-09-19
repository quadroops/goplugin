# Usages

- Please build/compile our plugin's example inside `plugin` director (see their `README.md` there)
- Copy file `conf.toml.tmpl` into `config.toml`
- Re-configure the value inside `config.toml`
- Run this example using:

```
// default port using 8181
GOPLUGIN_DIR=<path/to/your/config.toml> go run -race *.go -msg="from host to grpc" -addr="localhost:8181"

// or you can modify using custom port
GOPLUGIN_DIR=<path/to/your/config.toml> go run -race *.go -msg="from host to grpc" -port=8383 -addr="localhost:8383"
```

Response should be like:

```
2020/09/19 13:01:00 Response ping: pong
2020/09/19 13:01:00 Response exec: from host to grpc
2020/09/19 13:01:00 Error wait: signal: killed
```

Note:

- Each you have made some changes for the plugin, you have to get `md5sum` from the compiled binary
- Put new md5 value inside your `config.toml`

Command to get md5 value from a file

```
md5sum plugin
```