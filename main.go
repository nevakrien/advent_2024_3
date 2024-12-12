package main

import (
    "fmt"
    "os"
    "bufio"
    "io"
    "errors"
)

type PeekReader struct {
	reader *bufio.Reader
	peeked rune
	peekedValid bool
}

// NewPeekReader creates a new PeekReader wrapping the given io.Reader.
func NewPeekReader(r io.Reader) *PeekReader {
	return &PeekReader{
		reader: bufio.NewReader(r),
	}
}

// Peek returns the next character without consuming it.
func (pr *PeekReader) Peek() (rune, error) {
	if pr.peekedValid {
		return pr.peeked, nil
	}
	r, _, err := pr.reader.ReadRune()
	if err != nil {
		return -1, err
	}
	pr.peeked = r
	pr.peekedValid = true
	return r, nil
}

// Consume reads and consumes the next character from the stream.
func (pr *PeekReader) Consume() (rune, error) {
	if pr.peekedValid {
		pr.peekedValid = false
		return pr.peeked, nil
	}
	r, _, err := pr.reader.ReadRune()
	if err != nil {
		return -1, err
	}
	return r, nil
}

var FailedMatch = errors.New("did not match the pattern")

func isAsciiDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func ParseNum(pr *PeekReader) (int, error) {
	c,err := pr.Peek()
	if err!= nil {
		return 0,err
	}
	if !isAsciiDigit(c){
		return 0,FailedMatch
	}

	_,_ = pr.Consume()
	num := int(c-'0')

	for {
		c,err := pr.Peek()
		if err==io.EOF {
			return num,nil
		}
		if err!= nil {
			return num,err
		}
		if!isAsciiDigit(c){
			return num,nil
		}
		_,_ = pr.Consume()
		num = 10*num+int(c-'0')
	}	
}

func main() {
	// Open the file
    file, err := os.Open("input.txt")
    if err != nil {
        fmt.Println("Error opening file:", err)
        os.Exit(1)

    }
    defer file.Close() // Ensure the file is closed when the function exits
    reader := NewPeekReader(file)

    for {
    	num,err := ParseNum(reader)
    	if err==nil {
    		fmt.Println("\nnum(",num,")")
    		continue
    	}
    	c,err := reader.Consume()
    	if(err == io.EOF){
    		return
    	} else if err!=nil {
    		fmt.Println("Error reading:", err)
        	os.Exit(1)
    	}
    	fmt.Print(string(c))

    }

}