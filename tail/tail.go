package tail

import (
	"bufio"
	"io"
	"os"
	"strings"
	"time"
)

type Tail struct {
	Filename string
	Lines    chan string

	file   *os.File
	reader *bufio.Reader
	stop   bool
}

func TailFile(filename string) (tail *Tail, err error) {
	tail = &Tail{
		Filename: filename,
		Lines:    make(chan string, 500),
	}

	tail.stop = false

	go tail.processFile()
	return
}

func (tail *Tail) processFile() {
	offset := int64(0)
	var err error

	tail.file, err = os.Open(tail.Filename)
	check(err)
	defer tail.file.Close()
	tail.reader = bufio.NewReader(tail.file)

	tail.file.Seek(offset, io.SeekStart)
	for !tail.stop {
		offset, err := tail.getOffset()
		if err != nil {
			tail.stop = true
			return
		}
		line, err := tail.readLine()
		if err == nil {
			tail.Lines <- line
		} else {
			tail.waitForData()
			fileEnd, err := tail.file.Seek(0, io.SeekEnd)
			if err != nil || fileEnd < offset {
				// The file has changed since we opened it.
				// Start looking at the file again from the beginning.
				offset = 0
			}
			tail.file.Seek(offset, io.SeekStart)
			tail.reader.Reset(tail.file)
		}
	}
}

func (tail *Tail) getOffset() (offset int64, err error) {
	offset, err = tail.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return
	}
	offset -= int64(tail.reader.Buffered())
	return
}

func (tail *Tail) readLine() (string, error) {
	line, err := tail.reader.ReadString('\n')
	if err != nil {
		line = strings.TrimRight(line, "\n")
	}
	return line, err
}

func (tail *Tail) waitForData() {
	time.Sleep(1000 * time.Millisecond)
}

func (tail *Tail) Close() {
	tail.stop = true
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
