package main

import (
	"io"
	"os"
	"os/exec"
)

var (
	input    io.ReadCloser
	output   io.WriteCloser
	readBuf  *byteBuf
	writeBuf *byteBuf
)

const (
	maxReadBufSize  int = 1024
	maxWriteBufSize int = 1024
)

func main() {
	args := os.Args[1:]
	if err := parseFlags(&args); err != nil {
		os.Stderr.Write([]byte(err.Error() + "\n"))
		return
	}
	if *showHelpMsg {
		printUsage()
		return
	}
	if err := initPipes(); err != nil {
		os.Stderr.Write([]byte(err.Error() + "\n"))
		return
	}
	readBuf = newByteBuf(maxReadBufSize)
	writeBuf = newByteBuf(maxWriteBufSize)
	readingPage := 1
	for readingPage < *startPageN {
		if err := readOnePage(false); err != nil {
			os.Stderr.Write([]byte(err.Error() + "\n"))
			return
		} else {
			readingPage += 1
		}
	}
	for readingPage <= *endPageN {
		if err := readOnePage(true); err != nil {
			os.Stderr.Write([]byte(err.Error() + "\n"))
			return
		} else {
			readingPage += 1
		}
	}
	input.Close()
	output.Close()
}

func readOnePage(shouldWrite bool) error {
	if formFeed != nil && *formFeed {
		if err := readUntil('\f', shouldWrite); err != nil {
			return err
		}
	} else {
		if err := readForLine(*linesPerPage, shouldWrite); err != nil {
			return err
		}
	}
	return nil
}

func initPipes() error {
	var err error
	if filename != nil {
		if input, err = os.Open(*filename); err != nil {
			return err
		}
	} else {
		input = os.Stdin
	}
	if toDestination != nil {
		if output, err = exec.Command("lp", "-d"+*toDestination).StdinPipe(); err != nil {
			return err
		}
	} else {
		output = os.Stdout
	}
	return nil
}

func readForLine(lineN int, shouldWrite bool) error {
	for lineN > 0 {
		if readBuf.IsEmpty() && !readBuf.Read(input) {
			break
		}
		nextByte := readBuf.ReadNextByte()
		if nextByte == '\n' {
			lineN--
		}
		writeBuf.WriteNextByte(nextByte)
		if writeBuf.IsFull() {
			if shouldWrite {
				if err := writeBuf.Write(output); err != nil {
					return err
				}
			} else {
				writeBuf.Clear()
			}
		}
	}
	if !writeBuf.IsEmpty() {
		if shouldWrite {
			if err := writeBuf.Write(output); err != nil {
				return err
			}
		} else {
			writeBuf.Clear()
		}
	}
	return nil
}

func readUntil(delimiter byte, shouldWrite bool) error {
	for {
		if readBuf.IsEmpty() && !readBuf.Read(input) {
			break
		}
		nextByte := readBuf.ReadNextByte()
		if nextByte == delimiter {
			break
		}
		writeBuf.WriteNextByte(nextByte)
		if writeBuf.IsFull() {
			if shouldWrite {
				if err := writeBuf.Write(output); err != nil {
					return err
				}
			} else {
				writeBuf.Clear()
			}
		}
	}
	if !writeBuf.IsEmpty() {
		if shouldWrite {
			if err := writeBuf.Write(output); err != nil {
				return err
			}
		} else {
			writeBuf.Clear()
		}
	}
	return nil
}
