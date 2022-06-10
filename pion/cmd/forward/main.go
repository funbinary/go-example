package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/pion/sdp/v3"

	"github.com/bin-work/go-example/pkg/pion"
	"github.com/bin-work/go-example/pkg/pion/srs"
	"github.com/pion/webrtc/v3"
)

func FindHostInCandidate(raw string) string {
	split := strings.Fields(raw)
	// Foundation not specified: not RFC 8445 compliant but seen in the wild
	if len(raw) != 0 && raw[0] == ' ' {
		split = append([]string{" "}, split...)
	}
	if len(split) < 8 {
		return ""
	}

	address := split[4]

	return address

}

func main() {
	receiverPC, err := pion.NewPeerConnection(webrtc.Configuration{}, true, true)
	if err != nil {
		panic(err)
	}
	videoTrackChan := make(chan *webrtc.TrackLocalStaticRTP)
	ctx := context.Background()

	//var audioSSRC webrtc.SSRC
	//var videoSSRC webrtc.SSRC
	//var videoPayloadType webrtc.PayloadType
	//for _, v := range forwardPC.GetReceivers() {
	//	if v.Track().Kind() == webrtc.RTPCodecTypeAudio {
	//		//audioSSRC = v.Track().SSRC()
	//		//fmt.Println(audioSSRC)
	//	}
	//	if v.Track().Kind() == webrtc.RTPCodecTypeVideo {
	//		//videoSSRC = v.Track().SSRC()
	//		videoPayloadType = v.Track().PayloadType()
	//	}
	//}
	// 设置方向
	receiverPC.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly})
	receiverPC.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo, webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly})

	offer, err := receiverPC.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	if err = receiverPC.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	host := "192.168.3.249"
	port := "1985"
	room := "62a15f2cb692e53204f1c694"
	display := "6846a36b9dc2"
	rtcUrl := "webrtc://" + host + ":" + port + "/" + room + "/" + display
	answer, err := srs.RtcRequest(ctx, "rtc/v1/play", rtcUrl, offer.SDP)
	if err != nil {
		panic(err)
	}

	parsed := &sdp.SessionDescription{}
	parsed.Unmarshal([]byte(answer))

	for _, m := range parsed.MediaDescriptions {
		for _, a := range m.Attributes {
			if a.IsICECandidate() {
				answer = strings.ReplaceAll(answer, FindHostInCandidate(a.Value), host)
			}
		}
	}

	if err = receiverPC.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: answer}); err != nil {
		panic(err)
	}

	//packets := make(chan *rtp.Packet, 60)

	receiverPC.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		//go func() {
		//	if track.Kind() == webrtc.RTPCodecTypeAudio {
		//		return
		//	}
		//
		//	for {
		//		select {
		//		case <-ctx.Done():
		//			return
		//		case <-time.After(time.Duration(10) * time.Second):
		//			_ = receiverPC.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{
		//				MediaSSRC: uint32(track.SSRC()),
		//			}})
		//		}
		//	}
		//}()
		codec := track.Codec()

		trackDesc := fmt.Sprintf("channels=%v", codec.Channels)
		if track.Kind() == webrtc.RTPCodecTypeVideo {
			trackDesc = fmt.Sprintf("fmtp=%v", codec.SDPFmtpLine)
		}
		if headers := receiver.GetParameters().HeaderExtensions; len(headers) > 0 {
			trackDesc = fmt.Sprintf("%v, header=%v", trackDesc, headers)
		}
		fmt.Printf("Got track %v, pt=%v, tbn=%v, %v\n",
			codec.MimeType, codec.PayloadType, codec.ClockRate, trackDesc)
		fmt.Printf("Ignore track %v pt=%v\n", codec.MimeType, codec.PayloadType)
		//var lastTimestamp uint32
		var videoTrack *webrtc.TrackLocalStaticRTP
		if track.Kind() == webrtc.RTPCodecTypeVideo {
			videoTrack, err = webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "pion")
			if err != nil {
				panic(err)
			}
			videoTrackChan <- videoTrack

		}
		rtpBuf := make([]byte, 1500)
		//rtpPacket := &rtp.Packet{}
		for {
			i, _, err := track.Read(rtpBuf)
			if err != nil {
				panic(err)
			}
			//if err = rtpPacket.Unmarshal(rtpBuf[:i]); err != nil {
			//	panic(err)
			//}
			//rtpPacket.PayloadType = uint8(videoPayloadType)
			//rtpPacket.PayloadType = 1
			//var n int
			//if n, err = rtpPacket.MarshalTo(rtpBuf); err != nil {
			//	panic(err)
			//}
			switch track.Kind() {
			case webrtc.RTPCodecTypeVideo:
				if _, err = videoTrack.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
					panic(err)
				}
			case webrtc.RTPCodecTypeAudio:
				//if _, err = audioTrack.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
				//	panic(err)
				//}
			}

			// Read RTP packets being sent to Pion
			//rtp, _, readErr := track.ReadRTP()
			//if readErr != nil {
			//	panic(readErr)
			//}
			//
			//oldTimestamp := rtp.Timestamp
			//if lastTimestamp == 0 {
			//	rtp.Timestamp = 0
			//} else {
			//	rtp.Timestamp -= lastTimestamp
			//}
			//lastTimestamp = oldTimestamp
			////var writeErr error
			//switch track.Kind() {
			//case webrtc.RTPCodecTypeVideo:
			//	//if writeErr = forwardPC.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}}); writeErr != nil {
			//	//	fmt.Println(writeErr)
			//	//}
			//	packets <- rtp
			//	rtp.SSRC = uint32(videoSSRC)
			//	fmt.Println(rtp.SSRC)
			//	//packet.SSRC = f
			//	// Write out the packet, ignoring closed pipe if nobody is listening
			//	//fmt.Println(packet.SSRC)
			//
			//	if err := videoTrack.WriteRTP(rtp); err != nil {
			//		if errors.Is(err, io.ErrClosedPipe) {
			//			// The peerConnection has been closed.
			//			return
			//		}
			//
			//		panic(err)
			//	}
			//case webrtc.RTPCodecTypeAudio:
			//	//if writeErr = forwardPC.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}}); writeErr != nil {
			//	//	fmt.Println(writeErr)
			//	//}
			//	//packets <- rtp
			//
			//	rtp.SSRC = uint32(audioSSRC)
			//	fmt.Println(rtp.SSRC)
			//	//packet.SSRC = f
			//	// Write out the packet, ignoring closed pipe if nobody is listening
			//	//fmt.Println(packet.SSRC)
			//
			//	if err := audioTrack.WriteRTP(rtp); err != nil {
			//		if errors.Is(err, io.ErrClosedPipe) {
			//			// The peerConnection has been closed.
			//			return
			//		}
			//
			//		panic(err)
			//	}
			//}

		}

	})

	//go func() {
	//	var currTimestamp uint32
	//	for i := uint16(0); ; i++ {
	//		packet := <-packets
	//		// Timestamp on the packet is really a diff, so add it to current
	//		currTimestamp += packet.Timestamp
	//		packet.Timestamp = currTimestamp
	//		// Keep an increasing sequence number
	//		packet.SequenceNumber = i
	//		packet.SSRC = uint32(videoSSRC)
	//		//packet.SSRC = f
	//		// Write out the packet, ignoring closed pipe if nobody is listening
	//		//fmt.Println(packet.SSRC)
	//
	//		if err := videoTrack.WriteRTP(packet); err != nil {
	//			if errors.Is(err, io.ErrClosedPipe) {
	//				// The peerConnection has been closed.
	//				return
	//			}
	//
	//			panic(err)
	//		}
	//	}
	//}()

	receiverPC.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		logrus.Infof("ICE state %v", state)

	})

	videoTrack := <-videoTrackChan

	forwardPC, err := pion.NewPeerConnection(webrtc.Configuration{}, true, true)
	if err != nil {
		panic(err)
	}
	//audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "pion")
	//if err != nil {
	//	panic(err)
	//}
	//
	//audioSender, err := forwardPC.AddTrack(audioTrack)
	//if err != nil {
	//	panic(err)
	//}

	forwardSender, err := forwardPC.AddTrack(videoTrack)
	if err != nil {
		panic(err)
	}
	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := forwardSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
			//if _, _, rtcpErr := audioSender.Read(rtcpBuf); rtcpErr != nil {
			//	return
			//}
		}
	}()
	//forwardPC.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionSendonly})
	//forwardPC.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo, webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionSendrecv})
	forwardOffer, err := forwardPC.CreateOffer(nil)
	if err != nil {
		panic(err)
	}
	//fmt.Println("offer:", forwardOffer.SDP)
	if err = forwardPC.SetLocalDescription(forwardOffer); err != nil {
		panic(err)
	}
	//ctx := context.Background()

	forwardhost := "192.168.3.249"
	forwardport := "1985"
	forwardroom := "live"
	forwarddisplay := "livet"
	forwardrtcUrl := "webrtc://" + forwardhost + ":" + forwardport + "/" + forwardroom + "/" + forwarddisplay
	forwardAnswer, err := srs.RtcRequest(ctx, "rtc/v1/play", forwardrtcUrl, forwardOffer.SDP)
	if err != nil {
		panic(err)
	}
	forwardparsed := &sdp.SessionDescription{}
	forwardparsed.Unmarshal([]byte(forwardAnswer))
	for _, m := range forwardparsed.MediaDescriptions {
		for _, a := range m.Attributes {
			if a.IsICECandidate() {
				forwardAnswer = strings.ReplaceAll(forwardAnswer, FindHostInCandidate(a.Value), forwardhost)
			}
		}
	}
	fmt.Println(forwardAnswer)
	if err = forwardPC.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: forwardAnswer}); err != nil {
		panic(err)
	}

	select {}
}
