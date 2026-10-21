package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/arseniizyk/CoCMetaAnalyser/internal/service"
)

type CommandLine struct {
	svc           service.Service
	season        *string
	limit         *int
	filename      *string
	outDir        *string
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
  all -season SEASON -limit LIMIT -filename FILENAME [-outdir DIR]
      Outputs both meta and metapairs in a single API pass to <outdir>/{meta,metapairs}/SEASON/FILENAME.json
  meta -season SEASON -limit LIMIT -filename FILENAME [-outdir DIR]
      Outputs meta in JSON format to <outdir>/meta/SEASON/FILENAME.json
  metapairs -season SEASON -limit LIMIT -filename FILENAME [-outdir DIR]
      Outputs meta pairs in JSON format to <outdir>/metapairs/SEASON/FILENAME.json
  seasons
      Lists all available Legendary League seasons from CoC API
  backfill [-from SEASON] [-limit LIMIT] [-filename FILENAME] [-outdir DIR]
      Automatically discovers and collects missing seasons
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

	allCmd := flag.NewFlagSet("all", flag.ExitOnError)
	metaCmd := flag.NewFlagSet("meta", flag.ExitOnError)
	metaPairsCmd := flag.NewFlagSet("metapairs", flag.ExitOnError)

	switch os.Args[1] {
	case "all":
		cli.initCmdParams(allCmd)
		cli.parseCmdParams(allCmd)
		if err := cli.processAll(); err != nil {
			return err
		}
		fmt.Println("Successfully collected all meta data!")
	case "meta":
		cli.initCmdParams(metaCmd)
		cli.parseCmdParams(metaCmd)
		if err := cli.processMeta(); err != nil {
			return err
		}
		fmt.Println("Successfully collected meta items!")
	case "metapairs":
		cli.initCmdParams(metaPairsCmd)
		cli.parseCmdParams(metaPairsCmd)
		if err := cli.processMetaPairs(); err != nil {
			return err
		}
		fmt.Println("Successfully collected meta pairs!")
	case "seasons":
		return cli.processSeasons()
	case "backfill":
		return cli.processBackfill()
	default:
		cli.printUsage()
		return fmt.Errorf("unknown command: %s", os.Args[1])
	}

	return nil
}

