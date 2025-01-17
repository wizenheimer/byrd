package main

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

var (
	secretKey    string
	rotationTime time.Duration
	tokenManager *utils.TokenManager
	watchMode    bool
	compactMode  bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "tokenctl",
		Short: "A CLI tool for managing rotating bearer tokens",
		Run:   run,
	}

	rootCmd.PersistentFlags().StringVarP(&secretKey, "secret", "s", "", "Secret key for token generation")
	rootCmd.PersistentFlags().DurationVarP(&rotationTime, "rotation", "r", 5*time.Minute, "Token rotation interval")
	rootCmd.PersistentFlags().BoolVarP(&watchMode, "watch", "w", false, "Watch mode - continuously display token")
	rootCmd.PersistentFlags().BoolVarP(&compactMode, "compact", "c", false, "Compact display mode")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	if secretKey == "" {
		fmt.Println("Error: Secret key is required")
		os.Exit(1)
	}

	tokenManager = utils.NewTokenManager(secretKey, rotationTime)

	if watchMode {
		watchTokens()
	} else {
		displayToken()
	}
}

func displayToken() {
	token := tokenManager.GenerateToken()
	remaining := getRemainingTime()

	if compactMode {
		displayCompact(token, remaining)
	} else {
		displayFull(token, remaining)
	}
}

func watchTokens() {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "Monitoring tokens "
	s.Start()

	for {
		clearScreen()
		s.Stop()
		displayToken()
		time.Sleep(1 * time.Second)
		s.Start()
	}
}

func displayFull(token string, remaining time.Duration) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Token", "Status", "Expires In", "Rotation Interval"})
	table.SetAutoWrapText(false)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	status := color.New(color.FgCyan).Sprintf("Valid")
	if remaining.Seconds() < 30 {
		status = color.New(color.FgYellow).Sprintf("Rotating Soon")
	}

	table.Append([]string{
		token,
		status,
		fmt.Sprintf("%.0fs", remaining.Seconds()),
		fmt.Sprintf("%.0fm", rotationTime.Minutes()),
	})

	table.Render()
}

func displayCompact(token string, remaining time.Duration) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Token", "TTL"})
	table.SetAutoWrapText(false)
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColumnSeparator(" ")
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Append([]string{
		token,
		fmt.Sprintf("%.0fs", remaining.Seconds()),
	})

	table.Render()
}

func getRemainingTime() time.Duration {
	currentInterval := tokenManager.GetCurrentInterval()
	rotationTime := tokenManager.GetRotationTime()
	nextRotation := (currentInterval + 1) * int64(rotationTime.Seconds())
	return time.Duration(nextRotation-time.Now().Unix()) * time.Second
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
