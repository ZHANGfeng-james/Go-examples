package gobytes

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"unicode"
)

var (
	bComment = []byte{'#'}
	bEmpty   = []byte{}
	bEqual   = []byte{'='}
	bDQuote  = []byte{'"'}
)

type Config struct {
	filename string
	comment  map[int][]string
	data     map[string]string
	offset   map[string]int64

	sync.RWMutex
}

func LoadConfig(name string) (*Config, error) {
	// Open 和 OpenFile 的区别 --> OpenFile(name, O_RDONLY, 0)
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		file.Name(),
		make(map[int][]string),
		make(map[string]string),
		make(map[string]int64),
		sync.RWMutex{},
	}
	cfg.Lock()
	defer cfg.Unlock()
	defer file.Close()

	var comment bytes.Buffer
	buf := bufio.NewReader(file)
	for nComment, off := 0, int64(1); ; {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if bytes.Equal(line, bEmpty) {
			continue
		}

		off += int64(len(line))

		if bytes.HasPrefix(line, bComment) {
			line = bytes.TrimLeft(line, "#")
			line = bytes.TrimLeftFunc(line, unicode.IsSpace)
			comment.Write(line)
			comment.WriteByte('\n')
			continue
		}
		if comment.Len() != 0 {
			cfg.comment[nComment] = []string{strings.TrimSpace(comment.String())}
			comment.Reset()
			nComment++
		}

		vals := bytes.SplitN(line, bEqual, 2)
		if bytes.HasPrefix(vals[1], bDQuote) {
			vals[1] = bytes.Trim(vals[1], `"`)
		}

		key := strings.TrimSpace(string(vals[0]))
		cfg.comment[nComment-1] = append(cfg.comment[nComment-1], key)
		cfg.data[key] = strings.TrimSpace(string(vals[1]))
		cfg.offset[key] = off
	}
	return cfg, nil
}
