package transcoder

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/videocoin/transcode/transcoder/hlswatcher"

	"github.com/sirupsen/logrus"
	validatorv1 "github.com/videocoin/cloud-api/validator/v1"
	"github.com/videocoin/cloud-pkg/retry"
)

func checkSource(url string) error {
	if strings.HasPrefix(url, "file://") || strings.HasPrefix(url, "/") {
		fp := strings.TrimPrefix(url, "file://")
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			return err
		}
	} else if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		hc := http.Client{
			Timeout: 5 * time.Second,
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
	for idx, item := range a {
		if x.Cmp(item) == 0 {
			return idx
		}
	}
	return -1
}

func (t *Transcoder) waitGetInChunks(segmentNum uint64) (bool, error) {
	logger := t.logger.WithFields(logrus.Fields{
		"task_id": t.task.ID,
		"segment": segmentNum,
	})

	idx := -1
	counter := 0
	for {
		chunks, err := t.sc.GetInChunks()
		if err != nil {
			logger.Errorf("failed to get in chunks: %s", err)
			return false, err
		}

		logger.Debugf("GetInChunks: %+v\n", chunks)

		if len(chunks) > 0 {
			idx = SearchBigInts(chunks, big.NewInt(int64(segmentNum)))
			if idx >= 0 {
				break
			}
			if counter >= 10 {
				err = fmt.Errorf("failed to search in chunks: segment %d", segmentNum)
				logger.Error(err)
				return false, err
			}
		}

		counter++

		time.Sleep(time.Second * 5)
	}

	return idx >= 0, nil
}

func (t *Transcoder) submitAndValidateProof(segment *hlswatcher.SegmentInfo) error {
	logger := t.logger.WithFields(logrus.Fields{
		"task_id": t.task.ID,
		"segment": segment.Num,
	})
	logger.Info("submitting proof")

	inChunkID := big.NewInt(int64(segment.Num))
	outChunkID := inChunkID

	profileID := new(big.Int)
	profiles, _ := t.sc.GetProfiles()
	if len(profiles) > 0 {
		profileID = profiles[0]
	}

	tx, err := t.sc.SubmitProof(inChunkID, outChunkID, profileID)
	if err != nil {
		logger.Errorf("failed to submit proof: %s", err)
		return err
	}

	logger.Debugf("submitting proof tx %+v\n", tx.Hash().String())

	if t.task != nil {
		logger.Info("validating proof")

		ctx := context.Background()
		vpReq := &validatorv1.ValidateProofRequest{
			StreamContractAddress: t.task.StreamContractAddress,
			ProfileId:             profileID.Bytes(),
			OutputChunkId:         outChunkID.Bytes(),
			StreamId:              t.task.StreamID,
			IsLast:                segment.IsLast,
		}
		_, err = t.dispatcher.ValidateProof(ctx, vpReq)
		if err != nil {
			logger.Errorf("failed to validate proof: %s", err)
			return err
		}
	}

	return nil
}

func (t *Transcoder) uploadSegment(segment *hlswatcher.SegmentInfo) error {
	logger := t.logger.WithFields(logrus.Fields{
		"task_id": t.task.ID,
		"segment": segment.Num,
	})
	logger.Info("uploading segment")

	err := retry.RetryWithAttempts(5, time.Second*1, func() error {
		return t.uploadSegmentViaHTTP(t.task, segment)
	})
	if err != nil {
		logger.Errorf("failed to upload segment: %s", err)
		return err
	}

	logger.Info("segment has been uploaded")

	return nil
}
