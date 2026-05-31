package main

import (
	"strings"
)

type HandlerFunc func(req *HttpRequest, res *HttpResponse)
type Segment struct {
	Value   string
	isParam bool
}
type Route struct {
	Method  string
	Pattern []Segment
	Handler HandlerFunc
}

type HttpRouter struct {
	Routes []Route
}

func (h *HttpRouter) Handle(method string, path string, handler HandlerFunc) {
	h.Routes = append(
		h.Routes,
		Route{
			Method:  method,
			Pattern: ParsePattern(path),
			Handler: handler,
		},
	)
}

func (h *HttpRouter) Dispatch(req *HttpRequest, res *HttpResponse) {
	for _, route := range h.Routes {
		if route.Method != req.StartLine.method {
			continue
		}
		params, exists := matchPattern(route.Pattern, req.StartLine.path)
		if !exists {
			continue
		}
		req.Params = params
		route.Handler(req, res)
		return
	}
	res.SetStatus(HttpNotFound)
	res.Write([]byte(`{"error": "Invalid Path"}`))
}

func matchPattern(segments []Segment, path string) (map[string]string, bool) {
	reqSegment := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) != len(reqSegment) {
		return nil, false
	}
	params := make(map[string]string)
	for i, s := range segments {
		// check if parts & params match
		if s.isParam {
			params[s.Value] = reqSegment[i]
		} else if s.Value != reqSegment[i] {
			return nil, false
		}
	}

	return params, true
}

func ParsePattern(path string) []Segment {
	// path1/path2/{str}/path3/{str}
	seg := strings.Split(strings.Trim(path, "/"), "/")
	// path is equal to number of segments; like, path1 + path2 + ...
	// we will build the path with segment with param + no param
	segments := make([]Segment, 0, len(seg))
	for _, s := range seg {
		if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
			segments = append(segments, Segment{
				Value:   s[1 : len(s)-1],
				isParam: true,
			})
		} else {
			segments = append(segments, Segment{
				Value:   s,
				isParam: false,
			})
		}
	}
	return segments
}
