package CHelperUnique

import "github.com/rs/xid"

type Unique struct{}

func NewUnique() interface{} { return new(Unique) }

func (h *Unique) ID() string {
	guid := xid.New()
	return guid.String()
}
