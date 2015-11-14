package main

import (
	"flag"
	"os"
	"raft/coordinate"
	"raft/env"
)

var (
	path = flag.String("p", "conf/", "config path")
	file = flag.String("f", "raft.conf", "config file")
)

func main() {
	flag.Parse()

	if err := env.Init(*path, *file); err != nil {
		os.Exit(-1)
	}

	if cord, err := coordinate.NewCoordinator(env.Conf.Cord.Id, env.Conf.Cord.Addr); err != nil {
		env.Log.Error("create coordinator error : %s", err)
		os.Exit(-1)
	} else {
		cord.Work()
	}
}
