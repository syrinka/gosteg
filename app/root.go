package app

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gosteg",
	Short: "A mimic of zsteg",
	Run: func(cmd *cobra.Command, args []string) {
		dbg, _ := cmd.Flags().GetBool("debug")
		if dbg {
			logrus.SetLevel(logrus.DebugLevel)
		}
		tra, _ := cmd.Flags().GetBool("trace")
		if tra {
			logrus.SetLevel(logrus.TraceLevel)
		}

		if len(args) != 1 {
			panic("need input path")
		}
		var ipath = args[0]

		opath, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}
		if opath == "<nil>" {
			opath = ipath + ".o"
		}

		channel, err := cmd.Flags().GetString("channel")
		if err != nil {
			panic(err)
		}
		channel = strings.ToLower(channel)
		if !strings.Contains("rgba", channel) {
			panic("invalid channel")
		}

		bits, _ := cmd.Flags().GetUintSlice("bits")
		order, _ := cmd.Flags().GetString("order")
		xy, _ := cmd.Flags().GetString("xy")

		var opt = StegOption{
			channel: channel,
			bits:    bits,
			order:   order,
			xy:      xy,
		}
		logrus.Debug(opt)

		f, err := os.Open(ipath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		var data = extractData(f, opt)
		invert, _ := cmd.Flags().GetBool("invert")
		if invert {
			for i := range len(data) {
				data[i] = ^data[i]
			}
		}

		vis, _ := cmd.Flags().GetBool("visualize")
		if vis {
			visualize(data)
		}

		fo, err := os.OpenFile(opath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer fo.Close()

		_, err = fo.Write(data)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.Flags().StringP("output", "o", "<nil>", "output path")
	rootCmd.Flags().StringP("channel", "c", "<nil>", "selected channel")
	rootCmd.Flags().UintSliceP("bits", "b", []uint{1}, "bits")
	rootCmd.Flags().StringP("order", "s", "lsb", "bit order, lsb or msb")
	rootCmd.Flags().StringP("xy", "x", "xy", "determine scan dimension")
	rootCmd.Flags().BoolP("invert", "v", false, "invert result")

	rootCmd.Flags().Bool("visualize", false, "visualize extracted data")

	rootCmd.Flags().BoolP("debug", "d", false, "debug level")
	rootCmd.Flags().BoolP("trace", "D", false, "trace level")

	rootCmd.MarkFlagRequired("channel")
}

func Execute() {
	rootCmd.Execute()
}
