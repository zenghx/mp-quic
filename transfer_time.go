package quic

import (
	"github.com/lucas-clemente/quic-go/congestion"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"time"
)

//传输时间估算
func time_calculation(dataToSend protocol.ByteCount, sender congestion.OliaSender, pth *path) time.Duration {
	if dataToSend == 0 {
		return 0
	}
	freeCwnd := sender.GetCongestionWindow() - pth.sentPacketHandler.GetBytesInFlight()
	rtt := sender.SmoothedRTT()
	if freeCwnd > 0 && dataToSend < freeCwnd {
		return rtt / 2
	}
	transferTime := rtt
	dataToSend -= freeCwnd
	sender.MaybeIncreaseCwnd(pth.sentPacketHandler.GetBytesInFlight())
	for dataToSend > sender.GetCongestionWindow() {
		sender.MaybeIncreaseCwnd(0)
		transferTime += sender.SmoothedRTT()
		dataToSend = dataToSend - sender.GetCongestionWindow()
	}
	transferTime += sender.SmoothedRTT() / 2
	return transferTime
}
