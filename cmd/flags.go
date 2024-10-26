package cmd

import "flag"

func Flags() string {
	path := flag.String("path", "config.json", "use for load configuration file")

	flag.Parse()

	return *path
}
