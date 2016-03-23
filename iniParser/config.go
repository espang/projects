package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

type section struct {
	name   string
	values map[string]string
}

func (s *section) add(key, value string) {
	s.values[key] = value
}

type config struct {
	sections map[string]*section
}

func (c *config) addSection(s *section) {
	c.sections[s.name] = s
}

func NewConfig(f *os.File) (*config, error) {

	c := &config{
		sections: map[string]*section{},
	}

	reader := bufio.NewReader(f)
	eof := false
	var currentSeciton = &section{
		name:   "default",
		values: map[string]string{},
	}

	for nbr := 1; !eof; nbr++ {

		var ln []byte
		for {
			line, isPrefix, err := reader.ReadLine()
			eof = err == io.EOF
			if eof {
				//line end
				break
			} else if err != nil {
				log.Print("error : ", err)
				return &config{}, err
			}
			ln = append(ln, line...)
			if isPrefix == false {
				break
			}
		}

		if IsEmptyOrComment(ln) {
			// nothing to do here
			continue
		}

		if name, ok := GetSectionFromLine(ln); ok {
			if currentSeciton != nil {
				c.addSection(currentSeciton)
			}
			currentSeciton = &section{
				name:   name,
				values: map[string]string{},
			}
			continue
		}

		key, value, err := ParseKeyValuePair(nbr, ln)
		if err != nil {
			return &config{}, err
		}
		currentSeciton.add(key, value)
	}
	c.addSection(currentSeciton)

	return c, nil
}
