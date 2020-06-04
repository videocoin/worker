package hw

import (
	"io/ioutil"
	"os"
	"strings"
)

func GetDeviceModel() string {
	const filename = "/proc/device-tree/model"
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		f, err := os.Open(filename)
		if err != nil {
			return ""
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			return ""
		}

		return strings.ToLower(string(b))
	}

	return ""
}

func IsRaspberry() bool {
	return strings.Contains(GetDeviceModel(), "raspberry")
}

func IsJetson() bool {
	return strings.Contains(GetDeviceModel(), "jetson")
}
