package app

import (
	"fmt"

	"github.com/prgres/go2git-switch/profile"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

func ProfileSelect(c *cli.Context) error {
	profileList, err := getProfiles()
	if err != nil {
		return err
	}

	if err := profileSelect(c.String("target"), profileList); err != nil {
		return err
	}

	return ProfileCurrent(c)
}

func ProfileAdd(c *cli.Context) error {
	profileList, err := getProfiles()
	if err != nil {
		return err
	}

	if err := profileAdd(profileList); err != nil {
		return err
	}

	return ProfileSelect(c)
}

func ProfileRemove(c *cli.Context) error {
	profileList, err := getProfiles()
	if err != nil {
		return err
	}

	if err := profileRemove(profileList); err != nil {
		return err
	}

	return ProfileRemove(c)
}

func ProfileCurrent(c *cli.Context) error {
	target := c.String("target")
	activeProfile, err := profile.GetProfileActive(target)
	if err != nil {
		return err
	}

	fmt.Printf(">>> Active git profile in \"--%s\":\n\tNAME: %s\tEMAIL: %s", target, activeProfile.Name, activeProfile.Email)
	return nil
}

func getProfiles() ([]*profile.Profile, error) {
	var profileList []*profile.Profile
	if err := viper.UnmarshalKey("profiles", &profileList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return profileList, nil
}
