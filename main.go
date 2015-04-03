package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"unsafe"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/pocke/goevent"
)

func gthread(f func()) {
	gdk.ThreadsEnter()
	defer gdk.ThreadsLeave()
	f()
}

var (
	iconPixbufs = make([]*gdkpixbuf.Pixbuf, 5)
)

func getIcon(v int) *gdkpixbuf.Pixbuf {
	if v <= 0 {
		v = 1
	}
	n := int(math.Ceil(float64(v) / 20.0))
	if iconPixbufs[n-1] != nil {
		return iconPixbufs[n-1]
	}

	loader, _ := gdkpixbuf.NewLoaderWithType("png")
	f, err := Asset(fmt.Sprintf("assets/battery-bar-%d-icon.png", n))
	if err != nil {
		log.Println(err)
		f, _ = Asset("assets/battery-bar-1-icon.png")
	}
	loader.Write(f)
	iconPixbufs[n-1] = loader.GetPixbuf()
	return iconPixbufs[n-1]
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
		log.Printf("add BAT%d\n", n)
		var icon *gtk.StatusIcon

		gthread(func() {
			icon = gtk.NewStatusIconFromPixbuf(getIcon(100))
		})
		icon.SetTitle(fmt.Sprint("BAT%d", n))
		mu.Lock()
		defer mu.Unlock()
		statusIcons[n] = icon
	})

	e.On("delete", func(n int) {
		log.Printf("delete BAT%d\n", n)
		mu.Lock()
		defer mu.Unlock()
		gthread(func() {
			glib.ObjectFromNative(unsafe.Pointer(statusIcons[n].GStatusIcon)).Unref()
		})
		delete(statusIcons, n)
	})

	e.On("change", func(n, v int) {
		log.Printf("change BAT%d %d\n", n, v)
		gthread(func() {
			mu.Lock()
			defer mu.Unlock()
			icon := statusIcons[n]
			icon.SetTooltipText(fmt.Sprintf("BAT%d: %d", n, v))
			icon.SetFromPixbuf(getIcon(v))
		})
	})

	WatchBattery(e)

	gtk.Main()
}
