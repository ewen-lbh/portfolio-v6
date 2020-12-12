package main

import (
	"bufio"
	"os"
	"strings"
)

// Remove duplicated PO headers and weird lines

func main() {
	seenHeadersNames := make([]string, 0)
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
    var newContents string
    var sawMsgidDeclaration bool
	for scanner.Scan() {
        line := scanner.Text()
        // Remove ugly "#-#-#-#-#" lines
        if isUglyMarker(line) {
            continue
        }
        if strings.HasPrefix(line, "#: ") {
            sawMsgidDeclaration = true
        }
        if isHeader(line) && !sawMsgidDeclaration {
            headerName := getHeaderName(line)
            alredySeen := false
            for _, seenHeaderName := range seenHeadersNames {
                if headerName == seenHeaderName {
                    alredySeen = true
                }
            }
            if alredySeen {
                continue
            }
            seenHeadersNames = append(seenHeadersNames, headerName)
        }
        newContents += line + "\n"
    }
    file.Close()
    file, err = os.OpenFile(os.Args[2], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
    if err != nil {
        panic(err)
    }
    _, err = file.WriteString(newContents)
    if err != nil {
        panic(err)
    }
}

// isUglyMarker determines if the line is an ugly "#-#-#-#-#" line
func isUglyMarker(line string) bool {
    return strings.HasPrefix(line, "\"#-#-#-#-#")
}

func isHeader(line string) bool {
	return !HasPrefixes(line, "#", "msgid", "msgstr") && len(getHeaderName(line)) != len(line)
}

func getHeaderName(line string) string {
	// header has the shape [name, value]
    header := strings.SplitN(strings.Trim(line, "\""), ":", 2)
	return header[0]
}

// HasPrefixes returns true is any of the prefixes are prefixes of s
func HasPrefixes(s string, prefixes ...string) bool {
    for _, prefix := range prefixes {
        if strings.HasPrefix(s, prefix) {
            return true
        }
    }
    return false
}
