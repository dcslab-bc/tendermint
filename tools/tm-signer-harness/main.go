package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/Finschia/ostracon/crypto/ed25519"
	"github.com/Finschia/ostracon/libs/log"
	"github.com/Finschia/ostracon/privval"
	"github.com/Finschia/ostracon/tools/tm-signer-harness/internal"
	"github.com/Finschia/ostracon/version"
)

const (
	defaultAcceptRetries    = 100
	defaultBindAddr         = "tcp://127.0.0.1:0"
	defaultOCHome           = "~/.ostracon"
	defaultAcceptDeadline   = 1
	defaultConnDeadline     = 3
	defaultExtractKeyOutput = "./signing.key"
)

var logger = log.NewOCLogger(log.NewSyncWriter(os.Stdout))

// Command line flags
var (
	flagAcceptRetries int
	flagBindAddr      string
	flagOCHome        string
	flagKeyOutputPath string
)

// Command line commands
var (
	rootCmd       *flag.FlagSet
	runCmd        *flag.FlagSet
	extractKeyCmd *flag.FlagSet
	versionCmd    *flag.FlagSet
)

func init() {
	rootCmd = flag.NewFlagSet("root", flag.ExitOnError)
	rootCmd.Usage = func() {
		fmt.Println(`Remote signer test harness for Ostracon.

Usage:
  tm-signer-harness <command> [flags]

Available Commands:
  extract_key        Extracts a signing key from a local Ostracon instance
  help               Help on the available commands
  run                Runs the test harness
  version            Display version information and exit

Use "tm-signer-harness help <command>" for more information about that command.`)
		fmt.Println("")
	}

	runCmd = flag.NewFlagSet("run", flag.ExitOnError)
	runCmd.IntVar(&flagAcceptRetries,
		"accept-retries",
		defaultAcceptRetries,
		"The number of attempts to listen for incoming connections")
	runCmd.StringVar(&flagBindAddr, "addr", defaultBindAddr, "Bind to this address for the testing")
	runCmd.StringVar(&flagOCHome, "ochome", defaultOCHome, "Path to the Ostracon home directory")
	runCmd.Usage = func() {
		fmt.Println(`Runs the remote signer test harness for Ostracon.

Usage:
  tm-signer-harness run [flags]

Flags:`)
		runCmd.PrintDefaults()
		fmt.Println("")
	}

	extractKeyCmd = flag.NewFlagSet("extract_key", flag.ExitOnError)
	extractKeyCmd.StringVar(&flagKeyOutputPath,
		"output",
		defaultExtractKeyOutput,
		"Path to which signing key should be written")
	extractKeyCmd.StringVar(&flagOCHome, "ochome", defaultOCHome, "Path to the Ostracon home directory")
	extractKeyCmd.Usage = func() {
		fmt.Println(`Extracts a signing key from a local Ostracon instance for use in the remote
signer under test.

Usage:
  tm-signer-harness extract_key [flags]

Flags:`)
		extractKeyCmd.PrintDefaults()
		fmt.Println("")
	}

	versionCmd = flag.NewFlagSet("version", flag.ExitOnError)
	versionCmd.Usage = func() {
		fmt.Println(`
Prints the Ostracon version for which this remote signer harness was built.

Usage:
  tm-signer-harness version`)
		fmt.Println("")
	}
}

func runTestHarness(acceptRetries int, bindAddr, ochome string) {
	ochome = internal.ExpandPath(ochome)
	cfg := internal.TestHarnessConfig{
		BindAddr:         bindAddr,
		KeyFile:          filepath.Join(ochome, "config", "priv_validator_key.json"),
		StateFile:        filepath.Join(ochome, "data", "priv_validator_state.json"),
		GenesisFile:      filepath.Join(ochome, "config", "genesis.json"),
		AcceptDeadline:   time.Duration(defaultAcceptDeadline) * time.Second,
		AcceptRetries:    acceptRetries,
		ConnDeadline:     time.Duration(defaultConnDeadline) * time.Second,
		SecretConnKey:    ed25519.GenPrivKey(),
		ExitWhenComplete: true,
	}
	harness, err := internal.NewTestHarness(logger, cfg)
	if err != nil {
		logger.Error(err.Error())
		if therr, ok := err.(*internal.TestHarnessError); ok {
			os.Exit(therr.Code)
		}
		os.Exit(internal.ErrOther)
	}
	harness.Run()
}

func extractKey(ochome, outputPath string) {
	keyFile := filepath.Join(internal.ExpandPath(ochome), "config", "priv_validator_key.json")
	stateFile := filepath.Join(internal.ExpandPath(ochome), "data", "priv_validator_state.json")
	fpv := privval.LoadFilePV(keyFile, stateFile)
	pkb := []byte(fpv.Key.PrivKey.(ed25519.PrivKey))
	if err := ioutil.WriteFile(internal.ExpandPath(outputPath), pkb[:32], 0600); err != nil {
		logger.Info("Failed to write private key", "output", outputPath, "err", err)
		os.Exit(1)
	}
	logger.Info("Successfully wrote private key", "output", outputPath)
}

func main() {
	if err := rootCmd.Parse(os.Args[1:]); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}
	if rootCmd.NArg() == 0 || (rootCmd.NArg() == 1 && rootCmd.Arg(0) == "help") {
		rootCmd.Usage()
		os.Exit(0)
	}

	logger = log.NewFilter(logger, log.AllowInfo())

	switch rootCmd.Arg(0) {
	case "help":
		switch rootCmd.Arg(1) {
		case "run":
			runCmd.Usage()
		case "extract_key":
			extractKeyCmd.Usage()
		case "version":
			versionCmd.Usage()
		default:
			fmt.Printf("Unrecognized command: %s\n", rootCmd.Arg(1))
			os.Exit(1)
		}
	case "run":
		if err := runCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("Error parsing flags: %v\n", err)
			os.Exit(1)
		}
		runTestHarness(flagAcceptRetries, flagBindAddr, flagOCHome)
	case "extract_key":
		if err := extractKeyCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("Error parsing flags: %v\n", err)
			os.Exit(1)
		}
		extractKey(flagOCHome, flagKeyOutputPath)
	case "version":
		fmt.Println(version.OCCoreSemVer)
	default:
		fmt.Printf("Unrecognized command: %s\n", flag.Arg(0))
		os.Exit(1)
	}
}
