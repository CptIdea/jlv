package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

var COLORS = map[string]string{
	"info":    "\033[37m",
	"warning": "\033[36m",
	"fatal":   "\033[31m",
	"panic":   "\033[31m",
	"error":   "\033[33m",
}

func main() {
	showInfo := flag.Bool("i", false, "Show info logs")
	showWarn := flag.Bool("w", false, "Show warning logs")
	showErr := flag.Bool("e", false, "Show error logs")
	showFatal := flag.Bool("f", false, "Show fatal logs")
	showPanic := flag.Bool("p", false, "Show panic logs")

	showTime := flag.Bool("t", false, "Show only log time")

	showStats := flag.Bool("s", false, "Show log statistic")

	InvertShow := flag.Bool("I", false, "Inverted all show")

	fileAddr := flag.String("F", "", "file to view log") //Адрес файла который будем читать
	TCPAddr := flag.String("T", "", "TCP source of log")

	count := flag.Int("c", 0, "Count of viewing logs")

	cont := flag.Bool("C", false, "continuous reading")

	flag.Parse()
	if !*showErr && !*showFatal && !*showPanic && !*showWarn && !*showInfo {
		*showInfo = true
		*showWarn = true
		*showErr = true
		*showFatal = true
		*showPanic = true
	}
	if *InvertShow {
		*showInfo = !*showInfo
		*showWarn = !*showWarn
		*showErr = !*showErr
		*showFatal = !*showFatal
		*showPanic = !*showPanic
	}

	if *fileAddr == "" && *TCPAddr == "" {
		log.Fatal("Source not attached")
		return
	}
	var Src interface {
		io.Reader
		io.Writer
	}
	var err error
	var data []byte

	switch {
	case *fileAddr != "":
		Src, err = os.Open(*fileAddr)
		if err != nil {
			log.Fatal("Src open fail:", err)
			return
		}
		data, err = ioutil.ReadAll(Src)
		if err != nil {
			log.Fatal("Src read fail:", err)
			return
		}

	case *TCPAddr != "":
		Src, err = net.Dial("tcp", *TCPAddr)
		if err != nil {
			log.Fatal("Src open fail:", err)
			return
		}

		_, err = fmt.Fprint(Src, "getall\n")
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(Src)

		if err != nil {
			log.Fatal("Src read fail:", err)
			return
		}
	}

	if err != nil {
		log.Fatal("Read file fail:", err)
	}
	rawJSONs := strings.Split(string(data), "\n")
	var infoC, warnC, erroC, fataC, paniC int

	for i, n := range rawJSONs {
		if *count != 0 && i >= *count {
			break
		}
		CurOut := make(map[string]string)
		json.Unmarshal([]byte(n), &CurOut)
		lvl, ok := CurOut["level"]
		if !ok {
			continue
		}
		switch lvl {
		case "info":
			infoC++
			if !*showInfo {
				continue
			}
		case "error":
			erroC++
			if !*showErr {
				continue
			}
		case "fatal":
			fataC++
			if !*showFatal {
				continue
			}
		case "warning":
			warnC++
			if !*showWarn {
				continue
			}
		case "panic":
			paniC++
			if !*showPanic {
				continue
			}

		}
		msg, ok := CurOut["msg"]
		if !ok {

		}
		tm, ok := CurOut["time"]
		if !ok {
			continue
		}
		var fields string
		for s, cur := range CurOut {
			if s != "level" && s != "msg" && s != "time" && s != "prefix" {
				fields += s + "=" + cur + " "
			}
		}
		prefix, ok := CurOut["prefix"]
		if !ok {
			prefix = "ALL"
		}
		color := COLORS[lvl]
		lvl = string([]byte(lvl)[:4])
		tm = strings.Replace(tm, "T", " ", 1)
		tm = string([]byte(tm)[:19])
		if *showTime {
			tm = strings.Split(tm, " ")[1]
		}

		fmt.Printf("%s%s\033[0m\t%s\t[%s]%s\t%s\n", color, strings.ToTitle(lvl), tm, prefix, msg, fields)

	}
	if *showStats {
		fmt.Printf("\n\033[38;05;68mSTATS: \n\033[0m INFO:\t WARN:\t ERRO:\t FATA:\t PANI:\t  ALL:\n%6d\t%6d\t%6d\t%6d\t%6d\t%6d\n", infoC, warnC, erroC, fataC, paniC, infoC+warnC+erroC+fataC+paniC)
	}
	if *cont {
		if *TCPAddr != "" {
			Src, err = net.Dial("tcp", *TCPAddr)
			if err != nil {
				log.Fatal("Src open fail:", err)
				return
			}
			fmt.Fprint(Src, "connect\n")
			fmt.Println("\033[01mCONNECTED\033[0m")
		}
		for {
			time.Sleep(time.Millisecond * 10)
			var newdata []byte
			if *TCPAddr != "" {
				newdata, err = bufio.NewReader(Src).ReadBytes('\n')
			} else {
				newdata, err = ioutil.ReadAll(Src)
			}
			if err != nil {
				log.Fatal("jlv:", err)
			}
			if len(newdata) > 0 {
				rawJSONs := strings.Split(string(newdata), "\n")
				for i, n := range rawJSONs {
					if *count != 0 && i >= *count {
						break
					}
					CurOut := make(map[string]string)
					json.Unmarshal([]byte(n), &CurOut)
					lvl, ok := CurOut["level"]
					if !ok {
						continue
					}
					switch lvl {
					case "info":
						if !*showInfo {
							continue
						}
					case "error":
						if !*showErr {
							continue
						}
					case "fatal":
						if !*showFatal {
							continue
						}
					case "warning":
						if !*showWarn {
							continue
						}
					case "panic":
						if !*showPanic {
							continue
						}

					}
					msg, ok := CurOut["msg"]
					if !ok {

					}
					tm, ok := CurOut["time"]
					if !ok {
						continue
					}
					var fields string
					for s, cur := range CurOut {
						if s != "level" && s != "msg" && s != "time" && s != "prefix" {
							fields += s + "=" + cur + " "
						}
					}
					prefix, ok := CurOut["prefix"]
					if !ok {
						prefix = "ALL"
					}
					color := COLORS[lvl]
					lvl = string([]byte(lvl)[:4])
					tm = strings.Replace(tm, "T", " ", 1)
					tm = string([]byte(tm)[:19])
					if *showTime {
						tm = strings.Split(tm, " ")[1]
					}

					fmt.Printf("%s%s\033[0m\t%s\t[%s]%s\t%s\n", color, strings.ToTitle(lvl), tm, prefix, msg, fields)

				}
			}
		}
	}

}
