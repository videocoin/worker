package v1

import (
	fmt "fmt"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	profiles_v1 "github.com/videocoin/cloud-api/profiles/v1"
)

// CheckErrs used to ensure a job is valid before passing to transcoder
func (w *WorkOrder) CheckErrs() []error {
	var errs = make([]error, 0)

	err := fmt.Errorf("invalid profile")
	for _, p := range []profiles_v1.ProfileId{profiles_v1.ProfileIdFHD, profiles_v1.ProfileIdHD, profiles_v1.ProfileIdSD} {
		if p == w.ProfileId {
			err = nil
		}
	}

	if err != nil {
		errs = append(errs, err)
	}

	if w.ClientAddress == "" {
		errs = append(errs, fmt.Errorf("invalid client address"))
	}

	if w.StreamId <= 0 {
		errs = append(errs, fmt.Errorf("invalid stream id"))
	}

	if w.StreamAddress == "" {
		errs = append(errs, fmt.Errorf("invalid stream address"))
	}

	return errs
}

func (w *WorkOrder) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()

	err := scope.SetColumn("id", uuid.String())
	if err != nil {
		return err
	}

	return scope.SetColumn("created_at", time.Now().Unix())
}
