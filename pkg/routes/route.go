package routes

import (
	"strconv"
	"strings"
)

type Route interface {
	Route() string
	Path(...interface{}) string
}

type PlainRoute struct {
	route string
}

func (r PlainRoute) Route() string              { return r.route }
func (r PlainRoute) Path(...interface{}) string { return r.route }

type StringRoute struct {
	route string
	keys  []string
}

func (r StringRoute) Route() string {
	return r.route
}

func (r StringRoute) Path(ks ...interface{}) string {
	if len(ks) != 1 {
		return ""
	}

	path := r.route
	for i, key := range r.keys {
		idstr := "-"
		id, ok := ks[i].(uint)
		if ok {
			idstr = strconv.FormatUint(uint64(id), 10)
		}

		path = strings.Replace(path, key, idstr, 1)
	}
	return path
}
