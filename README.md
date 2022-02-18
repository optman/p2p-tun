A simple port forward tools build on libp2p with holepunch support.


server
```
$p2p-tun  -f target-host:80 -i 1234
```

client
```
$p2p-tun -s 12D3KooWE3AwZFT9zEWDUxhya62hmvEbRxYBWaosn7Kiqw5wsu73 -l :88
$curl localhost:88
```
