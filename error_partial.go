package gads

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type PartialFailureErrors []*PartialFailureError

// Error returns all the errors combined
func (p PartialFailureErrors) Error() string {
	var m []byte
	for i := range p {
		if i != 0 {
			m = append(m, " - "...)
		}
		m = append(m, p[i].Error()...)
	}
	return string(m)
}

// PartialFailureError represents a
type PartialFailureError struct {
	FieldPath string `xml:"fieldPath,omitempty"`
	Trigger   string `xml:"trigger,omitempty"`
	Type      string `xml:"type,attr,omitempty"`
	Reason    string `xml:"reason,omitempty"`
	Code      string `xml:"errorString,omitempty"`
	Offset    *int   `xml:"-"`
}

// Error return a summary of the error
func (p *PartialFailureError) Error() string {
	m := []byte(p.Code + " @ " + p.FieldPath)
	if p.Trigger != "" {
		m = append(m, " ; trigger:'"+p.Trigger+"'"...)
	}
	return string(m)
}

var offsetParse = regexp.MustCompile(`(?:operations)\[([0-9]+)\](?:\.(?:[^ ;]+))?`)

// GetRequestOffset returns the offset of the current partial error
func (p *PartialFailureError) GetRequestOffset() (int, error) {
	if p.Offset != nil {
		return *p.Offset, nil
	}

	a := offsetParse.FindStringSubmatch(p.FieldPath)

	if len(a) != 2 {
		return 0, errors.New("unable to find offset")
	}

	var offset int
	var err error
	// one capturing group, composed of digits
	offset, err = strconv.Atoi(a[1])
	if err != nil {
		return 0, fmt.Errorf("unable to parse offset from %v", a[1])
	}

	return offset, nil

}
