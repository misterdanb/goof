package main

import (
	"path/filepath"
	"runtime"
	gc "code.google.com/p/goncurses"
	"strings"
	"github.com/jessevdk/go-flags"
	"os"
	"net"
	"os/signal"
	"sync"
	"net/http"
	"strconv"
	"bufio"
)

var parser *flags.Parser

var cmdFlags struct {
	Help bool `short:"h" long:"help" description:"Show this help message"`
	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Self bool `short:"s" long:"self" description:"Goof the goof"`
	Port uint `short:"p" long:"port" description:"Port to serve on"`
    Count uint `short:"c" long:"count" description:"Amount of possible downloads"`
    EarWigglingSpeed uint `short:"e" long:"ear-wiggling-speed" description:"Speed of ear wiggling"`
}

var wg sync.WaitGroup

var cmdArgs []string

var self bool
var port uint
var counter uint
var path string
var customAnimationSpeed uint

var localIPs []string

var initErr error
var listener net.Listener
var running bool
var currentSpeed int64

var stdscr gc.Window

var win gc.Window
var rows, cols int

var frame int
var framesCount int

func init() {
	cmdFlags.Verbose = false
	cmdFlags.Self = false
	cmdFlags.Port = 1337
	cmdFlags.Count = 1
	cmdFlags.EarWigglingSpeed = 20
	
	parser = flags.NewParser(&cmdFlags, flags.PrintErrors | flags.PassDoubleDash)
	parser.Usage = "[OPTIONS] <filename>"
	
	if cmdArgs, initErr = parser.Parse(); initErr != nil {
		return
	} else {
		self = cmdFlags.Self
		counter = cmdFlags.Count
		port = cmdFlags.Port
		
		if len(cmdArgs) == 1 {
			path = cmdArgs[0]
		}
		
		customAnimationSpeed = cmdFlags.EarWigglingSpeed
	}
}

func main() {
	var (
		err error
	)
	
	if initErr != nil {
		return
	}
	
	if cmdFlags.Help {
		Haaaaalp()
		return
	}
	
	if len(cmdArgs) != 1 && !self {
		Haaaaalp()
		return
	}
	
	if self {
		path, err = GetSelfPath()
		path += "/goof"
		
		if err != nil {
			print("Something went wrong. RUN!!!\r\n")
		}
	}
	
	_, err = os.Stat(path)
	
	if os.IsNotExist(err)  {
		print("Yo, this file doesn't exist.\r\n")
		return
	}
	
	if customAnimationSpeed < 20 {
		print("Yo, the ear wiggling speed must be 20 or higher- It's because of goof's coolness.\r\n")
		return
	}
	
	localIPs, _ = GetIPs()
	
	wg.Add(1)
	
	// init ncurses and end it, if main returns
    _, err = gc.Init()
	
	// handle ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			EndProgram()
		}
	}()
	
	if err != nil {
		print("Could not initialize fancy ncurses graphics. Sorry, have to go now...\r\n")
	}
	
	// everything's fine, let's go
	running = true
	
	SetupServingInformation()
	SetStartServingInformation()
	
	go ListenAndServeFiles(":" + strconv.FormatInt(int64(port), 10))
	
	wg.Wait()
	
	gc.End()
	
	return
}

func EndProgram() {
	listener.Close()
	running = false
	gc.End()
}

func GetSelfPath() (string, error) {
	var (
		dir, p string
		err error
	)
	
	if runtime.GOOS == "linux" {
		pid := os.Getpid()
		link := "/proc/" + strconv.Itoa(pid) + "/exe"
		p, err = os.Readlink(link)
	} else {
		print("Your OS (" + runtime.GOOS + ") is not supported. Use Linux.\r\n")
	}
	
	if err != nil {
		return "", err
	}
	
	dir = filepath.Dir(p)
	dir = strings.Replace(dir, "\\", "/", -1)
	
	return dir, nil
}

func GetIPs() ([]string, error) {
	var (
		Fehler error
		Kandidaten []string
		WeltnetzVerfahrensAdressenZumTesten []string
	)
	
	WeltnetzVerfahrensAdressenZumTesten = append(WeltnetzVerfahrensAdressenZumTesten, "192.0.2.0", "198.51.100.0", "203.0.113.0")
	
	for Index := 0; Index < len(WeltnetzVerfahrensAdressenZumTesten); Index++ {
		Verbindung, Fehler := net.Dial("udp", WeltnetzVerfahrensAdressenZumTesten[Index] + ":80")
		
		if Fehler == nil {
			EigeneAufgespalteneWeltnetzVerfahrensAdresse := strings.Split(Verbindung.LocalAddr().String(), ":")
			
			WeltnetzVerfahrensAdresseBereitsEnthalten := false
			
			for IndexZwei := 0; IndexZwei < len(Kandidaten); IndexZwei++ {
				if EigeneAufgespalteneWeltnetzVerfahrensAdresse[0] == Kandidaten[IndexZwei] {
					WeltnetzVerfahrensAdresseBereitsEnthalten = true
					break
				}
			}
			
			if !WeltnetzVerfahrensAdresseBereitsEnthalten {
				Kandidaten = append(Kandidaten, EigeneAufgespalteneWeltnetzVerfahrensAdresse[0])
			}
		}
	}
	
	return Kandidaten, Fehler
}

