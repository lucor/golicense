package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
)

// Version can be used to set the version at link time
var Version string

type options struct {
	include          bool
	listLicenses     bool
	listLicenseNames bool
	output           string
	verbose          bool
	version          bool
}

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: golicense [OPTIONS] [PATH]
List information about the licenses of a Go module or binary and its dependencies.
Default is to look for a go.mod file into the current directory.

Options:
  -h, --help             show this help message
  -i, --include          include the licenses in the output
      --list-names       list the names of the license file can be detected and exit
      --list-licenses    list the licenses can be detected and exit
  -o, --output <file>    write to file instead of stdout
  -v, --verbose          make the tool verbose
  -V, --version          show the version number
`)
	}

	var opts options

	flag.BoolVar(&opts.include, "i", false, "")
	flag.BoolVar(&opts.include, "include", false, "")
	flag.StringVar(&opts.output, "o", "", "")
	flag.StringVar(&opts.output, "output", "", "")
	flag.BoolVar(&opts.verbose, "v", false, "")
	flag.BoolVar(&opts.verbose, "verbose", false, "")
	flag.BoolVar(&opts.version, "V", false, "")
	flag.BoolVar(&opts.version, "version", false, "")
	flag.BoolVar(&opts.listLicenseNames, "list-names", false, "")
	flag.BoolVar(&opts.listLicenses, "list-licenses", false, "")
	flag.Parse()

	if opts.version {
		fmt.Println("golicense", version())
		os.Exit(0)
	}

	if opts.listLicenses {
		listLicenses()
		os.Exit(0)
	}

	if opts.listLicenseNames {
		listLicenseNames()
		os.Exit(0)
	}

	path := "go.mod"
	if len(flag.Args()) > 0 {
		path = flag.Arg(0)
	}
	mi, err := getModuleInfo(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	licenses, err := getLicenses(getGoModCache(), mi, licenseNames)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	w := os.Stdout
	if opts.output != "" {
		f, err := os.Create(opts.output)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
		w = f
	}

	for _, l := range licenses {
		fmt.Fprintf(w, "## %s (https://pkg.go.dev/%s?tab=licenses)\n", l.ModuleInfo, l.ModuleInfo)
		if opts.verbose {
			fmt.Fprintf(w, "Path: %s\n", l.Path)
			fmt.Fprintf(w, "License: %s\n", l.Type)
		}
		fmt.Fprintln(w)
		if opts.include {
			fmt.Fprintf(w, "%s", l.Content)
			fmt.Fprintln(w)
		}
	}
}

func version() string {
	if Version != "" {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}
	return "(unknown)"
}
