package register

import "github.com/dati-mipt/consistency-algorithms/util"

type EpidemicRegister struct {
	rid int64

	current util.TimestampedValue

	replicas []util.Receiver
}

func (r EpidemicRegister) Write(value int64) bool {
	r.current.Val = value
	r.current.Ts = util.Timestamp{Number: r.current.Ts.Number + 1, Rid: r.rid}
	return true
}

func (r EpidemicRegister) Read() int64 {
	return r.current.Val
}

func (r EpidemicRegister) Periodically() {
	for _, rep := range r.replicas {
		rep.Message(r.current)
	}
}

func (r EpidemicRegister) Message(msg interface{}) {
	if cast, ok := msg.(util.TimestampedValue); ok {
		r.update(cast)
	}
}

func (r EpidemicRegister) update(u util.TimestampedValue) {
	if r.current.Ts.Less(u.Ts) {
		r.current = u
	}
}