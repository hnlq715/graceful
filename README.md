# graceful

Inspired by [overseer](https://github.com/jpillora/overseer) and [endless](https://github.com/fvbock/endless), with minimum codes and handy api to make http and grpc server graceful.

## Prerequisite

golang 1.8+

## Feature

- Graceful reload http servers, zero downtime on upgrade.
- Compatible with systemd, supervisor, etc.
- Drop-in placement for ```http.ListenAndServe```

## Example

``` golang
    type handler struct {
    }

    func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, port: %v, %q", r.Host, html.EscapeString(r.URL.Path))
    }

    func main(){
        graceful.ListenAndServe(":9222", &handler{})
    }
```

multi servers:

```golang
func listenMultiAddrs() {
	gs := grpc.NewService()
	pb.RegisterGreeterServer(gs.Server(), &server{})
	reflection.Register(gs.Server())

	hs := graceful.NewService()
	hs.Server().Handler = &handler{}

	server := graceful.NewServer()
	server.Register("0.0.0.0:9224", gs)
	server.Register("0.0.0.0:9225", hs)

	err := server.Run()
	fmt.Printf("error: %v\n", err)
}

```

```bash
grpcurl -plaintext -d '{"name":"sophos"}'  192.168.33.10:9224 helloworld.Greeter/SayHello
```

More example checkout example folder.

## Reload

```SIGHUP``` and ```SIGUSR1``` on master proccess are used as default to reload server. ```server.Reload()``` func works as well from your code.


## Drawbacks

```graceful``` starts a master process to keep pid unchaged for process managers(systemd, supervisor, etc.), and a worker proccess listen to actual addrs. That means ```graceful``` starts one more process. Fortunately, master proccess waits for signals and reload worker when neccessary, which is costless since reload is usually low-frequency action. 

## Default values

* ```StopTimeout```. Unfinished old connections will be drop in ```{StopTimeout}``` seconds, default 20s, after new server is up.

```golang
	server := graceful.NewServer(graceful.WithStopTimeout(time.Duration(4 * time.Hour)))
	server.Register(addr, handler)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
```
* ```Signals```. Default reload signals: ```syscall.SIGHUP, syscall.SIGUSR1``` and stop signals: ```syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT``` could be overwrited with:

```golang
	server := graceful.NewServer(graceful.WithStopSignals([]syscall.Signal{syscall.SIGKILL}), graceful.WithReloadSignals([]syscall.Signal{syscall.SIGHUP}))
	server.Register(addr, handler)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
```

## TODO

- ListenAndServeTLS
- Add alternative api: Run in single process without master-worker
