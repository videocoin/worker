package transcode

import (
	"testing"

	"github.com/VideoCoin/transcode"
)

func Test_generatePlaylist(t *testing.T) {

	s, err := transcode.New()
	if err != nil {
		t.Errorf("failed to make new transcod service: %s", err.Error())
	}

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"main",
			args{
				filename: "./test.m3u8",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := (tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("generatePlaylist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
