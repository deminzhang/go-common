package utils

import (
	"fmt"
	"regexp"
)

func ConstConflictCheck(src, pre string) {
	var constConflictSlice []string
	fmt.Printf("ConstConflictCheck %s*...\n", pre)
	r := fmt.Sprintf(`\s*(//)?\s*(%s\w+)\s*=\s*(\d+|"[^"]*").*$?`, pre)
	re := regexp.MustCompile(r)
	matches := re.FindAllStringSubmatch(src, -1)
	reg := make(map[string]string, len(matches))
	for _, kv := range matches {
		if kv[1] != "" {
			continue
		}
		k, v := kv[2], kv[3]
		if ko, b := reg[v]; b {
			constConflictSlice = append(constConflictSlice, fmt.Sprintf("Conflict %s[%s=%s=%s]", pre, k, ko, v))
		}
		reg[v] = k
	}
	if len(constConflictSlice) > 0 {
		fmt.Printf("ConstConflict:\n%s\n", constConflictSlice)
		panic(fmt.Sprintf("ConstConflictCheck %s* FAIL\n", pre))
	}
	fmt.Printf("ConstConflictCheck %s* OK\n", pre)
}
