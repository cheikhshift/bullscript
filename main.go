package main


import (
	"flag"
	"io/ioutil"
	"strings"
	"github.com/cheikhshift/gos/core"
	"fmt"
)

const launchername = "BSLauncher.go"
const handlerTemp = `http.HandleFunc("%s/%s", %s)`
const portConst = `go http.ListenAndServe(":%s", http.HandlerFunc(redirect(%s)))`
const formatterstr = `gofmt -w %s`
const listenStmt = `errgos := http.ListenAndServe(":%s", nil)
	if errgos != nil {
		log.Fatal(errgos)
	}`

const listenSecure = `errgos := http.ListenAndServeTLS(":%s", "%s", "%s", nil)
	if errgos != nil {
		log.Fatal("ListenAndServe: ", errgos)
	}`
const WebAppProg = `
package main


import (
	"net/http"
	"log"
	%s
)

const Hostsubstr = "%%s://%%s/%%s"



func redirect(port int) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request)  {
	// remove/add not default ports from req.Host
    schema := "http"
    if port == 443 {
    	schema = "https"
    }
	target := fmt.Sprintf(Hostsubstr, schema,req.URL.Host,req.URL.Path)
	if len(req.URL.RawQuery) > 0 {
		target += fmt.Sprintf("?%%s", req.URL.RawQuery)
	}

	http.Redirect(w, req, target, 301)
	}
}

func main(){

	%s

	%s

	%s

}

`


func main(){

	flagPath := flag.String("path", "", "path to bull script file")

	flag.Parse()

	fileData,err := ioutil.ReadFile(*flagPath)

	if err != nil {
		panic(err)
	}

	commands := strings.Split(string(fileData),"\n")
	lengthCommands := len(commands)
	var LaunchStmts []string
	var imports []string
	var prefix string
	var Handlers []string
	//var litpath = strings.Replace(*flagPath,"\\","/",-1)
	//var pathbits = strings.Split(litpath, "/")
	var Mainlistenner string
	
	for i := 0; i < lengthCommands ; i++ {
		line := commands[i]
		bits := strings.Split(line, ">")
		for o := 0; o < len(bits); o++ {
			bits[o] = strings.TrimSpace(bits[o])
		}
		if strings.Contains(line, "@onstart") {
			if len(bits) < 2 {
				panic(fmt.Sprintf("Error on line %v. @onstart requires 1 arguments:\n ie : @onstart> goPkg.SomeFunc(\"var_string\") ", i + 1))
			}
			LaunchStmts = append(LaunchStmts, bits[1])
		}

		if strings.Contains(line, "@i") {
			if len(bits) < 2 {
				panic(fmt.Sprintf("Error on line %v. @i requires 1 arguments:\n ie : @i>\"Package/path/import\" or @i> customname \"Package/path/import\" ", i + 1))
			}
			if bits[1] != "net/http" && bits[1] != "log" {
				imports = append(imports, bits[1])
			}
		}

		if strings.Contains(line, "@prefix") {
			if len(bits) < 2 {
				panic(fmt.Sprintf("Error on line %v. @prefix requires 1 arguments:\n ie : @prefix > /url/to/prefix/with/sub/handlers", i + 1))
			}
			prefix = bits[1]
		}

		if strings.Contains(line, "@path") {
			if len(bits) < 3 {
				panic(fmt.Sprintf("Error on line %v. @path requires 2 arguments:\n ie : @path /url/path package.Handler", i + 1))
			}
			hpath := fmt.Sprintf(handlerTemp, prefix, bits[1], bits[2]) 
			hpath = strings.Replace(hpath,"//","/", -1)

			Handlers = append(Handlers, hpath)
		}


		if strings.Contains(line, "@listen") {
			if len(bits) < 2 {
				panic(fmt.Sprintf("Error on line %v. @listen requires 1 arguments:\n ie : @listen > PORT", i + 1))
			}
			Mainlistenner = fmt.Sprintf(listenStmt , bits[1])
		}

		if strings.Contains(line, "@listensecure") {
			if len(bits) < 3 {
				panic(fmt.Sprintf("Error on line %v. @listensecure requires 3 arguments:\n ie : @listensecure > PORT >path_to_certificate_file > path_to_key_file", i + 1))
			}
			Mainlistenner = fmt.Sprintf(listenSecure, bits[1], bits[2], bits[3] )
		}

		if strings.Contains(line, "@redirect") {
			if len(bits) < 3 {
				panic(fmt.Sprintf("Error on line %v. @redirect requires 2 arguments:\n ie : @redirect > SOURCE_PORT > TARGET_PORT", i + 1))
			}
			bits[1] = strings.TrimSpace(bits[1])
			LaunchStmts = append(LaunchStmts, fmt.Sprintf(portConst, bits[1], bits[2]))
		}

		if strings.Contains(line, "@end") {
			prefix = ""
		}

		if strings.Contains(line, "@run") {
			if len(bits) < 2 {
				panic(fmt.Sprintf("Error on line %v. @run requires 1 arguments:\n ie : @run > echo \"Hello world\"", i + 1))
			}
			core.RunCmd(bits[1])
		}

	}


	finalStr := fmt.Sprintf(WebAppProg, strings.Join(imports,"\n"), strings.Join(LaunchStmts,"\n" ), strings.Join(Handlers,"\n"), Mainlistenner  )

	newfile := strings.Replace(*flagPath, ".bs", ".go" , -1)

	err = ioutil.WriteFile(newfile,[]byte(finalStr), 0700 )
	if err != nil {
		panic(err)
	}

	core.RunCmd(fmt.Sprintf(formatterstr, newfile))

}