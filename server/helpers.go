package server

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ridEmpty   = iota
	ridNumeric = iota
	ridString  = iota
)

var (
	reNumeric = regexp.MustCompile(`^[0-9]+$`)
)

type rid struct {
	kind int
	raw  string
}

func (r rid) IsNumeric() bool {
	return r.kind == ridNumeric
}

func (r rid) IsString() bool {
	return r.kind == ridString
}

func (r rid) String() string {
	return r.raw
}

func (r rid) Number() int64 {
	v, _ := strconv.ParseInt(r.raw, 10, 64)
	return v
}

func resourceID(c *gin.Context, key string) rid {
	val := strings.TrimSpace(c.Param(key))

	id := rid{raw: val, kind: ridString}
	if val == "" {
		id.kind = ridEmpty
	}
	if reNumeric.MatchString(val) {
		id.kind = ridNumeric
	}

	return id
}
