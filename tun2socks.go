package tun2socks

import (
	"log"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"

	vcore "github.com/v2ray/v2ray-core"

	// The following are necessary as they register handlers in their init functions.
	// Required features. Can't remove unless there is replacements.
	_ "github.com/v2ray/v2ray-core/app/dispatcher"
	_ "github.com/v2ray/v2ray-core/app/proxyman/inbound"
	_ "github.com/v2ray/v2ray-core/app/proxyman/outbound"

	// Other optional features.
	_ "github.com/v2ray/v2ray-core/app/dns"
	_ "github.com/v2ray/v2ray-core/app/log"
	_ "github.com/v2ray/v2ray-core/app/policy"
	_ "github.com/v2ray/v2ray-core/app/router"
	_ "github.com/v2ray/v2ray-core/app/stats"

	// Inbound and outbound proxies.
	_ "github.com/v2ray/v2ray-core/proxy/blackhole"
	_ "github.com/v2ray/v2ray-core/proxy/dokodemo"
	_ "github.com/v2ray/v2ray-core/proxy/freedom"
	_ "github.com/v2ray/v2ray-core/proxy/http"
	_ "github.com/v2ray/v2ray-core/proxy/mtproto"
	_ "github.com/v2ray/v2ray-core/proxy/shadowsocks"
	_ "github.com/v2ray/v2ray-core/proxy/socks"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/inbound"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/outbound"

	// Transports
	_ "github.com/v2ray/v2ray-core/transport/internet/domainsocket"
	_ "github.com/v2ray/v2ray-core/transport/internet/http"
	_ "github.com/v2ray/v2ray-core/transport/internet/kcp"
	_ "github.com/v2ray/v2ray-core/transport/internet/quic"
	_ "github.com/v2ray/v2ray-core/transport/internet/tcp"
	_ "github.com/v2ray/v2ray-core/transport/internet/tls"
	_ "github.com/v2ray/v2ray-core/transport/internet/udp"
	_ "github.com/v2ray/v2ray-core/transport/internet/websocket"

	// Transport headers
	_ "github.com/v2ray/v2ray-core/transport/internet/headers/http"
	_ "github.com/v2ray/v2ray-core/transport/internet/headers/noop"
	_ "github.com/v2ray/v2ray-core/transport/internet/headers/srtp"
	_ "github.com/v2ray/v2ray-core/transport/internet/headers/tls"
	_ "github.com/v2ray/v2ray-core/transport/internet/headers/utp"
	_ "github.com/v2ray/v2ray-core/transport/internet/headers/wechat"
	_ "github.com/v2ray/v2ray-core/transport/internet/headers/wireguard"

	// The following line loads JSON internally
	_ "github.com/v2ray/v2ray-core/main/jsonem"
)

type PacketFlow interface {
	WritePacket(packet []byte)
}

var lwipStack core.LWIPStack
var isStopped = false
var v *vcore.Instance

func InputPacket(data []byte) bool {
	if lwipStack != nil {
		log.Println("lwipStack is  not nil")
	} else {
		log.Println("lwipStack is  NIL")
	}

	lwipStack.Write(data)

	return lwipStack != nil
}

func StartSocks(packetFlow PacketFlow, proxyHost string, proxyPort int) {
	if packetFlow != nil {
		debug.SetGCPercent(5)
		lwipStack = core.NewLWIPStack()
		core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, uint16(proxyPort)))
		core.RegisterUDPConnHandler(socks.NewUDPHandler(proxyHost, uint16(proxyPort), 2*time.Minute))
		core.RegisterOutputFn(func(data []byte) (int, error) {
			if !isStopped {
				packetFlow.WritePacket(data)
			}
			runtime.GC()
			debug.FreeOSMemory()
			return len(data), nil
		})

		isStopped = false
	}
}

func StopSocks() {
	log.Println("StopSocks")

	isStopped = true
	if lwipStack != nil {
		log.Println("StopSocks 2")

		lwipStack.Close()
		lwipStack = nil
	}
}

func StartV2RRayWithJsonData(configBytes []byte) bool {
	v, err := vcore.StartInstance("json", configBytes)
	if err != nil {
		log.Fatalf("start V instance failed: %v", err)
		return false
	}
	_ = v.Type()

	return true
}

func StopV2Ray() {
	log.Println("StopV2Ray")

	if v != nil {
		log.Println("StopV2Ray 2")

		v.Close()
		v = nil
	}
}
