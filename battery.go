package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pocke/goevent"
)
import "time"

func WatchBattery(e goevent.Table) {
	go func() {
		b := NewBatteries(e)
		b.Update()

		t := time.Tick(1 * time.Minute)
		for {
			select {
			case <-t:
				b.Update()
			}
		}
	}()
}

const (
	Base = "/sys/class/power_supply"
)

type Batteres struct {
	e goevent.Table
	n []int
}

func NewBatteries(e goevent.Table) *Batteres {
	return &Batteres{
		e: e,
	}
}

func (b *Batteres) Update() {
	l := b.list()

	deleted := sliceDiff(b.n, l)
	added := sliceDiff(l, b.n)

	b.n = l
	for _, v := range added {
		b.e.Trigger("add", v)
	}
	for _, v := range deleted {
		b.e.Trigger("delete", v)
	}

	for _, v := range l {
		i, err := b.get(v)
		if err != nil {
			// Retry
			b.Update()
			return
		}

		b.e.Trigger("change", v, i)
	}

}

// TODO: error handling
func (_ *Batteres) list() []int {
	m, _ := filepath.Glob(fmt.Sprintf("%s/BAT*", Base))
	res := make([]int, 0, len(m))
	re := regexp.MustCompile(`BAT(\d+)$`)
	for _, v := range m {
		i, _ := strconv.Atoi(re.FindStringSubmatch(v)[1])
		res = append(res, i)
	}

	return res
}

func (_ *Batteres) get(i int) (int, error) {
	f, err := os.Open(fmt.Sprintf("%s/BAT%d/capacity", Base, i))
	if err != nil {
		return -1, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return -1, err
	}
	res, err := strconv.Atoi(strings.Trim(string(b), "\n"))
	if err != nil {
		return -1, err
	}

	return res, nil
}

func include(l []int, n int) int {
	for i, v := range l {
		if v == n {
			return i
		}
	}
	return -1
}

func sliceDiff(a, b []int) []int {
	res := make([]int, len(a))
	copy(res, a)
	for _, v := range b {
		if i := include(res, v); i != -1 {
			res = append(res[:i], res[i+1:]...)
		}
	}

	return res
}
