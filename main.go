package main

import (
	"fbeInstaller/icon"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	"github.com/getlantern/systray"
)

func init() {
	log.SetOutput(os.Stdout)
}

var installed []string

var (
	user32     = syscall.NewLazyDLL("user32.dll")
	messageBox = user32.NewProc("MessageBoxW")
	MB_OK      = 0x00000000
)

func msgBox(text, title string) {
	textPtr, _ := syscall.UTF16PtrFromString(text)
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	messageBox.Call(0, uintptr(unsafe.Pointer(textPtr)), uintptr(unsafe.Pointer(titlePtr)), uintptr(MB_OK))
}

func onReady() {
	systray.SetIcon(icon.Data)
	go func() {
		for {
			processFiles()
			time.Sleep(2 * time.Second)
		}
	}()

	systray.SetTitle("fbeInstaller")
	systray.SetTooltip("fbeInstaller")
	systray.SetIcon(icon.Data)
	mRemoveFiles := systray.AddMenuItem("Remove Installed Files", "remove installed .fbe")
	mQuit := systray.AddMenuItem("quit", "quit the application")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go func() {
		for {
			<-mRemoveFiles.ClickedCh
			if len(installed) > 0 {
				removeInstalledFiles()
			}
		}
	}()
}

func main() {
	systray.Run(onReady, onExit)
}

func onExit() {
	log.Println("*** fbeInstaller stopped")
}

func removeInstalledFiles() {
	for _, fileName := range installed {
		shipFilePath := filepath.Join(os.Getenv("APPDATA"), "Starbase", "ssc", "autosave", "ship_blueprints", fileName)
		err := os.Remove(shipFilePath)
		if err != nil {
			log.Printf("unable to remove file %s: %v", fileName, err)
		} else {
			log.Println("removed installed file:", fileName)
		}
	}
	installed = []string{} // clear installed files list

	// Show MessageBox
	msgBox(".fbe installed this session were removed", "success")
}

func processFiles() {
	// get files with .fbe extension in %userprofile%\Downloads folder
	fbes := filepath.Join(os.Getenv("USERPROFILE"), "Downloads")
	files, err := filepath.Glob(filepath.Join(fbes, "*.fbe"))
	if err != nil {
		log.Fatal(err)
	}

	// if there are any fbes in the downloads folder..
	if len(files) > 0 {
		// get newest fbe
		fbe := files[0]
		log.Println("hopefully moving:", fbe)

		// get files with .fbe extension in %appdata%\Starbase\ssc\autosave\ship_blueprints\ folder
		files, err = filepath.Glob(filepath.Join(os.Getenv("APPDATA"), "Starbase", "ssc", "autosave", "ship_blueprints", "*.fbe"))
		if err != nil {
			log.Fatal(err)
		}

		// find highest number file, format is ship_%d.fbe
		var highestNumber int
		for _, file := range files {
			filename := filepath.Base(file)
			var number int
			_, err := fmt.Sscanf(filename, "ship_%d.fbe", &number)
			if err == nil && number > highestNumber {
				highestNumber = number
			}
		}
		// increment by 1
		highestNumber++
		fmt.Println(highestNumber)

		// copy fbe to %appdata%\Starbase\ssc\autosave\ship_blueprints\ship_%d.fbe
		newPath := filepath.Join(os.Getenv("APPDATA"), "Starbase", "ssc", "autosave", "ship_blueprints", fmt.Sprintf("ship_%d.fbe", highestNumber))
		err = os.Rename(fbe, newPath)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("maybe finished moving to:", newPath)

		// append filename to installed
		installed = append(installed, filepath.Base(newPath))
	}
}
