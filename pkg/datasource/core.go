// Package datasource generates an arbitrary time series, which contains one data point per minute.
// The time series is supposed to be deterministic so that the value at a given timestamp is always the same
// for every query.
package datasource

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

type (
	DataSource struct{}

	DataPoint struct {
		Timestamp time.Time
		Value     float64
	}

	Cursor struct {
		rand          *rand.Rand
		tsRand        *rand.Rand
		base          float64
		altPeriod     time.Duration
		altBase       float64
		altMultiplier float64
		cycle         time.Time
		offset        time.Duration
		curTimestamp  time.Time
		endTimestamp  time.Time
		baseCache     map[time.Time]float64
		altBaseCache  map[time.Time]float64
	}
)

func New() *DataSource {
	src := &DataSource{}
	return src
}

const (
	timeStep   = time.Minute
	timePeriod = time.Hour * 24 // time series has a period of 24 hours
)

func (ds *DataSource) Query(begin, end time.Time) (*Cursor, error) {
	if !begin.Before(end) {
		return nil, errors.New("'end' is not after 'begin'")
	}

	chunkBegin := begin.Truncate(timePeriod)
	gen := rand.New(rand.NewSource(chunkBegin.Unix()))

	cursor := &Cursor{
		rand:          gen,
		cycle:         chunkBegin,
		base:          math.Cos(0.0),
		altBase:       math.Cos(1.0),
		altMultiplier: gen.Float64(),
		altPeriod:     time.Duration(gen.Intn(int(timePeriod) / 3)),
		curTimestamp:  chunkBegin,
		endTimestamp:  end,
		baseCache:     make(map[time.Time]float64),
		altBaseCache:  make(map[time.Time]float64),
	}

	for cursor.curTimestamp.Before(begin) {
		cursor.Next()
	}

	return cursor, nil
}

func (c *Cursor) resetCycle() {
	c.cycle = c.cycle.Add(timePeriod)
	c.rand = rand.New(rand.NewSource(c.cycle.Unix()))
	c.base = math.Cos(0.0)
	c.altBase = math.Cos(1.0)
	c.altMultiplier = c.rand.Float64()
	c.altPeriod = time.Duration(c.rand.Intn(int(timePeriod)/3)) + time.Hour*3
	c.offset = 0
}

func (c *Cursor) updateBaseValues() {
	if cachedBase, ok := c.baseCache[c.curTimestamp]; ok {
		c.base = cachedBase
	} else {
		c.base = math.Cos(math.Pi * 2.0 * (float64(c.offset) / float64(timePeriod)))
		c.baseCache[c.curTimestamp] = c.base
	}

	if cachedAltBase, ok := c.altBaseCache[c.curTimestamp]; ok {
		c.altBase = cachedAltBase
	} else {
		c.altBase = math.Cos(1.0 + math.Pi*2.0*(float64(c.offset)/float64(c.altPeriod)))
		c.altBaseCache[c.curTimestamp] = c.altBase
	}
}

func (c *Cursor) Next() (*DataPoint, bool) {
	if c.curTimestamp.After(c.endTimestamp) {
		return nil, false
	}

	val := (c.base+1.5)*50.0 + (c.altBase+1.2)*c.altMultiplier*50
	dp := &DataPoint{Timestamp: c.curTimestamp, Value: val}

	timeInc := timeStep + (time.Duration(c.rand.Intn(55)-17) * time.Second)
	c.curTimestamp = c.curTimestamp.Add(timeInc)

	c.offset += timeStep

	if c.offset >= timePeriod {
		// Reset for the next cycle.
		c.resetCycle()
	} else {
		// update base
		c.updateBaseValues()
	}

	return dp, true
}
