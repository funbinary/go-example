package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pion/rtcp"

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
	pc, err := pion.NewPeerConnection(webrtc.Configuration{}, true, true)
	if err != nil {
		panic(err)
	}

	// 设置方向
	pc.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly})
	pc.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo, webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly})

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	if err = pc.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	ctx := context.Background()
	host := "192.168.3.249"
	port := "1985"
	room := "62a151f438feb98851a01bbb"
	display := "48d12155ab78"
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

	if err = pc.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: answer}); err != nil {
		panic(err)
	}

	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		go func() {
			if track.Kind() == webrtc.RTPCodecTypeAudio {
				return
			}

			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Duration(10) * time.Second):
					_ = pc.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{
						MediaSSRC: uint32(track.SSRC()),
					}})
				}
			}
		}()
		codec := track.Codec()

		trackDesc := fmt.Sprintf("channels=%v", codec.Channels)
		if track.Kind() == webrtc.RTPCodecTypeVideo {
			trackDesc = fmt.Sprintf("fmtp=%v", codec.SDPFmtpLine)
		}
		if headers := receiver.GetParameters().HeaderExtensions; len(headers) > 0 {
			trackDesc = fmt.Sprintf("%v, header=%v", trackDesc, headers)
		}
		fmt.Println("Got track %v, pt=%v, tbn=%v, %v",
			codec.MimeType, codec.PayloadType, codec.ClockRate, trackDesc)
		fmt.Println("Ignore track %v pt=%v", codec.MimeType, codec.PayloadType)

	})

	pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		logrus.Infof("ICE state %v", state)

	})

	select {}
}
