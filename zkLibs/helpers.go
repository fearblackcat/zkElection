package zkLibs

import (
	"strconv"
	"strings"
)

type SequenceStrings []string

func (s SequenceStrings) Len() int {
	return len(s)
}
func (s SequenceStrings) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SequenceStrings) Less(i, j int) bool {
	iParts := strings.Split(s[i], "-")
	iNumParts := strings.Split(iParts[1], "_")
	iNum, _ := strconv.ParseInt(iNumParts[1], 10, 64)

	jParts := strings.Split(s[j], "-")
	jNumParts := strings.Split(jParts[1], "_")
	jNum, _ := strconv.ParseInt(jNumParts[1], 10, 64)

	return iNum < jNum
}

// CreateEndpoints creates a list of endpoints given the right scheme
func CreateEndpoints(addrs []string, scheme string) (entries []string) {
	for _, addr := range addrs {
		entries = append(entries, scheme+"://"+addr)
	}
	return entries
}

// Normalize the key for each store to the form:
//
//     /path/to/key
//
func Normalize(key string) string {
	return "/" + join(SplitKey(key))
}

// GetDirectory gets the full directory part of
// the key to the form:
//
//     /path/to/
//
func GetDirectory(key string) string {
	parts := SplitKey(key)
	parts = parts[:len(parts)-1]
	return "/" + join(parts)
}

// SplitKey splits the key to extract path informations
func SplitKey(key string) (path []string) {
	if strings.Contains(key, "/") {
		path = strings.Split(key, "/")
	} else {
		path = []string{key}
	}
	return path
}

// join the path parts with '/'
func join(parts []string) string {
	return strings.Join(parts, "/")
}
