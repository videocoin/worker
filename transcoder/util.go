package transcoder

import (
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

func checkSource(url string) error {
	if strings.HasPrefix(url, "file://") || strings.HasPrefix(url, "/") {
		fp := strings.TrimPrefix(url, "file://")
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			return err
		}
	} else if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		hc := http.Client{
			Timeout: time.Duration(5 * time.Second),
		}
		resp, err := hc.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp != nil && resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get %s, return status %s", url, resp.Status)
		}
	} else {
		return errors.New("unknown source type")
	}

	return nil
}

func SearchBigInts(a []*big.Int, x *big.Int) int {
	return sort.Search(len(a), func(i int) bool {
		return a[i].Cmp(x) >= 0
	})
}
