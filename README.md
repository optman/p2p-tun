# p2p-tun
A simple port forward and tun2socks tools build on libp2p with holepunch support.

## Usage
```
NAME:
   p2p-tun - port forward and tun2socks through libp2p

USAGE:
   p2p-tun [global options] command [command options] [arguments...]

COMMANDS:
   client   start client node
   server   start server node
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --id value           id seed (default: 0)
   --listen-port value  p2p listen port (default: 0)
   --secret value       authenticate user
   --debug              log debug (default: false)
   --help, -h           show help (default: false)
```

### port forward 

forward remote target:port to local port.

server
```
$ p2p-tun  --id 1234 --secret abcd server port --forward-addr 192.168.1.1:80
```

client
```
$ p2p-tun --secret abcd client --server-id 12D3KooWE3AwZFT9zEWDUxhya62hmvEbRxYBWaosn7Kiqw5wsu73  port --local-address :8888

$ curl localhost:8888
```


### tun2socks

tunnel tcp connection to remote built-in socks5 server (Full Access To Remote Network)


server
```
$p2p-tun --id 1234 --secret abcd server socks
```

client
```
$ sudo ip tuntap add mode tun tun0
$ sudo ip link set dev tun0 up
$ sudo ip route add 192.168.1.0/24 dev tun0

$ p2p-tun --secret abcd client --server-id 12D3KooWE3AwZFT9zEWDUxhya62hmvEbRxYBWaosn7Kiqw5wsu73  tun --tun-name tun0

$ curl 192.168.1.1:80
```


