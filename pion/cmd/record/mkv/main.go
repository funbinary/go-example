package mkv

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/bin-work/go-example/pkg/pion"
	"github.com/bin-work/go-example/pkg/pion/srs"

	"github.com/bin-work/go-example/pkg/pion/record"

	"github.com/pion/interceptor"
	"github.com/pion/rtcp"
	"github.com/pion/sdp/v3"
	"github.com/pion/webrtc/v3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Stream struct {
	Host          string
	Port          string
	Room          string
	Display       string
	rtcUrl        string
	SavePath      string
	IsOgg         bool
	ctx           context.Context
	cancel        context.CancelFunc
	pc            *webrtc.PeerConnection
	hasVideoTrack chan struct{}
	hasAudioTrack chan struct{}
	saver         *record.MkvSaver
}

func (self *Stream) Stop() {
	self.pc.Close()
	self.cancel()
	self.saver.Close()
}

func NewStream(host, port, room, display, savePath string) (*Stream, error) {
	var err error
	stream := &Stream{
		Host:          host,
		Port:          port,
		Room:          room,
		Display:       display,
		rtcUrl:        "webrtc://" + host + "/" + room + "/" + display,
		SavePath:      savePath,
		hasAudioTrack: make(chan struct{}, 1),
		hasVideoTrack: make(chan struct{}, 1),
	}
	if port != "" {
		stream.rtcUrl = "webrtc://" + host + ":" + port + "/" + room + "/" + display
	}
	logrus.Info("拉取", stream.rtcUrl+"rtc流")
	stream.ctx, stream.cancel = context.WithCancel(context.Background())
	stream.saver = record.NewMkvSaver(savePath)

	//创建PeerConncetion
	stream.pc, err = newPeerConnection(webrtc.Configuration{})
	//stream.pc.

	if err != nil {
		return nil, errors.Wrapf(err, "创建PeerConnection失败")
	}

	//设置方向
	stream.pc.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RTPTransceiverInit{
		Direction: webrtc.RTPTransceiverDirectionRecvonly,
	})
	stream.pc.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo, webrtc.RTPTransceiverInit{
		Direction: webrtc.RTPTransceiverDirectionRecvonly,
	})

	//创建offer
	offer, err := stream.pc.CreateOffer(nil)
	if err != nil {
		return nil, errors.Wrap(err, "创建Local offer失败")
	}

	// 设置本地sdp
	if err = stream.pc.SetLocalDescription(offer); err != nil {
		return nil, errors.Wrap(err, "设置Local SDP失败")
	}

	// 设置远端SDP
	answer, err := srs.RtcRequest(stream.ctx, "/rtc/v1/play", stream.rtcUrl, offer.SDP)
	//answer = strings.ReplaceAll(answer, "dev.beyondinfo.com.cn", "192.168.3.250")
	//answer = strings.ReplaceAll(answer, "dev.beyondinfo.com.cn", "117.28.132.153")

	if err != nil {
		return nil, errors.Wrap(err, "SDP协商失败")
	}

	// DNS解析

	parsed := &sdp.SessionDescription{}
	parsed.Unmarshal([]byte(answer))

	for _, m := range parsed.MediaDescriptions {
		for _, a := range m.Attributes {
			if a.IsICECandidate() {
				answer = strings.ReplaceAll(answer, pion.FindHostInCandidate(a.Value), host)
			}
		}
	}
	if err = stream.pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer, SDP: answer,
	}); err != nil {
		return nil, errors.Wrap(err, "设置Remote SDP失败")
	}

	stream.pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		if track.Kind() == webrtc.RTPCodecTypeVideo {
			go func() {
				ticker := time.NewTicker(time.Second * 3)
				select {
				case <-stream.ctx.Done():
					ticker.Stop()
					break
				case <-ticker.C:
					rtcpSendErr := stream.pc.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}})
					if rtcpSendErr != nil {
						logrus.Error(rtcpSendErr)
					}
				}
			}()
		}
		err = stream.onTrack(track, receiver)
		if err != nil {
			codec := track.Codec()
			logrus.Errorf("Handle  track %v, pt=%v\nerr %v", codec.MimeType, codec.PayloadType, err)
			stream.cancel()

		}
	})
	go func() {
		for {
			select {
			case <-stream.hasVideoTrack:
				runtime.Goexit()
			case <-time.After(6 * time.Second):
				stream.saver.InitWriter(0, 0)
				stream.IsOgg = true
				runtime.Goexit()
			}
		}

	}()
	stream.pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		logrus.Infof("ICE state %v", state)

		if state == webrtc.ICEConnectionStateFailed || state == webrtc.ICEConnectionStateClosed {
			if stream.ctx.Err() != nil {
				return
			}

			logrus.Warnf("Close for ICE state %v", state)
			stream.cancel()
		}
	})

	return stream, nil
}

func (self *Stream) onTrack(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) error {
	// Send a PLI on an interval so that the publisher is pushing a keyframe
	codec := track.Codec()

	trackDesc := fmt.Sprintf("channels=%v", codec.Channels)
	if track.Kind() == webrtc.RTPCodecTypeVideo {
		trackDesc = fmt.Sprintf("fmtp=%v", codec.SDPFmtpLine)
	}
	logrus.Infof("Got track %v, pt=%v tbn=%v, %v", codec.MimeType, codec.PayloadType, codec.ClockRate, trackDesc)

	var err error
	switch track.Kind() {
	case webrtc.RTPCodecTypeAudio:
		self.hasAudioTrack <- struct{}{}
	case webrtc.RTPCodecTypeVideo:
		self.hasVideoTrack <- struct{}{}
	}
	for self.ctx.Err() == nil {
		rtp, _, readErr := track.ReadRTP()
		if readErr != nil {
			return errors.Wrapf(err, "读取RTP失败")
		}
		switch track.Kind() {
		case webrtc.RTPCodecTypeAudio:
			self.saver.PushOpus(rtp)
		case webrtc.RTPCodecTypeVideo:
			self.saver.Push264(rtp)
		}
	}

	return self.ctx.Err()
}

func newPeerConnection(configuration webrtc.Configuration) (*webrtc.PeerConnection, error) {
	m := &webrtc.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		return nil, err
	}

	for _, extension := range []string{sdp.SDESMidURI, sdp.SDESRTPStreamIDURI, sdp.TransportCCURI} {
		if extension == sdp.TransportCCURI {
			continue
		}
		if err := m.RegisterHeaderExtension(webrtc.RTPHeaderExtensionCapability{URI: extension}, webrtc.RTPCodecTypeVideo); err != nil {
			return nil, err
		}
	}

	// https://github.com/pion/ion/issues/130
	// https://github.com/pion/ion-sfu/pull/373/files#diff-6f42c5ac6f8192dd03e5a17e9d109e90cb76b1a4a7973be6ce44a89ffd1b5d18R73
	for _, extension := range []string{sdp.SDESMidURI, sdp.SDESRTPStreamIDURI, sdp.AudioLevelURI} {
		if extension == sdp.AudioLevelURI {
			continue
		}
		if err := m.RegisterHeaderExtension(webrtc.RTPHeaderExtensionCapability{URI: extension}, webrtc.RTPCodecTypeAudio); err != nil {
			return nil, err
		}
	}

	i := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		return nil, err
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i))
	return api.NewPeerConnection(configuration)
}
