package CHelperUnique

import "github.com/rs/xid"

type Unique struct{}

func NewUnique() interface{} { return new(Unique) }

func (h *Unique) ID() string {
	guid := xid.New()
	return guid.String()
}

//x1:=65+rand.Intn(10)
//x2:=65+rand.Intn(26)
//fmt.Println(fmt.Sprintf("%c%c_%03d%s",x1,x2,9999,time.Now().Format("150405")))
