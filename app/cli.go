package app

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/prgres/go2git-switch/profile"
	"github.com/spf13/viper"
)

/* --- CLI commands--- */
func profileRemove(profileList []*profile.Profile) error {
	i, _, err := profileListPromt(profileList).Run()
	if err != nil {
		return fmt.Errorf("promp run failed %w", err)
	}

	profileList = append(profileList[:i], profileList[i+1:]...)

	viper.Set("profiles", profileList)
	return nil
}

func profileSelect(target string, profileList []*profile.Profile) error {
	i, _, err := profileListPromt(profileList).Run()
	if err != nil {
		return fmt.Errorf("promp run failed %w", err)
	}

	selected := profileList[i]
	if err := selected.SetProfile(target); err != nil {
		return fmt.Errorf("setting git config failed")
	}

	return nil
}

func profileAdd(profileList []*profile.Profile) error {
	label, err := (&promptui.Prompt{
		Label: "label <can be empty - the name will be used>",
	}).Run()
	if err != nil {
		return err
	}

	name, err := (&promptui.Prompt{
		Label: "name",
		Validate: func(input string) error {
			if input == "" {
				return errors.New("cannot be empty")
			}

			return nil
		},
	}).Run()
	if err != nil {
		return err
	}

	email, err := (&promptui.Prompt{
		Label: "email",
		Validate: func(input string) error {
			if input == "" {
				return errors.New("cannot be empty")
			}

			if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(input) {
				return errors.New("invalid email")
			}

			return nil
		},
	}).Run()
	if err != nil {
		return err
	}

	profileList = append(profileList, &profile.Profile{
		Label: label,
		Name:  name,
		Email: email,
	})

	viper.Set("profiles", profileList)
	return nil
}

// func ProfileCurrent(target string) (*profile.Profile, error) {
// 	return profile.GetProfileActive(target)
// }

/* --- CLI commands helpers func --- */
func getProfileListLabels(profileList []*profile.Profile) []string {
	result := make([]string, len(profileList))
	for i, cofnig := range profileList {
		if cofnig.Label != "" {
			result[i] = cofnig.Label
			continue
		}
		result[i] = cofnig.Name
	}

	return result
}

func profileListPromt(profileList []*profile.Profile) *promptui.Select {
	return &promptui.Select{
		Label:    ">>> Select GIT profile",
		HideHelp: true,
		Items:    getProfileListLabels(profileList),
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .Name }}",
			Selected: fmt.Sprintf(`{{ "%s " }} {{ . | green }}`, promptui.IconGood),
		},
	}
}
