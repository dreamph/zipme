package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dreamph/gozip"
	"log"
)

func chooseDirectory(w fyne.Window, onSelect func(filePath string, fileType int)) {
	dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
		if err != nil {
			return
		}
		onSelect(dir.Path(), Dir)
	}, w)
}

func chooseFiles(w fyne.Window, onSelect func(filePath string, fileType int)) {
	dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err error) {
		if err != nil {
			return
		}
		onSelect(closer.URI().Path(), File)
	}, w)
}

func chooseSaveFile(w fyne.Window, onSelect func(filePath string, fileType int)) {
	dialog.ShowFileSave(func(closer fyne.URIWriteCloser, err error) {
		onSelect(closer.URI().Path(), File)
	}, w)
}

const (
	Dir  = 0
	File = 1
)

type Items struct {
	Path string
	Type int
}

func run() {
	a := app.New()
	w := a.NewWindow("Zip Pro X")

	var items []Items

	//hello := widget.NewLabel("Please Select File or Directory")

	list := widget.NewList(
		func() int {
			return len(items)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			//o.(*widget.Label).SetText(items[i].Path)
			//t.Refresh()
			log.Println("List Update")
			t := o.(*widget.Label)
			t.Text = items[i].Path
			t.Refresh()
		},
	)

	passwordEntry := widget.NewEntry()
	confirmPasswordForm := []*widget.FormItem{
		{
			Text:   "Password",
			Widget: passwordEntry,
		},
	}

	zipFilesBtn := widget.NewButtonWithIcon("Zip Files", theme.DocumentIcon(), func() {
		chooseSaveFile(w, func(filePath string, fileType int) {

			log.Println("Form filePath:", filePath)

			err := runZipFiles(filePath, items, "")
			if err != nil {
				dialog.NewError(err, w)
				fmt.Println("Zip file error.")
			} else {
				dialog.ShowInformation("Success", "Zip file created successfully.\nPath: "+filePath, w)
				fmt.Println("Zip file created successfully.")
			}

			passwordEntry.Text = ""
			items = []Items{}
			list.Refresh()

		})
	})

	zipFilesWithPasswordBtn := widget.NewButtonWithIcon("Zip Files With Password", theme.DocumentIcon(), func() {
		chooseSaveFile(w, func(filePath string, fileType int) {

			log.Println("Form filePath:", filePath)

			dialog.ShowForm("Enter Password", "OK", "Cancel", confirmPasswordForm, func(s bool) {
				log.Println("Form submitted:", passwordEntry.Text)
				err := runZipFiles(filePath, items, passwordEntry.Text)
				if err != nil {
					dialog.NewError(err, w)
					fmt.Println("Zip file error.")
				} else {
					dialog.ShowInformation("Success", "Zip file created successfully.\nPath: "+filePath, w)
					fmt.Println("Zip file created successfully.")
				}
				passwordEntry.Text = ""
				items = []Items{}
				list.Refresh()
			}, w)
		})
	})

	addFilesBtn := widget.NewButtonWithIcon("Add File", theme.ContentAddIcon(), func() {
		chooseFiles(w, func(filePath string, fileType int) {
			items = append(items, Items{
				Path: filePath,
				Type: fileType,
			})

			list.Resize(fyne.NewSize(500, float32(100*len(items))))
			list.Refresh()

			zipFilesBtn.Enable()
			zipFilesWithPasswordBtn.Enable()
		})

	})
	addDirectoryBtn := widget.NewButtonWithIcon("Add Directory", theme.ContentAddIcon(), func() {
		chooseDirectory(w, func(filePath string, fileType int) {
			items = append(items, Items{
				Path: filePath,
				Type: fileType,
			})

			list.Resize(fyne.NewSize(500, float32(100*len(items))))
			list.Refresh()

			zipFilesBtn.Enable()
			zipFilesWithPasswordBtn.Enable()
		})
	})

	zipFilesBtn.Disable()
	zipFilesWithPasswordBtn.Disable()

	w.SetContent(container.NewVBox(

		//hello,
		addFilesBtn,

		addDirectoryBtn,

		zipFilesBtn,

		zipFilesWithPasswordBtn,

		list,
	))

	w.Resize(fyne.NewSize(700, 700))
	w.ShowAndRun()
}

func runZipFiles(zipFilePath string, items []Items, password string) error {
	var list []string

	for _, item := range items {
		list = append(list, item.Path)
	}

	if password != "" {
		err := gozip.Zip(zipFilePath, list, password)
		if err != nil {
			return err
		}
	} else {
		err := gozip.Zip(zipFilePath, list)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	run()
}
