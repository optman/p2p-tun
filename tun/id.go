package tun

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/protocol"
)

var (
	log = logging.Logger("p2p-tun")
)

const TcpProtocolID = protocol.ID("/tun/tcp")
const UdpProtocolID = protocol.ID("/tun/udp")
