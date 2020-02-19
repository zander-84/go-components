package CCron

import "errors"

var (
	ErrID = errors.New("Error ID")
	ErrIDExist = errors.New("ID exists")
)


/* Spec
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
*/
type Job struct {
	ID   string
	Desc  string
	Spec  string
	Cmd   Cmd
	Obj interface{}
}


type Crontab interface {
	// 添加任务
	AddJob(jobInfo *Job) error

	// 移除
	Remove(id string)

	//
	Status() map[string]*Job
	Start()
	Stop()

	Restart(id string) error

	Obj() interface{}
}
