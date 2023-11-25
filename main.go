package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

const (
	fw     = "fw"  // find word command
	log_   = "-l"  // add found word to log file
	remove = "-rm" // remove found word/word in a file
	help   = "-help"

	// shortcuts
	desk = "C:\\Users\\User\\OneDrive\\Masaüstü\\"
)

type WordInfo struct {
	word  string
	count int
}

// Operations Words
type Operations interface {
	Find(word string) (WordInfo, string)
	RemoveWordS() error
	LogFoundWordS() error
}

type OperationWords struct {
	// file     *os.File
	wordInfo WordInfo
	result   string
}

func NewOperationWords() *OperationWords {
	return &OperationWords{}
}

func (ow *OperationWords) Find(fileName string, word string) (WordInfo, string) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	// ow.file = file
	if err != nil {
		log.Fatalf(color.New(color.FgRed).Sprintf("[ Dosya okunurken bir hata oluştu ] %s", err))
	}
	defer file.Close()

	// read buffer file content
	reader := bufio.NewReader(file)
	bufferedContent := make([]byte, 1000000)
	reader.Read(bufferedContent)

	contentWords := strings.Split(string(bufferedContent), "\r\n")
	wordCount := 0
	var result string
	for _, v := range contentWords {
		result += "\n"
		for _, v2 := range strings.Split(v, " ") {
			if v2 == word {
				wordCount++
				result += color.New(color.FgRed).Sprint(" " + v2)
			} else {
				result += " " + v2
				continue
			}
		}
	}
	ow.wordInfo = WordInfo{word: word, count: wordCount}
	ow.result = result
	return ow.wordInfo, ow.result
}

// Todo Fix Log func
func (ow *OperationWords) Log() {
	file, err := os.OpenFile("C:\\Users\\User\\OneDrive\\Masaüstü\\log.txt", os.O_CREATE, os.ModeAppend)
	if err != nil {
		log.Fatalf(color.New(color.FgRed).Sprintf("[ Dosya okunurken bir hata oluştu ] %s", err))
	}
	defer file.Close()
	content := ow.wordInfo.word + " " + fmt.Sprint(ow.wordInfo.count) + ow.result
	if _, err := file.WriteString(content); err != nil {
		log.Fatalf(color.New(color.FgRed).Sprintf("[ Dosya okunurken bir hata oluştu ] %s", err))
	}
	color.Cyan("log dosyası oluşturuldu")
}

// Handle Commands

type Handle interface {
	Parse() error
	Start() error
	Help() string
}

type HandleCommand struct {
	command []string
	log_    bool
	// remove  bool
	result []string
	isHelp bool
}

func NewHandleCommand(command []string) *HandleCommand {
	return &HandleCommand{
		command: command,
		log_:    false,
		// remove:  false,
		isHelp: false,
	}
}

//en az fw word file_name.txt 3 tane

func (hc *HandleCommand) Parse() error {
	if len(hc.command) < 3 {
		if hc.command[0] == fw && len(hc.command) == 2 && hc.command[1] == help {
			hc.Help()
			hc.isHelp = true
			return nil
		}
		return fmt.Errorf(color.New(color.FgRed).Sprintf("[ Eksik Komut girildi ] %v. -help yazarak komutlara bakabilirsiniz.", hc.command))
	}
	prefix := hc.command[0]
	if prefix != fw {
		hc.Help()
		fmt.Println("")
		return fmt.Errorf(color.New(color.FgRed).Sprintf("[ Yanlış prefix ] %v beklendi %v bulundu. -help yazarak komutlara bakabilirsiniz.", fw, prefix))
	}

	for _, v := range hc.command {
		switch v {
		case log_:
			{
				hc.log_ = true
			}
		case help:
			{
				return fmt.Errorf(color.New(color.FgRed).Sprintf("[ Yanlış komut iskeleti ] %v. -help yazarak komutlara bakabilirsiniz.", hc.command))

			}
		// case remove:
		// 	{
		// 		hc.remove = true
		// 	}
		default:
			{
				if v != fw {
					hc.result = append(hc.result, v)
				}
			}
		}
	}

	return nil
}

// Todo = log dosyasında renkli yazmıyor

func (hc *HandleCommand) Start() {
	if hc.isHelp {
		return
	}
	word := strings.Join(hc.result[:len(hc.result)-1], " ")
	file := hc.result[len(hc.result)-1]
	if strings.Split(file, `\`)[0] == "desk" {
		file = desk + strings.Split(file, `\`)[1]
	}
	ow := NewOperationWords()
	wordInfo, result := ow.Find(file, word)
	if hc.log_ {
		ow.Log()
	}

	fmt.Println(result, "\r\n", color.New(color.FgGreen).Sprint("\r\n", wordInfo.count, " tane", ` "`, wordInfo.word, `"`, " bulundu."))
}

func (hc *HandleCommand) Help() {
	commands := [][]string{
		[]string{fw, "Uygulamanın çalışmasını sağlar."},
		[]string{log_, "Masaüstüne log.txt olarak çıktı verir."},
		[]string{remove, "Bulunan kelimeleri dosyadan kaldırır."},
	}

	table1 := tablewriter.NewWriter(os.Stdout)
	table1.SetHeader([]string{"Komut", "Komut açıklaması"})
	table1.SetBorder(false)
	table1.AppendBulk(commands)
	table1.Render()

	fmt.Println("")

	shortcuts := [][]string{
		[]string{`deks\`, "Masaüstü dizini için bir kısayoldur. (desk\\a.txt)"},
	}

	table2 := tablewriter.NewWriter(os.Stdout)
	table2.SetHeader([]string{"Komut", "Komut açıklaması"})
	table2.SetBorder(false)
	table2.SetFooter([]string{"", "", fmt.Sprint(time.Now().Date())})
	table2.AppendBulk(shortcuts)
	table2.Render()
}

func main() {
	handleCommand := NewHandleCommand(os.Args[1:])
	if err := handleCommand.Parse(); err != nil {
		log.Fatal(err)
	}
	handleCommand.Start()
}
