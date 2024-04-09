package internal

import (
	"fmt"
	"github.com/henges/envy/internal/parse"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"os/exec"
)

func RunWithEnvFile(envPath, procPath string, procArgs []string) error {

	env, err := readEnvFile(envPath)
	if err != nil {
		log.Err(err).Send()
		return err
	}
	envSl := buildEnviron(env)
	return runProc(procPath, procArgs, envSl)
}

func readEnvFile(envPath string) (map[string]string, error) {
	open, err := os.Open(envPath)
	if err != nil {
		return nil, err
	}
	defer open.Close()
	bs, err := io.ReadAll(open)
	if err != nil {
		return nil, err
	}

	return parse.EnvFile(string(bs))
}

func buildEnviron(input map[string]string) []string {

	original := os.Environ()
	for k, v := range input {
		original = append(original, fmt.Sprintf("%s=%s", k, v))
	}

	return original
}

func runProc(procPath string, args []string, env []string) error {

	command := exec.Command(procPath, args...)
	command.Env = env
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
