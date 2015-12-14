package main

import (
	"fmt"
	"io"
	"sort"
	"sync"
)

type EC2List struct {
	sync.RWMutex
	list []InstanceInfo
}

func (l *EC2List) Add(infos ...InstanceInfo) {
	l.Lock()
	defer l.Unlock()
	l.list = append(l.list, infos...)
}

func (l *EC2List) Sort() {
	sort.Sort(l)
}

func (l *EC2List) Output(w io.Writer) {
	for _, info := range l.list {
		fmt.Fprintf(w, info.String())
	}
}

func (l *EC2List) Len() int {
	return len(l.list)
}

func (l *EC2List) Swap(i, j int) {
	l.list[i], l.list[j] = l.list[j], l.list[i]
}

func (l *EC2List) Less(i, j int) bool {
	return l.list[i].LowerName() < l.list[j].LowerName()
}
