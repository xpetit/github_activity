package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	. "github.com/xpetit/x/v2"
)

const (
	width  = 52
	height = 7
)

func main() {
	repoPath := flag.String("repo", "empty", "Path to the GitHub repository")
	filename := flag.String("f", "grid.txt", "Path to the text file")
	nbCommits := flag.Int("commits", 67, `Number of commits per active "pixel"`)
	usage := flag.Usage
	flag.Usage = func() {
		usage()
		fmt.Println(`Create a GitHub repository and clone it next to this program`)
	}
	flag.Parse()

	lines := strings.Split(string(C2(os.ReadFile(*filename))), "\n")

	git := func(env []string, args ...string) string {
		cmd := exec.Command("git", append([]string{"-C", *repoPath}, args...)...)
		cmd.Env = append(os.Environ(), append(env, "TZ=etc/UTC")...)
		return strings.TrimSpace(string(C2(Output(cmd))))
	}
	origin := git(nil, "remote", "get-url", "origin")
	branch := git(nil, "branch", "--show-current")

	C(os.RemoveAll(*repoPath))
	C(os.Mkdir(*repoPath, 0o755))
	git(nil, "init")
	git(nil, "remote", "add", "origin", origin)

	year, month, day := time.Now().UTC().Date()
	date := time.Date(year, month, day, 12, 0, 0, 0, time.UTC)
	date = date.AddDate(0, 0, -int(date.Weekday()+1))

	fmt.Println("Creating commits...")
	for x := width - 1; x >= 0; x-- {
		for y := height - 1; y >= 0; y-- {
			if y < len(lines) && x < len(lines[y]) && lines[y][x] == '#' {
				s := date.Format("2006-01-02T15:04:05")
				for i := 0; i < *nbCommits; i++ {
					git(
						[]string{"GIT_AUTHOR_DATE=" + s, "GIT_COMMITTER_DATE=" + s},
						"commit", "--allow-empty", "--allow-empty-message", "-m", "",
					)
				}
			}
			date = date.AddDate(0, 0, -1)
		}
	}
	git(nil, "gc", "--aggressive", "--prune=now")
	git(nil, "push", "-f", "origin", branch)
}
