package tun2socks

import (
	"log"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"

	vcore "v2ray.com/core"

	// The following are necessary as they register handlers in their init functions.
	// Required features. Can't remove unless there is replacements.
	_ "v2ray.com/core/app/dispatcher"
	_ "v2ray.com/core/app/proxyman/inbound"
	_ "v2ray.com/core/app/proxyman/outbound"

	// Other optional features.
	_ "v2ray.com/core/app/dns"
	_ "v2ray.com/core/app/log"
	_ "v2ray.com/core/app/policy"
	_ "v2ray.com/core/app/router"
	_ "v2ray.com/core/app/stats"

	// Inbound and outbound proxies.
	_ "v2ray.com/core/proxy/blackhole"
	_ "v2ray.com/core/proxy/dokodemo"
	_ "v2ray.com/core/proxy/freedom"
	_ "v2ray.com/core/proxy/http"
	_ "v2ray.com/core/proxy/mtproto"
	_ "v2ray.com/core/proxy/shadowsocks"
	_ "v2ray.com/core/proxy/socks"
	_ "v2ray.com/core/proxy/vmess/inbound"
	_ "v2ray.com/core/proxy/vmess/outbound"

	// Transports
	_ "v2ray.com/core/transport/internet/domainsocket"
	_ "v2ray.com/core/transport/internet/http"
	_ "v2ray.com/core/transport/internet/kcp"
	_ "v2ray.com/core/transport/internet/quic"
	_ "v2ray.com/core/transport/internet/tcp"
	_ "v2ray.com/core/transport/internet/tls"
	_ "v2ray.com/core/transport/internet/udp"
	_ "v2ray.com/core/transport/internet/websocket"

	// Transport headers
	_ "v2ray.com/core/transport/internet/headers/http"
	_ "v2ray.com/core/transport/internet/headers/noop"
	_ "v2ray.com/core/transport/internet/headers/srtp"
	_ "v2ray.com/core/transport/internet/headers/tls"
	_ "v2ray.com/core/transport/internet/headers/utp"
	_ "v2ray.com/core/transport/internet/headers/wechat"
	_ "v2ray.com/core/transport/internet/headers/wireguard"

	// The following line loads JSON internally
	_ "v2ray.com/core/main/jsonem"
)

type PacketFlow interface {
	WritePacket(packet []byte)
}

var lwipStack core.LWIPStack
var isStopped = false
var v *vcore.Instance

func InputPacket(data []byte) {
	if lwipStack != nil {
		log.Println("lwipStack is  not nil")
	} else {
		log.Println("lwipStack is  NIL")
	}

	lwipStack.Write(data)
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
