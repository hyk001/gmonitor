package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/simplejia/utils"

	"github.com/zhaochuanyun/gmonitor/comm"
	"os/exec"
)

type Conf struct {
	Port     int
	RootPath string
	Environ  string
	Svrs     map[string]string
	Clog     *struct {
		Name  string
		Addr  string
		Mode  int
		Level int
	}
}

var (
	Envs                                       map[string]*Conf
	Env                                        string
	C                                          *Conf
	Start, Stop, Restart, GraceRestart, Status string
)

func init() {
	flag.StringVar(&Start, comm.START, "", "start a svr")
	flag.StringVar(&Stop, comm.STOP, "", "stop a svr")
	flag.StringVar(&Restart, comm.RESTART, "", "restart a svr")
	flag.StringVar(&GraceRestart, comm.GRESTART, "", "grace restart a svr")
	flag.StringVar(&Status, comm.STATUS, "", "status a svr")
	flag.StringVar(&Env, "env", "dev", "set env")

	var conf string
	flag.StringVar(&conf, "conf", "", "set custom conf")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Another process monitor\n")
		fmt.Fprintf(os.Stderr, "version: 1.7\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	dir := filepath.Dir(path)
	fcontent, err := ioutil.ReadFile(filepath.Join(dir, "conf", "conf.json"))
	//fcontent, err := ioutil.ReadFile(filepath.Join("/Users/mvpzhao/go/src/github.com/zhaochuanyun/gmonitor", "conf", "conf.json"))
	if err != nil {
		println("conf.json not found")
		os.Exit(-1)
	}

	fcontent = utils.RemoveAnnotation(fcontent)
	if err := json.Unmarshal(fcontent, &Envs); err != nil {
		fmt.Println("conf.json wrong format:", err)
		os.Exit(-1)
	}

	C = Envs[Env]
	if C == nil {
		fmt.Println("env not right:", Env)
		os.Exit(-1)
	}

	func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("conf not right:", err)
				os.Exit(-1)
			}
		}()

		matchs := regexp.MustCompile(`[\w|\.]+|".*?[^\\"]"`).FindAllString(conf, -1)
		for n, match := range matchs {
			matchs[n] = strings.Replace(strings.Trim(match, "\""), `\"`, `"`, -1)
		}
		for n := 0; n < len(matchs); n += 2 {
			name, value := matchs[n], matchs[n+1]

			rv := reflect.Indirect(reflect.ValueOf(C))
			for _, field := range strings.Split(name, ".") {
				rv = reflect.Indirect(rv.FieldByName(strings.Title(field)))
			}
			switch rv.Kind() {
			case reflect.String:
				rv.SetString(value)
			case reflect.Bool:
				b, err := strconv.ParseBool(value)
				if err != nil {
					panic(err)
				}
				rv.SetBool(b)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				i, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					panic(err)
				}
				rv.SetInt(i)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				u, err := strconv.ParseUint(value, 10, 64)
				if err != nil {
					panic(err)
				}
				rv.SetUint(u)
			case reflect.Float32, reflect.Float64:
				f, err := strconv.ParseFloat(value, 64)
				if err != nil {
					panic(err)
				}
				rv.SetFloat(f)
			}
		}
	}()

	fmt.Printf("Env: %s\nC: %s\n", Env, utils.Iprint(C))

	return
}
