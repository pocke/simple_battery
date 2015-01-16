package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/pocke/goevent"
)

func main() {
	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(&os.Args)
	log.SetFlags(log.Llongfile)

	var mu sync.Mutex
	statusIcons := make(map[int]*gtk.StatusIcon)

	e := goevent.NewTable()
	e.On("add", func(n int) {
		log.Printf("add %d\n", n)
		gdk.ThreadsEnter()
		icon := gtk.NewStatusIconFromStock(gtk.STOCK_FILE)
		gdk.ThreadsLeave()
		icon.SetTitle(fmt.Sprint("%d", n))
		mu.Lock()
		defer mu.Unlock()
		statusIcons[n] = icon
	})

	e.On("delete", func(n int) {
		log.Printf("delete %d\n", n)
		mu.Lock()
		defer mu.Unlock()
		delete(statusIcons, n)
	})

	e.On("change", func(n, v int) {
		log.Printf("change %d %d\n", n, v)
		mu.Lock()
		mu.Unlock()
		gdk.ThreadsEnter()
		gdk.ThreadsLeave()
		statusIcons[n].SetTooltipText(fmt.Sprintf("%d", v))
	})

	WatchBattery(e)

	gtk.Main()
}
