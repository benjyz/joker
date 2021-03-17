package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	. "github.com/candid82/joker/core"
)

var client *http.Client

func extractMethod(request Map) string {
	if ok, m := request.Get(MakeKeyword("method")); ok {
		switch m := m.(type) {
		case String:
			return m.S
		case Keyword:
			return m.ToString(false)[1:]
		case Symbol:
			return m.ToString(false)
		default:
			panic(RT.NewError(fmt.Sprintf("method must be a string, keyword or symbol, got %s", m.GetType().ToString(false))))
		}
	}
	return "get"
}

func getOrPanic(m Map, k Object, errMsg string) Object {
	if ok, v := m.Get(k); ok {
		return v
	}
	panic(RT.NewError(errMsg))
}

func mapToReq(request Map) *http.Request {
	method := strings.ToUpper(extractMethod(request))
	url := EnsureObjectIsString(getOrPanic(request, MakeKeyword("url"), ":url key must be present in request map"), "url: %s").S
	var reqBody io.Reader
	if ok, b := request.Get(MakeKeyword("body")); ok {
		reqBody = strings.NewReader(EnsureObjectIsString(b, "body: %s").S)
	}
	req, err := http.NewRequest(method, url, reqBody)
	PanicOnErr(err)
	if ok, headers := request.Get(MakeKeyword("headers")); ok {
		h := EnsureObjectIsMap(headers, "headers: %s")
		for iter := h.Iter(); iter.HasNext(); {
			p := iter.Next()
			req.Header.Add(EnsureObjectIsString(p.Key, "header name: %s").S, EnsureObjectIsString(p.Value, "header value: %s").S)
		}
	}
	if ok, host := request.Get(MakeKeyword("host")); ok {
		req.Host = EnsureObjectIsString(host, "host: %s").S
	}
	return req
}

func reqToMap(host String, port String, req *http.Request) Map {
	defer req.Body.Close()
	res := EmptyArrayMap()
	body, err := ioutil.ReadAll(req.Body)
	PanicOnErr(err)
	res.Add(MakeKeyword("request-method"), MakeKeyword(strings.ToLower(req.Method)))
	res.Add(MakeKeyword("body"), MakeString(string(body)))
	res.Add(MakeKeyword("uri"), MakeString(req.URL.Path))
	res.Add(MakeKeyword("query-string"), MakeString(req.URL.RawQuery))
	res.Add(MakeKeyword("server-name"), host)
	res.Add(MakeKeyword("server-port"), port)
	res.Add(MakeKeyword("remote-addr"), MakeString(req.RemoteAddr[:strings.LastIndexByte(req.RemoteAddr, byte(':'))]))
	res.Add(MakeKeyword("protocol"), MakeString(req.Proto))
	res.Add(MakeKeyword("scheme"), MakeKeyword("http"))
	res.Add(MakeKeyword("host"), MakeString(req.Host))
	headers := EmptyArrayMap()
	for k, v := range req.Header {
		headers.Add(MakeString(strings.ToLower(k)), MakeString(strings.Join(v, ",")))
	}
	res.Add(MakeKeyword("headers"), headers)
	return res
}

func respToMap(resp *http.Response) Map {
	defer resp.Body.Close()
	res := EmptyArrayMap()
	body, err := ioutil.ReadAll(resp.Body)
	PanicOnErr(err)
	res.Add(MakeKeyword("body"), MakeString(string(body)))
	res.Add(MakeKeyword("status"), MakeInt(resp.StatusCode))
	respHeaders := EmptyArrayMap()
	for k, v := range resp.Header {
		respHeaders.Add(MakeString(k), MakeString(strings.Join(v, ",")))
	}
	res.Add(MakeKeyword("headers"), respHeaders)
	// TODO: 32-bit issue
	res.Add(MakeKeyword("content-length"), MakeInt(int(resp.ContentLength)))
	return res
}

func mapToResp(response Map, w http.ResponseWriter) {
	status := 0
	if ok, s := response.Get(MakeKeyword("status")); ok {
		status = EnsureObjectIsInt(s, "HTTP response status: %s").I
	}
	body := ""
	if ok, b := response.Get(MakeKeyword("body")); ok {
		body = EnsureObjectIsString(b, "HTTP response body: %s").S
	}
	if ok, headers := response.Get(MakeKeyword("headers")); ok {
		header := w.Header()
		h := EnsureObjectIsMap(headers, "HTTP response headers: %s")
		for iter := h.Iter(); iter.HasNext(); {
			p := iter.Next()
			hname := EnsureObjectIsString(p.Key, "HTTP response header name %s").S
			switch pvalue := p.Value.(type) {
			case String:
				header.Add(hname, pvalue.S)
			case Seqable:
				s := pvalue.Seq()
				for !s.IsEmpty() {
					header.Add(hname, EnsureObjectIsString(s.First(), "HTTP response header value: %s").S)
					s = s.Rest()
				}
			default:
				panic(RT.NewError("HTTP response header value must be a string or a seq of strings"))
			}
		}
	}
	if status != 0 {
		w.WriteHeader(status)
	}
	io.WriteString(w, body)
}

func sendRequest(request Map) Map {
	req := mapToReq(request)
	RT.GIL.Unlock()
	resp, err := client.Do(req)
	RT.GIL.Lock()
	PanicOnErr(err)
	return respToMap(resp)
}

func startServer(addr string, handler Callable) Object {
	i := strings.LastIndexByte(addr, byte(':'))
	host, port := MakeString(addr), MakeString("")
	if i != -1 {
		host = MakeString(addr[:i])
		port = MakeString(addr[i+1:])
	}
	RT.GIL.Unlock()
	defer RT.GIL.Lock()
	err := http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		RT.GIL.Lock()
		defer func() {
			RT.GIL.Unlock()
			if r := recover(); r != nil {
				w.WriteHeader(500)
				io.WriteString(w, "Internal server error")
				fmt.Fprintln(os.Stderr, r)
			}
		}()
		response := handler.Call([]Object{reqToMap(host, port, req)})
		mapToResp(EnsureObjectIsMap(response, "HTTP response: %s"), w)
	}))
	PanicOnErr(err)
	return NIL
}

func startFileServer(addr string, root string) Object {
	err := http.ListenAndServe(addr, http.FileServer(http.Dir(root)))
	PanicOnErr(err)
	return NIL
}

func initNative() {
	client = &http.Client{}
}
