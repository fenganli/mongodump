// Main package for the mongodump tool.
package main

import (
	"github.com/mongodb/mongo-tools/common/log"
	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/common/util"
	"github.com/mongodb/mongo-tools/mongodump"
	"github.com/mongodb/mongo-tools/mongorestore"
	"os"
	"fmt"
)

func main() {
	// initialize command-line opts
	// --ts
	// ts_string := flag.String("ts", "1:1", "time stamp string")
	// flag.Parse()
	// flag.StringVar(

	opts := options.New("mongodump", mongodump.Usage, options.EnabledOptions{true, true, true})

	inputOpts := &mongodump.InputOptions{}
	opts.AddOptions(inputOpts)
	outputOpts := &mongodump.OutputOptions{}
	opts.AddOptions(outputOpts)

	args, err := opts.Parse()
	if err != nil {
		log.Logf(log.Always, "error parsing command line options: %v", err)
		log.Logf(log.Always, "try 'mongodump --help' for more information")
		os.Exit(util.ExitBadOptions)
	}

	if len(args) > 0 {
		log.Logf(log.Always, "positional arguments not allowed: %v", args)
		log.Logf(log.Always, "try 'mongodump --help' for more information")
		os.Exit(util.ExitBadOptions)
	}

	// print help, if specified
	if opts.PrintHelp(false) {
		return
	}

	// print version, if specified
	if opts.PrintVersion() {
		return
	}
	
	// ts_string := flag.String("ts", "1:1", "time stamp string")
	// flag.Parse()
	// ts_string := "1450325712:1"

	// init logger
	log.SetVerbosity(opts.Verbosity)

	// connect directly, unless a replica set name is explicitly specified
	_, setName := util.ParseConnectionString(opts.Host)
	opts.Direct = (setName == "")
	opts.ReplicaSetName = setName

	dump := mongodump.MongoDump{
		ToolOptions:   opts,
		OutputOptions: outputOpts,
		InputOptions:  inputOpts,
	}
	// fmt.Printf(opts.TS)
	time_stamp, err := mongorestore.ParseTimestampFlag(opts.TS)
	if err != nil {
		fmt.Errorf("couldn't get the timestamp from command arguments")
		os.Exit(util.ExitError)
	}
	// dump.oplogStart = time_stamp 

	if err = dump.Init(time_stamp); err != nil {
		log.Logf(log.Always, "Failed: %v", err)
		os.Exit(util.ExitError)
	}
	//dump.OplogStart = time_stamp

	if err = dump.Dump(); err != nil {
		log.Logf(log.Always, "Failed: %v", err)
		if err == util.ErrTerminated {
			os.Exit(util.ExitKill)
		}
		os.Exit(util.ExitError)
	}
}
