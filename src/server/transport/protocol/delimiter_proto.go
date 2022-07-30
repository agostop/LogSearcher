package protocol

import (
	"strings"
)

type DelimiterProto string

func (o DelimiterProto) Input(data []byte) int {
    if len(data) == 0 {
        return 0
    }
    stringData := string(data)
    i := strings.Index(stringData, string(o))
    if stringData[:1] == string(o) {
        i++
    }

    return i

}

func (o DelimiterProto) Package(data []byte) ([][]byte, []byte, error) {
    var result [][]byte
    stringData := string(data)
    i := o.Input(data)
    if i <= 0 {
        return nil, data, nil
    }

    for ; i > 0; i = strings.Index(stringData, string(o)) {
        if i+1 > len(stringData) {
            break
        }
        result = append(result, []byte(stringData[:i]))
        stringData = stringData[i+1:]
    }

    return result, []byte(stringData), nil

}
