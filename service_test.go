package transcode

import (
	"reflect"
	"testing"
)

func Test_newService(t *testing.T) {
	tests := []struct {
		name    string
		want    *Service
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newService()
			if (err != nil) != tt.wantErr {
				t.Errorf("newService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newService() = %v, want %v", got, tt.want)
			}
		})
	}
}
