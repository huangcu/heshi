package main

import (
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

func uniq(strings []string) (ret []string) {
	return
}

func authors() []string {
	if b, err := exec.Command("git", "log").Output(); err == nil {
		lines := strings.Split(string(b), "\n")

		var a []string
		r := regexp.MustCompile(`^Author:\s*([^ <]+).*$`)
		for _, e := range lines {
			ms := r.FindStringSubmatch(e)
			if ms == nil {
				continue
			}
			a = append(a, ms[1])
		}
		sort.Strings(a)
		var p string
		lines = []string{}
		for _, e := range a {
			if p == e {
				continue
			}
			lines = append(lines, e)
			p = e
		}
		return lines
	}
	return []string{"Yasuhiro Matsumoto <mattn.jp@gmail.com>"}
}

func main() {
	var menuitem *gtk.MenuItem
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("GTK Go!")
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		// println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")

	//--------------------------------------------------------
	// GtkVBox
	//--------------------------------------------------------
	vbox := gtk.NewVBox(false, 1)

	//--------------------------------------------------------
	// GtkMenuBar
	//--------------------------------------------------------
	menubar := gtk.NewMenuBar()
	vbox.PackStart(menubar, false, false, 0)

	//--------------------------------------------------------
	// GtkVPaned
	//--------------------------------------------------------
	vpaned := gtk.NewVPaned()
	vbox.Add(vpaned)

	//--------------------------------------------------------
	// GtkFrame
	//--------------------------------------------------------
	frame1 := gtk.NewFrame("")
	framebox1 := gtk.NewVBox(false, 1)
	frame1.Add(framebox1)

	frame2 := gtk.NewFrame("Demo")
	framebox2 := gtk.NewVBox(false, 1)
	frame2.Add(framebox2)

	//--------------------------------------------------------
	// GtkTextView
	//--------------------------------------------------------
	swinContent := gtk.NewScrolledWindow(nil, nil)
	swinContent.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swinContent.SetShadowType(gtk.SHADOW_IN)

	textviewContent := gtk.NewTextView()
	var start, end gtk.TextIter
	bufferContent := textviewContent.GetBuffer()
	bufferContent.GetStartIter(&start)
	bufferContent.Insert(&start, "Hello ")
	bufferContent.GetEndIter(&end)
	bufferContent.Insert(&end, "World!")
	tag := bufferContent.CreateTag("bold", map[string]string{"background": "#FF0000", "weight": "700"})
	bufferContent.GetStartIter(&start)
	bufferContent.GetEndIter(&end)
	bufferContent.ApplyTag(tag, &start, &end)
	swinContent.Add(textviewContent)
	framebox2.Add(swinContent)
	framebox2.SetSizeRequest(600, 400)

	swinInput := gtk.NewScrolledWindow(nil, nil)
	swinInput.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swinInput.SetShadowType(gtk.SHADOW_IN)

	textviewInput := gtk.NewTextView()
	swinInput.Add(textviewInput)
	framebox1.Add(swinInput)

	//--------------------------------------------------------
	// GtkHBox
	//--------------------------------------------------------
	buttons := gtk.NewHBox(false, 1)
	//--------------------------------------------------------
	// GtkButton
	//--------------------------------------------------------
	buttonSend := gtk.NewButtonWithLabel("Send")
	buttonSend.Clicked(func() {
		bufferContent.GetEndIter(&end)
		bufferInput := textviewInput.GetBuffer()
		var startInput, endInput gtk.TextIter
		bufferInput.GetStartIter(&startInput)
		bufferInput.GetEndIter(&endInput)
		bufferContent.Insert(&end, "\n"+bufferInput.GetText(&startInput, &endInput, false))
		bufferInput.Delete(&startInput, &endInput)
	})
	buttonSend.SetSizeRequest(50, 10)
	buttonClose := gtk.NewButtonWithLabel("Close")
	buttonClose.Clicked(func() {
		window.Destroy()
	})
	buttonClose.SetSizeRequest(50, 10)
	buttons.SetSizeRequest(100, 10)
	buttons.Add(buttonSend)
	buttons.Add(buttonClose)
	framebox1.Add(buttons)
	framebox1.SetSizeRequest(600, 100)
	// framebox1.PackStart(buttons, false, false, 0)
	// buffer.Connect("changed", func() {
	// 	println("changed")
	// })

	vpaned.Pack1(frame2, false, false)
	vpaned.Pack2(frame1, false, false)
	//--------------------------------------------------------
	// GtkMenuItem
	//--------------------------------------------------------
	cascademenu := gtk.NewMenuItemWithMnemonic("_File")
	menubar.Append(cascademenu)
	submenu := gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	menuitem = gtk.NewMenuItemWithMnemonic("E_xit")
	menuitem.Connect("activate", func() {
		gtk.MainQuit()
	})
	submenu.Append(menuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_View")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	checkmenuitem := gtk.NewCheckMenuItemWithMnemonic("_Disable")
	checkmenuitem.Connect("activate", func() {
		vpaned.SetSensitive(!checkmenuitem.GetActive())
	})
	submenu.Append(checkmenuitem)

	menuitem = gtk.NewMenuItemWithMnemonic("_Font")
	menuitem.Connect("activate", func() {
		fsd := gtk.NewFontSelectionDialog("Font")
		// fsd.SetFontName(fontbutton.GetFontName())
		fsd.Response(func() {
			println(fsd.GetFontName())
			// fontbutton.SetFontName(fsd.GetFontName())
			fsd.Destroy()
		})
		fsd.SetTransientFor(window)
		fsd.Run()
	})
	submenu.Append(menuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_Help")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	menuitem = gtk.NewMenuItemWithMnemonic("_About")
	menuitem.Connect("activate", func() {
		dialog := gtk.NewAboutDialog()
		dialog.SetName("Go-Gtk Demo!")
		dialog.SetProgramName("demo")
		dialog.SetAuthors(authors())
		dir, _ := path.Split(os.Args[0])
		imagefile := path.Join(dir, "../../data/mattn-logo.png")
		pixbuf, _ := gdkpixbuf.NewPixbufFromFile(imagefile)
		dialog.SetLogo(pixbuf)
		dialog.SetLicense("The library is available under the same terms and conditions as the Go, the BSD style license, and the LGPL (Lesser GNU Public License). The idea is that if you can use Go (and Gtk) in a project, you should also be able to use go-gtk.")
		dialog.SetWrapLicense(true)
		dialog.Run()
		dialog.Destroy()
	})
	submenu.Append(menuitem)

	//--------------------------------------------------------
	// GtkStatusbar
	//--------------------------------------------------------
	statusbar := gtk.NewStatusbar()
	contextID := statusbar.GetContextId("go-gtk")
	statusbar.Push(contextID, "GTK binding for Go!")

	framebox1.PackStart(statusbar, false, false, 0)

	//--------------------------------------------------------
	// Event
	//--------------------------------------------------------
	window.Add(vbox)
	window.SetSizeRequest(600, 600)
	window.ShowAll()
	gtk.Main()
}
