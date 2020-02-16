package build

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// filled in during build
var (
	GitTag string
)

type Target struct {
	System string
	Arch   string
}

func Process(appName, mainPath string) {
	app := flag.String("app", appName, "the name of the app to build")
	system := flag.String("os", runtime.GOOS, "the OS to build for ('go tool dist list -json' for options)")
	arch := flag.String("arch", runtime.GOARCH, "the architecture to build for ('go tool dist list -json' for options)")
	suffix := flag.String("suffix", ".exe", "the filename suffix to use (\"\" for no suffix)")
	flag.Parse()

	PerformBuild(*app, *system, *arch, *suffix, mainPath)
}

func ProcessForTargets(appName string, targets []Target, suffix, mainPath string) {
	for _, t := range targets {
		PerformBuild(appName, t.System, t.Arch, suffix, mainPath)
	}
}

func ProcessDefaultTargets(appName, mainPath string) {
	ProcessForTargets(appName, []Target{
		{
			System: "windows",
			Arch:   "amd64",
		},
		{
			System: "darwin",
			Arch:   "amd64",
		},
		{
			System: "linux",
			Arch:   "amd64",
		},
	}, ".exe", mainPath)
}

func GetGitTag() (string, error) {
	gitTagCmd := exec.Command("git", "describe", "--always", "--tag")
	gitStdOut, _ := tieCommandToBuffers(gitTagCmd)
	if err := gitTagCmd.Run(); err != nil {
		return "", errors.New("while running git describe, error = " + err.Error())
	}
	return strings.TrimSpace(gitStdOut.String()), nil
}

func PerformBuild(app, system, arch, suffix, mainPath string) {
	gitTag, err := GetGitTag()
	if err != nil {
		log.Fatal("while running git describe, error = " + err.Error())
	}

	cmd := getBuildCommand(system, app, arch, suffix, gitTag, mainPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("built " + app + suffix + " v" + gitTag + " for OS \"" + system + "\" with arch \"" + arch + "\"")
}

func getBuildCommand(platform, appName, arch, suffix, gitTag, mainPath string) *exec.Cmd {
	cmd := exec.Command("go",
		"build",
		"-ldflags",
		"-s -w -X github.com/markdicksonjr/nibbler/build.GitTag="+gitTag,
		"-o",
		"dist/"+platform+"/"+arch+"/"+appName+suffix,
		mainPath)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOOS="+platform)
	if len(arch) > 0 {
		cmd.Env = append(cmd.Env, "GOARCH="+arch)
	}
	return cmd
}

func tieCommandToBuffers(cmd *exec.Cmd) (*bytes.Buffer, *bytes.Buffer) {
	var out bytes.Buffer
	multi := io.MultiWriter(&out)
	cmd.Stdout = multi

	var outErr bytes.Buffer
	multi2 := io.MultiWriter(&outErr)
	cmd.Stderr = multi2
	return &out, &outErr
}
