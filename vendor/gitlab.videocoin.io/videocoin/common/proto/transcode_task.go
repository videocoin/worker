package proto

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"gitlab.videocoin.io/videocoin/common/protoutil"
)

func (task *TranscodeTask) Value() (driver.Value, error) {
	m := &protoutil.JSONPb{OrigName: true, EmitDefaults: true, EnumsAsInts: true}
	return m.Marshal(task)
}

func (task *TranscodeTask) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed.")
	}

	m := &protoutil.JSONPb{OrigName: true, EmitDefaults: true, EnumsAsInts: true}
	return m.Unmarshal(source, task)
}

func (task *TranscodeTask) IsOutputHLS() bool {
	return task.Output.Video[0].Format == "hls"
}

func (task *TranscodeTask) IsOutputMP4() bool {
	return task.Output.Video[0].Format == "mp4"
}

func (task *TranscodeTask) IsOutputFLV() bool {
	return task.Output.Video[0].Format == "flv"
}

func (task *TranscodeTask) IsOutputRTMP() bool {
	return task.IsOutputFLV() && task.Output.Video[0].GetRtmpAddress() != ""
}

func (task *TranscodeTask) GetAngles(video *VideoOutput) []uint32 {
	angles := []uint32{}

	angle := video.Yaw
	delta := 360 / video.TotalStreams
	step := video.TotalStreams / video.Streams

	for {
		angles = append(angles, angle)
		angle += delta * step
		if angle >= 360 {
			break
		}
	}

	return angles
}

func (task *TranscodeTask) GetOutputDir(baseDir string) string {
	return fmt.Sprintf(
		"%s/%d/%s",
		baseDir,
		task.Metadata.UserId,
		task.Metadata.BroadcastId,
	)
}

func (task *TranscodeTask) GetOutputFilename(video *VideoOutput, angle uint32) string {
	return fmt.Sprintf("%s-.%s", task.GetSessionID(video, angle), getExtForFormat(video.Format))
}

func (task *TranscodeTask) GetOutputStatPath(baseDir string) string {
	return fmt.Sprintf("%s/stat.txt", task.GetOutputDir(baseDir))
}

func (task *TranscodeTask) GetSessionID(video *VideoOutput, angle uint32) string {
	return fmt.Sprintf(
		"%s-%dx%d-%dK-%d",
		task.Metadata.PlatformId,
		video.Width,
		video.Height,
		video.Bitrate,
		angle,
	)
}

func (task *TranscodeTask) GetSegmentObjectName(video *VideoOutput, angle uint32, name string) string {
	sessionID := task.GetSessionID(video, angle)
	return fmt.Sprintf(
		"%s/%d/%d/%s/%s/%s",
		task.Output.Video[0].Format,
		task.Metadata.UserId,
		task.Metadata.EventId,
		task.Metadata.BroadcastId,
		sessionID,
		name,
	)
}

func (task *TranscodeTask) GetObjectName(video *VideoOutput, angle uint32) string {
	return fmt.Sprintf(
		"%s/%d/%d/%s/%s",
		task.Output.Video[0].Format,
		task.Metadata.UserId,
		task.Metadata.EventId,
		task.Metadata.BroadcastId,
		task.GetOutputFilename(video, angle),
	)
}

func getExtForFormat(format string) string {
	switch format {
	case "hls":
		return "m3u8"
	default:
		return format
	}
}
