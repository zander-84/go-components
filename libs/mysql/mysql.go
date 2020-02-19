package CMysql

type Mysql interface {
	Obj() interface{}
	Transaction(f func(tx interface{}) (int, error))(int, error)
}
