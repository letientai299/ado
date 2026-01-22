package profiling

import (
	"os"
	"runtime"
	"runtime/pprof"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type StopFn = func()

const flagName = "pprof-profiles"

var supportedProfiles = []string{
	"cpu",
	"block",
	"mutex",
	"goroutine",
	"threadcreate",
	"allocs",
	"heap",
}

func NoOp() {}

func RegisterFlag(cmd *cobra.Command) {
	all := strings.Join(supportedProfiles, ", ")
	cmd.PersistentFlags().
		StringSlice(flagName, nil, "collect profiles: "+all)
}

// Start initializes and starts the requested pprof profiles.
// It returns a function that should be deferred to stop profiling and flush data.
func Start(cmd *cobra.Command) func() {
	profiles := getProfiles(cmd)
	if len(profiles) == 0 {
		return NoOp
	}

	var stops []func()
	for _, p := range profiles {
		filename := getFileName(p)
		switch p {
		case "cpu":
			f, err := os.Create(filename)
			if err != nil {
				log.Errorf("could not create CPU profile: %v", err)
				continue
			}
			if err = pprof.StartCPUProfile(f); err != nil {
				log.Errorf("could not start CPU profile: %v", err)
				_ = f.Close()
				continue
			}
			stops = append(stops, func() {
				pprof.StopCPUProfile()
				_ = f.Close()
				log.Infof("CPU profile saved to %s", filename)
			})
		case "heap", "allocs", "block", "mutex", "goroutine", "threadcreate":
			// These are all lookup-based profiles
			stops = append(stops, func() {
				f, err := os.Create(filename)
				if err != nil {
					log.Errorf("could not create %s profile: %v", p, err)
					return
				}
				defer func(f *os.File) { _ = f.Close() }(f)

				// Write heap profile
				runtime.GC() // get up-to-date statistics
				if err = pprof.Lookup(p).WriteTo(f, 0); err != nil {
					log.Errorf("could not write %s profile: %v", p, err)
					return
				}
				log.Infof("%s profile saved to %s", p, filename)
			})
		default:
			log.Warnf("unknown profile type: %s", p)
		}
	}

	return func() {
		for i := len(stops) - 1; i >= 0; i-- {
			stops[i]()
		}
	}
}

func getFileName(p string) string {
	// use the same prefix to keep those output files near each other
	return "_" + p + ".pprof"
}

func getProfiles(cmd *cobra.Command) []string {
	if cmd == nil || !cmd.Flags().Changed(flagName) {
		return nil
	}

	profiles, err := cmd.Flags().GetStringSlice(flagName)
	if err != nil {
		log.Warnf("could not parse %s flag: %v", flagName, err)
		return nil
	}

	return slices.DeleteFunc(profiles, func(p string) bool {
		return !slices.Contains(supportedProfiles, p)
	})
}
