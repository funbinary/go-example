package opus

import (
	"context"
	"fmt"
	"strings"

	"github.com/bin-work/go-example/pkg/pion/srs"

	"github.com/bin-work/go-example/pkg/pion"
	"github.com/pion/sdp/v3"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RtcOgg struct {
	Host     string
	Room     string
	Display  string
	rtcUrl   string
	SavePath string
	ctx      context.Context
	cancel   context.CancelFunc
	pc       *webrtc.PeerConnection
	saver    *oggwriter.OggWriter
}

func NewOgg(host, room, display, savePath string) (*RtcOgg, error) {
	var err error
	stream := &RtcOgg{
		Host:     host,
		Room:     room,
		Display:  display,
		rtcUrl:   "webrtc://" + host + "/" + room + "/" + display,
		SavePath: savePath,
	}
	logrus.Info("拉取", stream.rtcUrl+"rtc流")
	stream.ctx, stream.cancel = context.WithCancel(context.Background())
	if err != nil {
		panic(err)
	}
	//创建PeerConncetion
	stream.pc, err = pion.NewPeerConnection(webrtc.Configuration{}, true, true)
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
		if codec.MimeType == "audio/opus" {
			stream.saver, err = oggwriter.New(savePath+".ogg", codec.ClockRate, codec.Channels)
			fmt.Println(err)
			err = stream.writeTrackToDisk(track)
			if err != nil {
				codec := track.Codec()
				err = errors.Wrapf(err, "Handle  track %v, pt=%v", codec.MimeType, codec.PayloadType)
			}
		}

	})
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

func (self *RtcOgg) Stop() {
	self.pc.Close()
	self.cancel()
	self.saver.Close()
}

func (self *RtcOgg) writeTrackToDisk(track *webrtc.TrackRemote) error {
	for self.ctx.Err() == nil {
		pkt, _, err := track.ReadRTP()
		if err != nil {
			if self.ctx.Err() != nil {
				return nil
			}
			return errors.Wrapf(err, "Read RTP")
		}
		if err = self.saver.WriteRTP(pkt); err != nil {
			if len(pkt.Payload) <= 2 {
				continue
			}
			return errors.Wrapf(err, "write rtp failed")
		}
	}

	return self.ctx.Err()
}
