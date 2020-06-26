package server

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type blockTimesParams struct {
	Limit int64 `form:"limit"`
}

type accountsIndexParams struct {
	Height int64 `form:"height"`
}

func (p *blockTimesParams) setDefaults() {
	if p.Limit <= 1 {
		p.Limit = 100
	}
	if p.Limit > 1000 {
		p.Limit = 1000
	}
}

type timeBucket struct {
	Interval string `form:"interval"`
	Period   uint   `form:"period"`
}

func (t *timeBucket) validate() error {
	if t.Interval == "" {
		t.Interval = "h"
	}

	switch t.Interval {
	case "h":
		if t.Period == 0 {
			t.Period = 48
		}
		if t.Period > 120 {
			t.Period = 120
		}
	case "d":
		if t.Period == 0 {
			t.Period = 30
		}
		if t.Period > 100 {
			t.Period = 100
		}
	default:
		return errors.New("invalid interval: " + t.Interval)
	}

	return nil
}

func getTimeBucket(c *gin.Context) (t timeBucket, err error) {
	if err = c.BindQuery(&t); err != nil {
		return
	}
	err = t.validate()
	return
}
