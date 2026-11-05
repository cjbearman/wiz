package wiz

import "time"

// zeroableTicker mimics a ticker, but can be configured with a zero duration to effectively disable ticking.
type zeroableTicker struct {
	C                <-chan time.Time
	underlyingTicker *time.Ticker
}

// NewZeroableTicker returns a new ticker that will not tick if the duration is zero
func NewZeroableTicker(t time.Duration) zeroableTicker {

	zt := &zeroableTicker{}

	if t != 0 {
		zt.underlyingTicker = time.NewTicker(t)
		zt.C = zt.underlyingTicker.C
	} else {
		zt.C = make(chan time.Time)
	}

	return *zt
}

// Stop the ticker
func (nt *zeroableTicker) Stop() {
	if nt.underlyingTicker != nil {
		nt.underlyingTicker.Stop()
	}
}

// Reset the ticker
func (nt *zeroableTicker) Reset(d time.Duration) {
	if nt.underlyingTicker != nil {
		nt.underlyingTicker.Reset(d)
	} else {
		panic("cannot reset a zeroableTicker that was configured with zero duration")
	}
}
