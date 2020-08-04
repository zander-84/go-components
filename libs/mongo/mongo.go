package CMongo

type Mongo interface {
	Obj() interface{}
	DB() interface{} //*Database
}
