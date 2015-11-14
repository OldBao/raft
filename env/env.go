package env

import (
	"code.google.com/p/gcfg"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"path"
)

type Config struct {
	Peer struct {
		Addr []string
	}
	Cord struct {
		Addr string
		Id   int
	}
}

var (
	Log  l4g.Logger
	Conf Config
)

func Init(cpath, file string) error {
	Log = make(l4g.Logger)

	Log.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())
	//Log.AddFilter("log", log4go.FINE, log4go.NewFileLogWriter("example.log", true))

	if err := gcfg.ReadFileInto(&Conf, path.Join(cpath, file)); err != nil {
		fmt.Println("Load Conf error : ", err)
		return err
	}

	return nil
}