func Haaaaalp() {
	print("Goof - 20% cooler than woof!\n\n")
	print("(furthermore it's written in Go, which\r\nbtw. is a horrible language. Don't write\r\nprograms in Go...)\r\n\r\n")
	
	print("                   .\r\n")
	print("                 \\ | /\r\n")
	print("                 _\\|/_\r\n")
	print("               .' ' ' '.        ___\r\n")
	print("             _.|.--.--.|.___.--'___`-.\r\n")
	print("           .'.'||  |  ||`----'\"`   ``'`\r\n")
	print("         .'.'  ||()|()||\r\n")
	print(" .___..-'.'    /       \\\r\n")
	print(" `----'\"`     /   .-.   \\\r\n")
	print("             (.'.(___).'.)\r\n")
	print("              `.__.-.__.'\r\n")
	print("               |_|   |_|\r\n")
	print("                `.`-'.'\r\n")
	print("                  `\"`\r\n\r\n")
	
	parser.WriteHelp(os.Stdout)
}

func PrintEmptyLinesOnWindow(win gc.Window, n int) {
	for i := 0; i < n; i++ {
		win.Println()
	}
}

func PrintSpacesOnWindow(win gc.Window, n int) {
	for i := 0; i < n; i++ {
		win.Print(" ")
	}
}

func PrintCenteredOnWindow(win gc.Window, cols, lineLength int, line string) {
	PrintSpacesOnWindow(win, (cols - lineLength) / 2)
	win.Println(line)
}

func createWindow(h, w, y, x int) gc.Window {
	new, _ := gc.NewWindow(h, w, y, x)
	new.Refresh()
	
	return *new
}

func SetupServingInformation() {
	gc.Echo(false)
	gc.CBreak(true)
	gc.Cursor(0)
	
    stdscr.Keypad(true)

    rows, cols = stdscr.MaxYX()
	
	win = createWindow(rows, cols, 0, 0)
	
	frame = 0
	framesCount = 2
	
	return
}

func SetStartServingInformation() {
	win.Clear()
	
	PrintEmptyLinesOnWindow(win, (rows - 14) / 2)
	
	PrintCenteredOnWindow(win, cols, 42, "                     .                     ")
	PrintCenteredOnWindow(win, cols, 42, "                   \\ | /                   ")
	PrintCenteredOnWindow(win, cols, 42, "                   _\\|/_                   ")
	PrintCenteredOnWindow(win, cols, 42, "                 .' ' ' '.                 ")
	PrintCenteredOnWindow(win, cols, 42, "     ____________|.--.--.|____________     ")
	PrintCenteredOnWindow(win, cols, 42, "   <´____________||  |  ||____________`>    ")
	PrintCenteredOnWindow(win, cols, 42, "                 ||()|()||                 ")
	PrintCenteredOnWindow(win, cols, 42, "                 /       \\                 ")
	PrintCenteredOnWindow(win, cols, 42, "                /   .-.   \\                ")
	PrintCenteredOnWindow(win, cols, 42, "               (.'.(___).'.)               ")
	PrintCenteredOnWindow(win, cols, 42, "                `.__.-.__.'                ")
	PrintCenteredOnWindow(win, cols, 42, "                 |_|   |_|                 ")
	PrintCenteredOnWindow(win, cols, 42, "                  '.`-´.'                 ")
	PrintCenteredOnWindow(win, cols, 42, "                    `\"´                   ")
	
	splittedPath := strings.Split(path, "/")
	
	waitingString := "Waiting for clients."
	offeringString := "Offering: " + path
	wgetString := "wget http://" + localIPs[0] + ":" + strconv.FormatInt(int64(port), 10) + "/" + splittedPath[len(splittedPath) - 1]
	
	PrintCenteredOnWindow(win, cols, 0, "")
	PrintCenteredOnWindow(win, cols, len(waitingString), waitingString)
	PrintCenteredOnWindow(win, cols, len(offeringString), offeringString)
	PrintCenteredOnWindow(win, cols, 0, "")
	PrintCenteredOnWindow(win, cols, len(wgetString), wgetString)
	
	win.Refresh()
}

