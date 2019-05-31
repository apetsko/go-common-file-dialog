package main

import (
	cfd "github.com/harry1453/go-common-file-dialog"
	"log"
)

func main() {
	cfd.Initialize()
	defer cfd.UnInitialize()
	openDialog, err := cfd.NewOpenFileDialog(cfd.DialogConfig{
		Title: "Open Text File",
		Role:  "OpenTextExample",
		FileFilter: []cfd.FileFilter{
			{
				DisplayName: "Text Files (*.txt)",
				Pattern:     "*.txt",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := openDialog.Show(); err != nil {
		log.Fatal(err)
	}
	result, err := openDialog.GetResult()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Chosen file: %s\n", result)
}