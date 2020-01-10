package lib

import (
	"fmt"
	"time"
)

type memoryUtilizationPoint struct {
	time  *time.Time
	value *float64
}

func (mup memoryUtilizationPoint) toString() string {
	return fmt.Sprintf("%v %v", mup.time, *mup.value)
}
