package profile

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

/* --- Profile ---*/
type Profile struct {
	Label string `mapstructure:"label"`
	Name  string `mapstructure:"name"`
	Email string `mapstructure:"email"`
}

func (c Profile) SetProfile(target string) error {
	if err := execCmdGitConfigSetUserName(target, c.Name); err != nil {
		return err
	}

	if err := execCmdGitConfigSetUserEmail(target, c.Email); err != nil {
		return err
	}

	return nil
}

func GetProfileActive(target string) (*Profile, error) {
	nameByte, err := execCmdGitConfigGetUserName(target)
	if err != nil {
		return nil, err
	}

	emailByte, err := execCmdGitConfigGetUserEmail(target)
	if err != nil {
		return nil, err
	}

	return &Profile{
		Name:  string(nameByte),
		Email: string(emailByte),
	}, nil
}

/* --- helpers func --- */

/* -- user name --*/
func execCmdGitConfigGetUserName(target string) ([]byte, error) {
	return execCmdGitConfigGetUser(target, "name")
}

func execCmdGitConfigSetUserName(target string, name string) error {
	return execCmdGitConfigSetUser(target, "name", name)
}

/* -- user email --*/
func execCmdGitConfigGetUserEmail(target string) ([]byte, error) {
	return execCmdGitConfigGetUser(target, "email")
}

func execCmdGitConfigSetUserEmail(target string, email string) error {
	return execCmdGitConfigSetUser(target, "email", email)
}

/* -- ._. --*/
func execCmdGitConfigSetUser(target string, key, value string) error {
	_, err := execCmdGitConfigSet(target, fmt.Sprintf("user.%s", key), value)
	if err != nil {
		return err
	}

	return nil
}

func execCmdGitConfigSet(target string, commad ...string) ([]byte, error) {
	cmdExtended := append(
		[]string{"config", fmt.Sprintf("--%s", target), "--replace-all"},
		commad...,
	)

	return exeCmdGit(cmdExtended...)
}

func execCmdGitConfigGetUser(target string, key string) ([]byte, error) {
	return execCmdGitConfigGet(target, fmt.Sprintf("user.%s", key))
}

func execCmdGitConfigGet(target string, commad ...string) ([]byte, error) {
	cmdExtended := append(
		[]string{"config", fmt.Sprintf("--%s", target)},
		commad...,
	)

	return exeCmdGit(cmdExtended...)
}

func exeCmdGit(commad ...string) ([]byte, error) {
	zap.L().Debug(fmt.Sprintf("Executing git cmd: %s", strings.Join(commad, " ")))
	cmd := exec.Command("git", commad...)
	cmd.Stderr = os.Stdout

	return cmd.Output()
}
