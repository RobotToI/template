package util

import (
	"fmt"
	"strings"
	"time"
)

// CustomTime is a custom time format
type CustomTime time.Time

const ctLayout = "2006-01-02"

// ToTime returns the time.Time
func (ct *CustomTime) ToTime() time.Time {
	return time.Time(*ct)
}

// UnmarshalJSON unmarshals the time
func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(ctLayout, s)
	*ct = CustomTime(nt)
	return
}

// MarshalJSON marshals the time
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

// String returns the time in the custom format
func (ct *CustomTime) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(ctLayout))
}
