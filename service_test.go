package transcode

import "testing"

func Test_generatePlaylist(t *testing.T) {
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
			if err := generatePlaylist(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("generatePlaylist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
