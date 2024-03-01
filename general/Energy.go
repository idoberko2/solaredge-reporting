package general

import "time"

type Energy struct {
	DateTime time.Time `db:"t"`
	Value    int
}

func (e Energy) Empty() bool {
	return e.DateTime.IsZero()
}
