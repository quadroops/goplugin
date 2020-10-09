# Usages

- `go get -u -v` to install and update all dependencies
- Please build/compile our plugin's example inside `plugin` director (see their `README.md` there)
- Copy file `conf.toml.tmpl` into `config.toml`
- Re-configure the value inside `config.toml`
- Run this example using:

```
GOPLUGIN_DIR=<path/to/your/config.toml> go run -race *.go -msg=world3
```

Response should be like:

```
2020/09/18 00:18:42 Response ping: pong
2020/09/18 00:18:42 Response exec: world3
2020/09/18 00:18:42 Error wait: signal: killed
```

Note:

- Each you have make some changes for the plugin, you have to get `md5sum` from the compiled binary
- Put new md5 value inside your `config.toml`

To get md5 value from compiled binary file

```
md5sum plugin
```

This example, will give you an example how to use `goplugin.Supervisor`, when the process is running, it will print the plugin process ID (PID),
you can use this command :

```
sudo kill -9 <pid>
```

To demonstrate that our supervisor is working fine.  When you kill plugin process ID, our Supervisor will automatically detect and will auto restart
plugin process.