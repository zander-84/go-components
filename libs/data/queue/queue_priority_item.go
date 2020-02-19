package CQueue

import  "github.com/zander-84/go-components/libs/data/link"

type tempPriorityItem struct {
	value interface{}
	priority int
}

func newTempPriorityItem(value interface{}, priority int) *tempPriorityItem {
	return &tempPriorityItem{
		value: value,
		priority:priority,
	}
}
type priorityItem struct {
	value    CLink.LinkList
	priority int
}
func newPriorityItem(value interface{}, priority int) *priorityItem {
	link:=CLink.NewLinkList(CLink.DoubleLink)
	link.Push(value)
	return &priorityItem{
		value:    link,
		priority: priority,
	}
}
func (this *priorityItem) Less(than *priorityItem) bool {
	return this.priority > than.priority
}