func (cli *CommandLine) processAll() error {
	res, err := cli.svc.GetMetaFull(*cli.limit, *cli.season)
	if err != nil {
		log.Println("Something went wrong, err:", err)
		return err
	}

	if err := writeJSONToFile(cli.metaPath, res.ItemsTotal); err != nil {
		return err
	}
	if err := writeJSONToFile(cli.metaPairsPath, res.PairsTotal); err != nil {
		return err
	}

	// If collecting 10k or large sample, automatically save top-100 files as well
	if *cli.limit >= 100 && *cli.filename != "meta100" && len(res.ItemsTop100) > 0 {
		baseDir := resolveBaseDir(*cli.outDir)
		cleanSeason := NormalizeSeason(*cli.season)
		meta100Path := filepath.Join(baseDir, "meta", cleanSeason, "meta100.json")
		pairs100Path := filepath.Join(baseDir, "metapairs", cleanSeason, "meta100.json")

		_ = writeJSONToFile(meta100Path, res.ItemsTop100)
		_ = writeJSONToFile(pairs100Path, res.PairsTop100)
		log.Printf("Also automatically generated Top-100 files for season %s", cleanSeason)
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

func resolveBaseDir(customDir string) string {
	if customDir != "" {
		return customDir
	}
	candidates := []string{
		"frontend/public",
		"../../frontend/public",
		"../frontend/public",
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			return c
		}
	}
	return "frontend/public"
}

// NormalizeSeason converts various CoC season ID formats into a clean YYYY-MM representation.
func NormalizeSeason(raw string) string {
	raw = strings.TrimSpace(raw)
	if matched, _ := regexp.MatchString(`^\d{4}-\d{2}$`, raw); matched {
		return raw
	}

	re := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`)
	m := re.FindStringSubmatch(raw)
	if len(m) == 4 {
		year, _ := strconv.Atoi(m[1])
		month, _ := strconv.Atoi(m[2])
		day, _ := strconv.Atoi(m[3])

		if strings.HasPrefix(raw, "v2-") {
			if day >= 28 {
				month++
				if month > 12 {
					month = 1
					year++
				}
			}
			return fmt.Sprintf("%04d-%02d", year, month)
		}

		if day <= 7 {
			month--
			if month == 0 {
				month = 12
				year--
			}
		}
		return fmt.Sprintf("%04d-%02d", year, month)
	}

	return raw
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir() && fi.Size() > 0
}

func (cli *CommandLine) processSeasons() error {
	seasons, err := cli.svc.GetLeagueSeasons()
	if err != nil {
		return fmt.Errorf("fetching league seasons: %w", err)
	}

	fmt.Println("Available Legendary League seasons:")
	for _, s := range seasons.Items {
		fmt.Println(" -", s.ID)
	}
	return nil
}

func (cli *CommandLine) processBackfill() error {
	backfillCmd := flag.NewFlagSet("backfill", flag.ExitOnError)
	limit := backfillCmd.Int("limit", 10000, "Number of players (min 100, max 25k)")
	filename := backfillCmd.String("filename", "meta10k", "Name for output file (without .JSON)")
	outDir := backfillCmd.String("outdir", "", "Path to frontend/public directory (optional)")
	fromSeason := backfillCmd.String("from", "2025-01", "Start season filter (optional, e.g. 2025-01)")

	if len(os.Args) > 2 {
		if err := backfillCmd.Parse(os.Args[2:]); err != nil {
			return err
		}
	}

	baseDir := resolveBaseDir(*outDir)
	seasons, err := cli.svc.GetLeagueSeasons()
	if err != nil {
		return fmt.Errorf("fetching seasons: %w", err)
	}

	log.Printf("Found %d total seasons in CoC API. Checking for missing seasons starting from %s...", len(seasons.Items), *fromSeason)

	for _, s := range seasons.Items {
		if *fromSeason != "" && s.ID < *fromSeason {
			continue
		}

		cleanSeason := NormalizeSeason(s.ID)
		metaFile := filepath.Join(baseDir, "meta", cleanSeason, *filename+".json")
		pairsFile := filepath.Join(baseDir, "metapairs", cleanSeason, *filename+".json")

		if fileExists(metaFile) && fileExists(pairsFile) {
			log.Printf("Season %s (%s) already collected, skipping", cleanSeason, s.ID)
			continue
		}

		log.Printf("Collecting missing data for season %s (API: %s, limit: %d)...", cleanSeason, s.ID, *limit)
		res, err := cli.svc.GetMetaFull(*limit, s.ID)
		if err != nil {
			log.Printf("Failed to collect season %s: %v", s.ID, err)
			continue
		}

		if err := writeJSONToFile(metaFile, res.ItemsTotal); err != nil {
			return err
		}
		if err := writeJSONToFile(pairsFile, res.PairsTotal); err != nil {
			return err
		}

		// Also auto-save top 100 if collecting a larger sample
		if *limit >= 100 && *filename != "meta100" && len(res.ItemsTop100) > 0 {
			meta100File := filepath.Join(baseDir, "meta", cleanSeason, "meta100.json")
			pairs100File := filepath.Join(baseDir, "metapairs", cleanSeason, "meta100.json")
			_ = writeJSONToFile(meta100File, res.ItemsTop100)
			_ = writeJSONToFile(pairs100File, res.PairsTop100)
		}

		log.Printf("Successfully saved season %s to %s and %s", cleanSeason, metaFile, pairsFile)
	}

	fmt.Println("Backfill completed!")
	return nil
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

	baseDir := resolveBaseDir(*cli.outDir)
	cleanSeason := NormalizeSeason(*cli.season)
	cli.metaPath = filepath.Join(baseDir, "meta", cleanSeason, *cli.filename+".json")
	cli.metaPairsPath = filepath.Join(baseDir, "metapairs", cleanSeason, *cli.filename+".json")
}

func (cli *CommandLine) initCmdParams(cmd *flag.FlagSet) {
	cli.season = cmd.String("season", "", `Season in format yyyy-mm(2025-05)`)
	cli.limit = cmd.Int("limit", 0, "Number of players(min 100, max 25k)")
	cli.filename = cmd.String("filename", "", "Name for output file(without .JSON)")
	cli.outDir = cmd.String("outdir", "", "Path to frontend/public directory (optional, auto-detected)")
}
