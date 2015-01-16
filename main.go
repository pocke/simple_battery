package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"unsafe"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/pocke/goevent"
)

func gthread(f func()) {
	gdk.ThreadsEnter()
	defer gdk.ThreadsLeave()
	f()
}

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
		var icon *gtk.StatusIcon
		gthread(func() {
			icon = gtk.NewStatusIconFromStock(gtk.STOCK_FILE)
		})
		icon.SetTitle(fmt.Sprint("%d", n))
		mu.Lock()
		defer mu.Unlock()
		statusIcons[n] = icon
	})

	e.On("delete", func(n int) {
		log.Printf("delete %d\n", n)
		mu.Lock()
		defer mu.Unlock()
		gthread(func() {
			glib.ObjectFromNative(unsafe.Pointer(statusIcons[n].GStatusIcon)).Unref()
		})
		delete(statusIcons, n)
	})

	e.On("change", func(n, v int) {
		log.Printf("change %d %d\n", n, v)
		gthread(func() {
			mu.Lock()
			mu.Unlock()
			statusIcons[n].SetTooltipText(fmt.Sprintf("%d", v))
		})
	})

	WatchBattery(e)

	gtk.Main()
}
