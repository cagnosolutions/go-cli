package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// The comment section below is an example of what we would like to do

/*
	$ my-program -h
	Usage of my-program:
  		-status string
    		see the status of the service
  		-service string
    		specify the service you wish to
  		-name string
    		specify your name (default "Guest")
*/

func main() {
	// prog0()
	// prog1()
	prog2()
}

func ClassicUsage(w io.Writer, prog string, varOpts bool) {
	if !varOpts {
		fmt.Fprintf(w, "Usage:\n  %s [options]\n\nOptions:\n", prog)
	} else {
		fmt.Fprintf(w, "Usage:\n  %s [options] [variable[=value] ...]\n\nOptions:\n", prog)
	}
}

type CLITool struct {
	Name  string
	Usage func()
	*flag.FlagSet
}

func NewCLIProg(errorHandling flag.ErrorHandling, varOpts bool) *CLITool {
	name := filepath.Base(os.Args[0])
	c := &CLITool{
		Name:    name,
		FlagSet: flag.NewFlagSet(name, errorHandling),
	}
	c.Usage = func() {
		ClassicUsage(c.Output(), c.Name, varOpts)
		c.PrintDefaults()
	}
	var helpFlag bool
	c.BoolVar(&helpFlag, "help", false, "print usage and help")
	return c
}

func (c *CLITool) NoArgFlag(name string, usage string) *bool {
	p := new(bool)
	c.BoolVar(p, name, false, usage)
	return p
}

func (c *CLITool) ChoiceFlag() {
	// TODO: implement
}

func (c *CLITool) ParseFlags() {
	err := c.Parse(os.Args[1:])
	if c.NArg() == 0 && err == nil {
		c.Usage()
	}
}

func (c *CLITool) String() string {
	ss := fmt.Sprintf("Flags=[")
	var sss []string
	c.Visit(
		func(f *flag.Flag) {
			sss = append(sss, fmt.Sprintf("%q", f.Name))
		},
	)
	ss += strings.Join(sss, ",")
	ss += "]"
	ss += fmt.Sprintf("\nName=%q\nNArgs=%d\nArgs=%+v\n", c.Name, c.NArg(), c.Args())
	return ss
}

func prog0() {
	// Initialize a new CLI tool
	cli := NewCLIProg(flag.ContinueOnError, true)

	// Add any flags you wish
	cli.NoArgFlag("start", "start a service")
	cli.NoArgFlag("stop", "stop a service")
	cli.NoArgFlag("restart", "restart a service")
	cli.NoArgFlag("list", "list the available services")

	// Call parse flags
	cli.ParseFlags()

	fmt.Println(cli)
}

func prog1() {

	// add help flag to override the defaults
	var helpShort, helpLong bool
	flag.BoolVar(&helpShort, "h", false, "print help or usage details")
	flag.BoolVar(&helpLong, "help", false, "print help or usage details")

	// set up the usage function
	usage := func() {
		ss := fmt.Sprintf("Usage:\n  %s [options] [variable[=value] ...]\n\n", filepath.Base(os.Args[0]))
		ss += fmt.Sprintf("Options:\n")
		fmt.Print(ss)
		flag.PrintDefaults()
	}

	// declare flags
	var startFlag bool
	var stopFlag bool
	var restartFlag bool
	var statusFlag bool
	var listFlag bool
	var serviceFlag string

	// set up the flags
	flag.BoolVar(&startFlag, "start", false, "start a service")
	flag.BoolVar(&stopFlag, "stop", false, "stop a service")
	flag.BoolVar(&restartFlag, "restart", false, "restart a service")
	flag.BoolVar(&statusFlag, "status", false, "see this status of a service")
	flag.BoolVar(&listFlag, "list", false, "list the available services")
	flag.StringVar(&serviceFlag, "service", "none", "service you wish to take action on")

	// parse the flags
	flag.Parse()

	// check for an initial command
	if flag.NArg() == 0 {
		usage()
	}

	// Print output
	fmt.Printf("flags:\n")
	fmt.Printf("\th=%v\n", helpShort)
	fmt.Printf("\thelp=%v\n", helpLong)
	fmt.Printf("\trestart=%v\n", restartFlag)
	fmt.Printf("\tstart=%v\n", startFlag)
	fmt.Printf("\tstatus=%v\n", statusFlag)
	fmt.Printf("\tstop=%v\n", stopFlag)
}

func prog2() {
	// Subcommands
	countCommand := flag.NewFlagSet("count", flag.ExitOnError)
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)

	// Count subcommand flag pointers
	// Adding a new choice for --metric of 'substring' and a new --substring flag
	countTextPtr := countCommand.String("text", "", "Text to parse. (Required)")
	countMetricPtr := countCommand.String("metric", "chars", "Metric {chars|words|lines|substring}. (Required)")
	countSubstringPtr := countCommand.String(
		"substring", "", "The substring to be counted. Required for --metric=substring",
	)
	countUniquePtr := countCommand.Bool("unique", false, "Measure unique values of a metric.")

	// List subcommand flag pointers
	listTextPtr := listCommand.String("text", "", "Text to parse. (Required)")
	listMetricPtr := listCommand.String("metric", "chars", "Metric <chars|words|lines>. (Required)")
	listUniquePtr := listCommand.Bool("unique", false, "Measure unique values of a metric.")

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		fmt.Println("list or count subcommand is required")
		os.Exit(1)
	}

	// Switch on the subcommand
	// Parse the flags for appropriate FlagSet
	// FlagSet.Parse() requires a set of arguments to parse as input
	// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
	switch os.Args[1] {
	case "list":
		listCommand.Parse(os.Args[2:])
	case "count":
		countCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Check which subcommand was Parsed using the FlagSet.Parsed() function. Handle each case accordingly.
	// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
	if listCommand.Parsed() {
		// Required Flags
		if *listTextPtr == "" {
			listCommand.PrintDefaults()
			os.Exit(1)
		}
		// Choice flag
		metricChoices := map[string]bool{"chars": true, "words": true, "lines": true}
		if _, validChoice := metricChoices[*listMetricPtr]; !validChoice {
			listCommand.PrintDefaults()
			os.Exit(1)
		}
		// Print
		fmt.Printf("textPtr: %s, metricPtr: %s, uniquePtr: %t\n", *listTextPtr, *listMetricPtr, *listUniquePtr)
	}

	if countCommand.Parsed() {
		// Required Flags
		if *countTextPtr == "" {
			countCommand.PrintDefaults()
			os.Exit(1)
		}
		// If the metric flag is substring, the substring flag is required
		if *countMetricPtr == "substring" && *countSubstringPtr == "" {
			countCommand.PrintDefaults()
			os.Exit(1)
		}
		// If the metric flag is not substring, the substring flag must not be used
		if *countMetricPtr != "substring" && *countSubstringPtr != "" {
			fmt.Println("--substring may only be used with --metric=substring.")
			countCommand.PrintDefaults()
			os.Exit(1)
		}
		// Choice flag
		metricChoices := map[string]bool{"chars": true, "words": true, "lines": true, "substring": true}
		if _, validChoice := metricChoices[*listMetricPtr]; !validChoice {
			countCommand.PrintDefaults()
			os.Exit(1)
		}
		// Print
		fmt.Printf(
			"textPtr: %s, metricPtr: %s, substringPtr: %v, uniquePtr: %t\n", *countTextPtr, *countMetricPtr,
			*countSubstringPtr, *countUniquePtr,
		)
	}

}
