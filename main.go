package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

type GitWrapper struct{}

func (g *GitWrapper) RunCommand(args ...string) (strings.Builder, int) {

	fmt.Println("Running command:", "git", args)

	if len(args) > 0 && args[0] == "log" {
		args = append(args, "--pretty=format:%ai|%H|%an|%s")
	}
	cmd := exec.Command("git", args...)
	cmd.Env = append(cmd.Env, "LANG=en_US.UTF-8") // Set the LANG environment variable
	stdoutStderr, _ := cmd.CombinedOutput()

	if len(args) < 1 {
		args = append(args, "")
	}

	return FormatOutput(args[0], string(stdoutStderr)), cmd.ProcessState.ExitCode()
}

func FormatOutput(command, output string) strings.Builder {
	lines := strings.Split(output, "\n")
	var formattedOutput strings.Builder

	switch command {
	case "status":
		tableData := [][]string{{"FILE", "STATUS"}}
		otherLines := []string{}
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "modified:") {
				tableData = append(tableData, []string{strings.TrimPrefix(line, "modified: "), pterm.FgBlue.Sprint("Modified")})
			} else if strings.HasPrefix(line, "new file:") {
				tableData = append(tableData, []string{strings.TrimPrefix(line, "new file: "), pterm.FgGreen.Sprint("New File")})
			} else if strings.HasPrefix(line, "deleted:") {
				tableData = append(tableData, []string{strings.TrimPrefix(line, "deleted: "), pterm.FgRed.Sprint("Deleted")})
			} else {
				line = resolveLines(line)
				otherLines = append(otherLines, line)
			}
		}
		formattedOutput.WriteString(strings.Join(otherLines, ""))
		output, _ := pterm.DefaultTable.WithHasHeader().WithData(tableData).Srender()
		formattedOutput.WriteString(output)

	case "log":
		tableData := [][]string{{"DATE", "HASH", "AUTHOR", "MESSAGE"}}
		for _, line := range lines {
			parts := strings.Split(line, "|")
			if len(parts) < 4 {
				continue
			}
			date, hash, author, message := parts[0], parts[1], parts[2], parts[3]
			tableData = append(tableData, []string{date, hash[:10], author, message})
		}
		output, _ := pterm.DefaultTable.WithHasHeader().WithData(tableData).Srender()
		formattedOutput.WriteString(output)
	case "branch":
		tableData := [][]string{{"BRANCH", "CURRENT"}}
		for _, line := range lines {
			if strings.HasPrefix(line, "*") {
				line = strings.Replace(line, "*", " ", 1)
				tableData = append(tableData, []string{pterm.FgGreen.Sprint(line), "📍"})
			} else if line != "" {
				tableData = append(tableData, []string{line, ""})
			}
		}
		output, _ := pterm.DefaultTable.WithHasHeader().WithData(tableData).Srender()
		formattedOutput.WriteString(output)
	case "diff":
		for _, line := range lines {
			switch {
			case strings.HasPrefix(line, "index"):
				formattedOutput.WriteString(pterm.FgYellow.Sprintln(line))
			case strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++"):
				formattedOutput.WriteString(pterm.FgCyan.Sprintln(line))
			case strings.HasPrefix(line, "@@"):
				formattedOutput.WriteString(pterm.FgBlue.Sprintln(line))
			case strings.HasPrefix(line, "-"):
				formattedOutput.WriteString(pterm.FgRed.Sprintln(line))
			case strings.HasPrefix(line, "+"):

				formattedOutput.WriteString(pterm.FgGreen.Sprintln(line))
			default:
				formattedOutput.WriteString(pterm.FgGray.Sprintln(line))
			}
		}
	default:

		var infoLines strings.Builder
		for _, line := range lines {

			if strings.HasPrefix(line, "hint:") {
				line = strings.Replace(line, "hint:", "", 1)
				infoLines.WriteString(line + "\n")
				continue
			}
			if infoLines.Len() > 0 {
				formattedOutput.WriteString(pterm.Info.Sprintln(infoLines.String()))
				infoLines.Reset()
			}

			line = resolveLines(line)
			formattedOutput.WriteString(line)
		}

	}

	return formattedOutput
}

func resolveLines(line string) string {

	if strings.HasPrefix(line, "On branch") {
		line = pterm.DefaultBox.Sprintf(pterm.Green(line))
	} else if strings.Contains(line, "modified:") {
		line = "🛠️ " + pterm.LightBlue(line)
	} else if strings.Contains(line, "new file:") {
		line = "✨" + pterm.LightGreen(line)
	} else if strings.Contains(line, "deleted:") {
		line = "❌ " + pterm.LightRed(line)
	} else if strings.Contains(line, "Untracked files:") {
		line = "📁 " + line
	}

	if strings.HasPrefix(line, "usage: git") {
		line = strings.Replace(line, "usage: git", "usage: fit", 1)
	}
	if strings.HasPrefix(line, "fatal:") {
		line = strings.Replace(line, "fatal:", "", 1)
		line = pterm.Error.Sprintln(line)
	}

	if strings.HasPrefix(line, "error:") {
		line = strings.Replace(line, "error:", "", 1)
		line = pterm.Error.Sprintln(line)
	}

	if strings.HasPrefix(line, "warning:") {
		line = strings.Replace(line, "warning:", "", 1)
		line = pterm.Warning.Sprintln(line)
	}

	if strings.HasPrefix(line, "No commits yet") {
		line = pterm.Info.Sprintln(line)
	}

	if strings.HasPrefix(line, "nothing added to commit but untracked files present") {
		line = pterm.Info.Sprintln(line)
	}

	if strings.HasPrefix(line, "no changes added to commit") {
		line = pterm.Info.Sprintln(line)
	}

	if strings.HasPrefix(line, "nothing to commit, working tree clean") {
		line = pterm.Info.Sprintln(line)
	}

	if strings.HasPrefix(line, "HEAD") {
		line = pterm.Info.Sprintln(line)
	}

	if strings.HasPrefix(line, "Initialized empty Git repository") {
		line = pterm.DefaultBox.Sprintln(pterm.LightGreen("✅ " + line))
	}

	if strings.HasPrefix(line, "Switched to") {
		line = pterm.DefaultBox.Sprintln(pterm.LightGreen("✅ " + line))
	}

	if strings.HasPrefix(line, "Already on") {
		line = pterm.DefaultBox.Sprintln(pterm.LightBlue("ℹ️ " + line))
	}

	if strings.HasPrefix(line, "Your branch is") {
		line = pterm.Info.Sprintln(line)
	}

	return strings.TrimRight(line, "\n") + "\n"
}

func main() {
	pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("F", pterm.FgLightBlue.ToStyle()),
		putils.LettersFromStringWithStyle("it", pterm.FgYellow.ToStyle()),
	).Render()

	git := &GitWrapper{}
	args := os.Args[1:]

	output, exitCode := git.RunCommand(args...)

	fmt.Println(output.String())
	os.Exit(exitCode)
}
