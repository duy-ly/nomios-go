package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-mysql-org/go-mysql/mysql"
)

func BuildEventPos(pos *mysql.Position) string {
	return fmt.Sprintf("%s:%d", pos.Name, pos.Pos)
}

func ParseEventPos(s string) mysql.Position {
	var (
		binlogFile string
		binlogPos  uint64
	)
	pp := strings.SplitN(s, ":", 2)
	if len(pp) == 2 {
		binlogFile = pp[0]
		binlogPos, _ = strconv.ParseUint(pp[1], 10, 32)
	}

	return mysql.Position{
		Name: binlogFile,
		Pos:  uint32(binlogPos),
	}
}
