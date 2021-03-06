package http

import (
	"reflect"
	"fmt"
	"log"
	"time"
	"strconv"
	"strings"
        "os"
	// "errors"
	nhttp "net/http"
)

// An http result can be anything, the dispatcher will try to figure out
// what to do with the result based on it's type.
type Result interface{}

type Context []interface{}

type ContextMap map[string]interface{}

// Metadata stored about the controllers at startup time.
type ControllerMetadata struct {
	Type reflect.Type
	ContextKeys []string
	Fields map[string]string
}

var controllerRegistry = make(map[string]*ControllerMetadata)

// Bind action key in the form of "Controller.Action" to a controller type.
func RegisterAction(key string, t reflect.Type, keys []string, fields map[string]string) {
	// TODO: handle error case
	log.Printf("REGISTER %v: %v", key, keys)
	controllerRegistry[key] = &ControllerMetadata{
		t,
		keys,
		fields,
	}
}

// Returns a net/http handler function for dispatching ego actions using the given router.
func ActionDispatchHandler(r *Router) nhttp.HandlerFunc {
	// setup static file server   
	var publicMux = nhttp.NewServeMux()
        var publicFs = justFilesFilesystem{nhttp.Dir("./public")}
	publicMux.Handle("/", nhttp.FileServer(publicFs))

	var assetsMux = nhttp.NewServeMux()
	assetsMux.Handle("/", nhttp.StripPrefix("/assets", nhttp.FileServer(justFilesFilesystem{nhttp.Dir("./.ego-genfiles/assets")})))

	return func(w nhttp.ResponseWriter, httpReq *nhttp.Request) {
		startTime := time.Now().UnixNano() / 1000000
		// reqType := ""

		route, pathParams, found := r.Lookup(httpReq.URL.Path, httpReq.Method)
		if !found {
			// try the wildcard tree
			route, _, found = r.Lookup(httpReq.URL.Path, "*")
			if !found {
                                wrapper := NewStaticFileResponseWriter(w)
                                if (strings.HasPrefix(httpReq.URL.Path, "/assets")) {
                                  if strings.HasSuffix(httpReq.URL.Path, "/assets") || strings.HasSuffix(httpReq.URL.Path, "/assets/") {
                                    wrapper.WriteHeader(nhttp.StatusNotFound);
                                    return;
                                  }
                                  assetsMux.ServeHTTP(wrapper, httpReq)
                                } else {
                                  if (strings.HasSuffix(httpReq.URL.Path, "/")) {
                                    if _, err := publicFs.Open(httpReq.URL.Path + "/index.html"); err != nil {
                                      wrapper.WriteHeader(nhttp.StatusNotFound);
                                      return;
                                    }
                                  }
                                  publicMux.ServeHTTP(wrapper, httpReq)
                                }
				// NotFoundAction.Dispatch(w, httpReq, nil, reqType)
				return;
			}
		}
		metadata := controllerRegistry[fmt.Sprintf("%s.%s", route.ControllerName, route.ActionName)]
		ctrlType := metadata.Type
		ctrlVal := reflect.New(ctrlType)
		cfgMethod := ctrlVal.MethodByName("Configure")
		if cfgMethod.IsValid() {
			cfgMethod.Call([]reflect.Value{})
		}
		req := NewRequest()
		req.Parse(httpReq)
		params := make([]reflect.Value, 0)
		log.Printf("mfields: %v", metadata.Fields);
		for key, t := range metadata.Fields {
			log.Printf("FIELD %v:%v", key, t)
			switch(t) {
			case "int":
				p := httpReq.FormValue(key)
				if pathParam, ok := pathParams[key]; ok {
					p = pathParam
				}
				if i, err := strconv.ParseInt(p, 10, 0); err == nil {
					params = append(params, reflect.ValueOf(int(i)))
				}
			case "string":
				p := httpReq.FormValue(key)
				params = append(params, reflect.ValueOf(p))
			case "&amp;{http Request}":
				params = append(params, reflect.ValueOf(*req))
			}
		}
		method := ctrlVal.MethodByName(route.ActionName)
		resultVal := method.Call(params)[0]
		result := resultVal.Interface()
		var resp *Response
		if ctx, ok := result.(Context); ok {
			ctxMap := make(map[string]interface{})
			log.Printf("******")
			log.Printf("ctx: %v", ctx)
			for i, obj := range ctx {
				ctxMap[metadata.ContextKeys[i]] = obj
			}
			resp = &Response{
				Context: ctxMap,
			}
			log.Printf("resp: %v", resp)
		}
		if resp == nil {
			if r, ok := result.(*Response); ok {
				resp = r
			} else {
				resp = &Response{
				}
			}
		}
		log.Printf("resp: %v", resp)
		ctrlName := strings.Replace(route.ControllerName, "Controller", "", -1)
		fmt.Sprintf("ctrl: %v", strings.ToLower(ctrlName))
		fmt.Sprintf("str: %v", strings.ToLower(route.ActionName))
		resp.View = fmt.Sprintf("app/views/%v/%v.html.hbs",
			strings.ToLower(ctrlName), strings.ToLower(route.ActionName))
		if (resp.StatusCode != 0) {
			w.WriteHeader(resp.StatusCode)
		}
		log.Printf("%s", resp)
		resp.WriteHTML(w)
		log.Printf("%s %s -> %s.%s (%dms)",
			route.Path.Method,
			route.Path.Value,
			route.ControllerName,
			route.ActionName,
			(time.Now().UnixNano() / 1000000) - startTime);
	}
}

type staticFileResponseWriter struct {
	writer nhttp.ResponseWriter
}

func NewStaticFileResponseWriter(rw nhttp.ResponseWriter) staticFileResponseWriter {
	return staticFileResponseWriter{
		writer: rw,
	}
}

func (rw staticFileResponseWriter) Header() nhttp.Header {
	return rw.writer.Header()
}

func (rw staticFileResponseWriter) Write(data []byte) (int, error) {
	if string(data) == "404 page not found\n" {
		return 0, nil
	}
	return rw.writer.Write(data)
}

func (rw staticFileResponseWriter) WriteHeader(code int) {
	rw.writer.WriteHeader(code)
	if code == nhttp.StatusNotFound {
		fmt.Fprintln(rw, "404 bitches!")
	}
}

type justFilesFilesystem struct {
  fs nhttp.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (nhttp.File, error) {
  f, err := fs.fs.Open(name)
  if err != nil {
    return nil, err
  }
  return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
  nhttp.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
fmt.Println("read dir??")
  return nil, nil
}
