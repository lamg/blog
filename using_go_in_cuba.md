# Using [Go][0] in Cuba

Some packages and tools related to [Go][0] programming cannot be downloaded from Cuba, due some restrictions Google applies. The savior for such cases turns out to be [Tor Browser][1], with a little addition.

These restrictions hit your face when you do something like:

```sh
git clone git@github.com:lamg/proxy.git
cd proxy/cmd/proxy
go install
```

since that project has dependencies hosted at the [official Go site][0], blocked from Cuba, which are golang.org/x/net and golang.org/x/text.

In this case, it is easy to download them manually from https://github.com/golang/net and https://github.com/golang/text respectively, and to put them under golang.org/x at the GOPATH. But is better to let the `go` command do that work for you, and it can be done if you replace the last command at the previous script by `https_proxy=http://localhost:8080 go install`.

This, of course, will not work right now since there's no such HTTP proxy tunneling requests made by the `go` command through another country, with access to [Go][0], at localhost:8080; which is what that line should do. For achieving that download and run [Tor Browser][1], and run `proxy -p socks5://localhost:9150 -f`. However this will not work right now, also you need to download the `proxy` command from https://github.com/lamg/proxy/releases, and that's all. Currently you will find there only the [Linux AMD64 binary](https://github.com/lamg/proxy/releases/download/v3.0.1/proxy_3.0.1_linux-amd64.zip), due to rational lazyness.

[0](https://golang.org)
[1](https://www.torproject.org/projects/torbrowser.html.en)
