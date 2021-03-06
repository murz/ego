package tmpl

import (
	// "html/template"
	"io/ioutil"
	// "io"
	"fmt"
	"log"
	"path"
	"net/http"
)

type TemplateEngine interface {
	Compile(string) CompiledTemplate
}

type CompiledTemplate interface {
	Execute(map[string]interface{}) (*string, error)
	ExecuteRW(http.ResponseWriter, map[string]interface{}) error
}

var templateEngine TemplateEngine = &HandlebarsTemplateEngine{}
var templates = make(map[string]CompiledTemplate)

func SetTemplateEngine(e TemplateEngine) {
	templateEngine = e
}

func Render(w http.ResponseWriter, name string, ctx map[string]interface{}) {
	if (templates[name] == nil) {
		panic(fmt.Sprintf("Unknown template: %s", name))
	}
log.Printf("execute rw..")
	err := templates[name].ExecuteRW(w, ctx)
	if (err != nil) {
		panic(err)
	}
}

func Parse(filename string) {
	// basePath := "app/views/"
	// dirlist, err := ioutil.ReadDir(basePath + "layouts")
	// if err != nil {
	// 	log.Fatalf("Error reading /app/views/layouts")
	// }
	// for _, f := range dirlist {
	// 	file, _ := ioutil.ReadFile(basePath + "layouts/" + f.Name())
	// 	templateEngine.Compile(string(file))
	// }
	// file, _ = iotuil.ReadFile(filename)
	log.Printf("YARRR: %v", filename)
	templates[filename] = templateEngine.Compile(filename)
}

func ParseDir(dir string) {
	// ParsePartials(dirname)
	dirlist, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error reading %s: %s\n", dir, err)
	}
	for _, f := range dirlist {
		filename := path.Join(dir, f.Name())
		if f.Name() == "layouts" || f.Name()[0:1] == "_" {
			continue
		}
		if f.IsDir() {
			ParseDir(filename)
		} else {
			Parse(filename)
		}
	}
}

///////////////////////////////////////////////////////

// var pkgName string
// var templates = make(map[string]*template.Template)
// var helpers = make([]*Helper, 0)
// var partials = make([]*Partial, 0)

// func ParseDir(filename string) {
// 	basePath := pkgName + "/app/views/"
// 	dirlist, err := ioutil.ReadDir(basePath + "layouts")
// 	if err != nil {
// 		log.Fatalf("Error reading /app/views/layouts")
// 	}
// 	for _, f := range dirlist {
// 		filenames := []string {
// 			// basePath + "layouts/" + f.Name(),
// 			filename,
// 		}
// 		key := f.Name() + "||" + filename
// 		keynl := filename
// 		if pkgName != "" {
// 			key = strings.Replace(key, basePath, "", -1)
// 			keynl = strings.Replace(keynl, basePath, "", -1)
// 		}
// 		t := template.New(key)
// 		tnl := template.New(filename)
// 		for _, helper := range helpers {
// 			t.Funcs(template.FuncMap{
// 				helper.Name: helper.Execute,
// 	        })
// 	        tnl.Funcs(template.FuncMap{
// 				helper.Name: helper.Execute,
// 	        })
// 		}
// 		for _, partial := range partials {
// 			partialTemplate := t.New(partial.Name)
// 			partialTemplate.Parse(partial.TemplateString)
// 			partialTemplateNoLayout := tnl.New(partial.Name)
// 			partialTemplateNoLayout.Parse(partial.TemplateString)
// 		}
// 		file, _ := ioutil.ReadFile(basePath + "layouts/" + f.Name())
// 		filenl, _ := ioutil.ReadFile(basePath + keynl)
// 		t.Parse(string(file))
// 		t.ParseFiles(filenames...)
// 		tnl.Parse(string(filenl))
// 		if (err != nil) {
// 			log.Panic(err)
// 		}
// 		templates[key] = t
// 		log.Printf("writing %v to %v", filename, tnl)
// 		templates[keynl] = tnl
// 	}
// }

// func ParseDir(dirname string) {
// 	ParsePartials(dirname)
// 	dirlist, err := ioutil.ReadDir(pkgName + dirname)
// 	if err != nil {
// 		log.Fatalf("Error reading %s: %s\n", dirname, err)
// 	}
// 	for _, f := range dirlist {
// 		filename := path.Join(pkgName, dirname, f.Name())
// 		if f.Name() == "layouts" || f.Name()[0:1] == "_" {
// 			continue
// 		}
// 		if f.IsDir() {
// 			ParseDir(path.Join(dirname, f.Name()))
// 		} else {
// 			Parse(filename)
// 		}
// 	}
// }

// func ParsePartials(dirname string) {
// 	basePath := pkgName + dirname
// 	dirlist, err := ioutil.ReadDir(basePath)
// 	if err != nil {
// 		log.Fatalf("Error reading %s: %s\n", dirname, err)
// 	}
// 	for _, f := range dirlist {
// 		filename := path.Join(pkgName, dirname, f.Name())
// 		if f.Name()[0:1] == "_" {
// 			key := f.Name()[1:len(f.Name())]
// 			key = strings.Replace(key, ".html", "", -1)
// 			// t := template.New(key)
// 			file, _ := ioutil.ReadFile(filename)
// 			// t.Parse(string(file))
// 			partial := &Partial{
// 				Name: key,
// 				TemplateString: string(file),
// 			}
// 			log.Printf("adding partial: %v", key)
// 			partials = append(partials, partial)
// 		}
// 	}
// }

// func Render(wr io.Writer, key string, layout string, context interface{}) {
// 	if (layout == "none") {
// 		layout = ""
// 	}
// 	if (layout != "") {
// 		layout = layout+"||"
// 	}
// 	log.Printf("looking up: %v", layout+key)
// 	t := templates[layout + key]
// 	log.Printf("found: %v", t)

// 	log.Printf("context %v", context)
// 	err := t.Execute(wr, context)
// 	if (err != nil) {
// 		panic(err)
// 	}
// 	// TODO: do something with err instead of _
// }

// func SetPackageName(name string) {
// 	pkgName = name
// }

// func RegisterHelper(helper *Helper) *Helper {
// 	log.Printf("## REGISTERED helper: %v", helper.Name)
// 	helpers = append(helpers, helper)
// 	return helper
// }