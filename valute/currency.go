package valute

import (
	"fmt"
	"strconv"
	"strings"
)

type ValuteMeta struct {
	Min, Max CurrentValute
	Mean     float64
}

type CurrentValute struct {
	Date  string
	Value float64
}

type Valute struct {
	Name  string `xml:"Name"`
	Value string `xml:"Value"`
}

type Valutes map[string]*ValuteMeta

func NewValutes() Valutes {
	return make(Valutes)
}

func New(date string, value float64) *ValuteMeta {
	return &ValuteMeta{
		Min: CurrentValute{
			Date:  date,
			Value: value,
		},
		Max: CurrentValute{
			Date:  date,
			Value: value,
		},
		Mean: value,
	}
}

func (c *ValuteMeta) Add(date string, value float64) {
	c.CheckMin(CurrentValute{
		Date:  date,
		Value: value,
	})
	c.CheckMax(CurrentValute{
		Date:  date,
		Value: value,
	})
	c.Mean += value
}

func (c *ValuteMeta) GetMean(limit int64) float64 {
	return c.Mean / float64(limit)
}

func (c *ValuteMeta) CheckMin(dv CurrentValute) {
	if c.Min.Value > dv.Value {
		c.Min.Value = dv.Value
		c.Min.Date = dv.Date
	}
}

func (c *ValuteMeta) CheckMax(dv CurrentValute) {
	if c.Max.Value < dv.Value {
		c.Max.Value = dv.Value
		c.Max.Date = dv.Date
	}
}

func (v *Valute) ParseFloat() float64 {
	str := strings.Replace(v.Value, ",", ".", 1)
	val, err := strconv.ParseFloat(str, 64)

	if err != nil {
		panic(fmt.Errorf("parse error %w", err))
	}

	return val
}