func UpdateServingInformation() {
	win.Clear()
	
	PrintEmptyLinesOnWindow(win, (rows - 14) / 2)
	
	switch (frame) {
		case 0:
			PrintCenteredOnWindow(win, cols, 42, "                     .                     ")
			PrintCenteredOnWindow(win, cols, 42, "                   \\ | /                   ")
			PrintCenteredOnWindow(win, cols, 42, "                   _\\|/_                   ")
			PrintCenteredOnWindow(win, cols, 42, "                 .' ' ' '.        ___      ")
			PrintCenteredOnWindow(win, cols, 42, "               _.|.--.--.|.___.--'___`-.   ")
			PrintCenteredOnWindow(win, cols, 42, "             .'.'||  |  ||`----'\"`   ``'`  ")
			PrintCenteredOnWindow(win, cols, 42, "           .'.'  ||()|()||                 ")
			PrintCenteredOnWindow(win, cols, 42, "   .___..-'.'    /       \\                 ")
			PrintCenteredOnWindow(win, cols, 42, "   `----'\"`     /   .-.   \\                ")
			PrintCenteredOnWindow(win, cols, 42, "               (.'.(___).'.)               ")
			PrintCenteredOnWindow(win, cols, 42, "                `.__.-.__.'                ")
			PrintCenteredOnWindow(win, cols, 42, "                 |_|   |_|                 ")
			PrintCenteredOnWindow(win, cols, 42, "                  '.`-´.'                 ")
			PrintCenteredOnWindow(win, cols, 42, "                    `\"´                   ")
		case 1:
			PrintCenteredOnWindow(win, cols, 42, "                     .                     ")
			PrintCenteredOnWindow(win, cols, 42, "                   \\ | /                   ")
			PrintCenteredOnWindow(win, cols, 42, "                   _\\|/_                   ")
			PrintCenteredOnWindow(win, cols, 42, "      ___        .' ' ' '.                 ")
			PrintCenteredOnWindow(win, cols, 42, "   .-´___'--.___.|.--.--.|._              ")
			PrintCenteredOnWindow(win, cols, 42, "  ´'´´   `\"'----´||  |  ||'.'.         ")
			PrintCenteredOnWindow(win, cols, 42, "                 ||()|()||  '.'.           ")
			PrintCenteredOnWindow(win, cols, 42, "                 /       \\    '.'-..___.   ")
			PrintCenteredOnWindow(win, cols, 42, "                /   .-.   \\     `\"'----´   ")
			PrintCenteredOnWindow(win, cols, 42, "               (.'.(___).'.)               ")
			PrintCenteredOnWindow(win, cols, 42, "                `.__.-.__.'                ")
			PrintCenteredOnWindow(win, cols, 42, "                 |_|   |_|                 ")
			PrintCenteredOnWindow(win, cols, 42, "                  '.`-´.'                 ")
			PrintCenteredOnWindow(win, cols, 42, "                    `\"´                   ")
	}
	
	servingString := "I'm serving with " + strconv.FormatInt(int64(customAnimationSpeed), 10) + "% more coolness than woof."
	
	PrintCenteredOnWindow(win, cols, 0, "")
	PrintCenteredOnWindow(win, cols, len(servingString), servingString)
	
	win.Refresh()
	
	frame++
	frame %= framesCount
}

func ListenAndServeFiles(addressPortPart string) { 
	var (
		err error
	)
	
	defer wg.Done()
	
	listener, err = net.Listen("tcp", "0.0.0.0" + addressPortPart)
	
	if err != nil {
		println("Error listening: ", err.Error())
		running = false
		
		return
	}
	
	http.HandleFunc("/", RequestHandler)
	http.Serve(listener, nil)
	
	if err != nil {
		println("Error serving: ", err.Error())
		running = false
		
		return
	}
	
	return
}

func RequestHandler(writer http.ResponseWriter, request *http.Request) {
	var (
		file *os.File
		fileInfo os.FileInfo
		reader *bufio.Reader
		buffer []byte
		bufferSize int
		readBytes int
		baseAnimationRefreshRate uint
		err error
	)
	
	counter--
	
	if file, err = os.Open(path); err != nil {
		return
	}
	
	defer file.Close()
	
	baseAnimationRefreshRate = (1000000 * 20) / customAnimationSpeed
	
	bufferSize = 2^16
	buffer = make([]byte, bufferSize)
	
	reader = bufio.NewReaderSize(file, bufferSize)
	writer.Header().Set("Content-Type", "application/octet-stream")
	fileInfo, err = file.Stat()
	writer.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	
	for err == nil {
		// make a timestamp
		//t0 := time.Now()
		
		for i := 0; i < int(baseAnimationRefreshRate); i++ {
			// read some bytes from file
			readBytes, err = reader.Read(buffer)
			
			// write as many bytes as read
			writer.Write(buffer[:readBytes])
			
			if err != nil {
				break
			}
		}
		
		UpdateServingInformation()
		
		/*// make another timestamp
		t1 := time.Now()
		
		// calculate elapsed time
		dt := t1.Sub(t0)
		
		// update speed
		currentSpeed += int64(readBytes) / dt.Nanoseconds()
		currentSpeed /= 2*/
	}
	
	//reader.WriteTo(writer)
	
	if counter <= 0 {
		EndProgram()
	}
	
	return
}
