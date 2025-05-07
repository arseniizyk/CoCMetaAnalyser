package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/arseniizyk/CoCMetaAnalyser/internal/service"
)

type CommandLine struct {
	svc           service.Service
	season        *string
	limit         *int
	filename      *string
	metaPath      string
	metaPairsPath string
}

func New(svc service.Service) *CommandLine {
	return &CommandLine{
		svc: svc,
	}
}

func (cli *CommandLine) printUsage() {
	text := `
Usage:
  meta -season SEASON -limit LIMIT -filename FILENAME - outputs meta in JSON format to cmd/app/meta/SEASON/FILENAME.json
  metapairs -season SEASON -limit LIMIT -filename FILENAME - outputs meta pairs in JSON format to cmd/app/metapairs/SEASON/FILENAME.json
`
	fmt.Println(text)
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CommandLine) Run() error {
	cli.validateArgs()

	metaCmd := flag.NewFlagSet("meta", flag.ExitOnError)
	metaPairsCmd := flag.NewFlagSet("metapairs", flag.ExitOnError)

	switch os.Args[1] {
	case "meta":
		cli.initCmdParams(metaCmd)
		cli.parseCmdParams(metaCmd)
		if err := cli.processMeta(); err != nil {
			return err
		}
		fmt.Println("Successfully!")
	case "metapairs":
		cli.initCmdParams(metaPairsCmd)
		cli.parseCmdParams(metaPairsCmd)
		if err := cli.processMetaPairs(); err != nil {
			return err
		}
		fmt.Println("Successfully!")
	}

	return nil
}

func (cli *CommandLine) processMeta() error {
	meta, err := cli.svc.GetMetaItems(*cli.limit, *cli.season)
	if err != nil {
		log.Println("Something went wrong, err:", err)
		return err
	}

	return writeJSONToFile(cli.metaPath, meta)
}

func (cli *CommandLine) processMetaPairs() error {
	meta, err := cli.svc.GetMetaItemPairs(*cli.limit, *cli.season)
	if err != nil {
		log.Println("Something went wrong, err:", err)
		return err
	}

	return writeJSONToFile(cli.metaPairsPath, meta)
}

func writeJSONToFile(path string, data any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(formatted)
	return err
}

func (cli *CommandLine) parseCmdParams(cmd *flag.FlagSet) {
	if err := cmd.Parse(os.Args[2:]); err != nil {
		cli.printUsage()
		os.Exit(1)
	}

	if cmd.Parsed() {
		if *cli.season == "" || *cli.limit == 0 || *cli.filename == "" {
			cli.printUsage()
			os.Exit(1)
		}
	}
	cli.metaPath = fmt.Sprintf("meta/%s/%s.json", *cli.season, *cli.filename)
	cli.metaPairsPath = fmt.Sprintf("metapairs/%s/%s.json", *cli.season, *cli.filename)
}

func (cli *CommandLine) initCmdParams(cmd *flag.FlagSet) {
	cli.season = cmd.String("season", "", `Season in format yyyy-mm(2025-05)`)
	cli.limit = cmd.Int("limit", 0, "Number of players(min 100, max 25k)")
	cli.filename = cmd.String("filename", "", "Name for output file(without .JSON)")
}
