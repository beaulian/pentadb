// Contains the interface and implementation of Log

/* BSD 3-Clause License

Copyright (c) 2017, Guan Jiawen, Li Lundong
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the copyright holder nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package log

import (
	"os"
	"io"
	"fmt"
	"time"
	"sync"
	"bytes"
	"runtime"
	"strconv"
)

// log level
const (
	// integer to control what color is printed
	// each is on behalf of a kind of level, making
	// the text that will be printed has different color
	InfoColor uint32 = 37               // for info level, color is white
	WarningColor uint32 = 33            // for warning level, color is yellow
	ErrorColor uint32 = 31              // for error level, color is red
)

// These flags define which text to prefix to each log entry generated by the Log.
const (
	// Bits or'ed together to control what's printed.
	// There is no control over the order they appear (the order listed
	// here) or the format they present (as described in the comments).
	// is specified.
	// For example, flags Ldate | Ltime produce,
	//	2009/01/23 01:23:23 message
	// while flags Ldate | Ltime | Lshortfile produce,
	//	2009/01/23 01:23:23 d.go:23: message
	Ldate         = 1         // the date in the local time zone: 2009/01/23
	Ltime         = 2         // the time in the local time zone: 01:23:23
	Lshortfile    = 4         // final file name element and line number: d.go:23.
)

type Log struct {
	out io.Writer          // destination for output
	flag int               // properties
	mu *sync.Mutex         // used for synchronization
	// color use a special way to store three kinds of color
	// for a 32-bit integer, color use the first 8-bit as info color
	// the subsequent 8-bit as warning color, and then the subsequent 8-bit
	// as error color
	color uint32           // define three colors for info, warning and error level
	colorTemplate string   // define how to print colored text
}


func NewLog(out io.Writer, flag int) *Log {
	color := InfoColor <<  24 + WarningColor << 16 + ErrorColor << 8
	return &Log {
		out:            out,
		flag:           flag,
		mu:             new(sync.Mutex),
		colorTemplate:  "\x1b[1;%dm%s\x1b[0m",
		color:          color,
	}
}

func (l *Log) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.flag = flag
}

// set color template to replace the old
func (l *Log) SetColorTemplate(ct string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.colorTemplate = ct
}

// set color for each level
func (l *Log) SetInfoColor(color uint32) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.color = l.color & 0x00ffffff & (color << 24)
}

func (l *Log) SetWarningColor(color uint32) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.color = l.color & 0xff00ffff & (color << 16)
}

func (l *Log) SetErrorColor(color uint32) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.color = l.color & 0xffff00ff & (color << 8)
}

func (l *Log) wrapper(text string) string {
	var buf bytes.Buffer
	if l.flag & (Ldate | Ltime) != 0 {
		t := time.Now()
		if l.flag & Ldate != 0 {
			year, month, day := t.Date()
			buf.WriteString(strconv.Itoa(year))
			buf.WriteString("/")
			buf.WriteString(strconv.Itoa(int(month)))
			buf.WriteString("/")
			buf.WriteString(strconv.Itoa(day))
		}
		if (l.flag & Ltime) != 0 {
			buf.WriteString(" ")
			hour, min, sec := t.Clock()
			buf.WriteString(strconv.Itoa(hour))
			buf.WriteString(":")
			buf.WriteString(strconv.Itoa(min))
			buf.WriteString(":")
			buf.WriteString(strconv.Itoa(sec))
		}
		buf.WriteString(" ")
	}
	if l.flag & Lshortfile != 0 {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		buf.WriteString(file)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(line))
		buf.WriteString(": ")
	}
	buf.WriteString(text)

	return buf.String()
}

func (l *Log) Info(format string, v ...interface{}) {
	text := fmt.Sprintf(l.colorTemplate, l.color >> 24, l.wrapper(fmt.Sprintf(format, v...)))
	fmt.Println(text)
}

func (l *Log) Warning(format string, v ...interface{}) {
	text := fmt.Sprintf(l.colorTemplate, (l.color & 0x00ff0000) >> 16, l.wrapper(fmt.Sprintf(format, v...)))
	fmt.Println(text)
	panic(text)
}

func (l *Log) Error(format string, v ...interface{}) {
	text := fmt.Sprintf(l.colorTemplate, (l.color & 0x0000ff00) >> 8, l.wrapper(fmt.Sprintf(format, v...)))
	fmt.Println(text)
	os.Exit(1)
}