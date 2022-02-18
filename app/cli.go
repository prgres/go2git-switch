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
	i, _, err := profileListPrompt(profileList).Run()
	if err != nil {
		return fmt.Errorf("prompt run failed %w", err)
	}

	profileList = append(profileList[:i], profileList[i+1:]...)

	viper.Set("profiles", profileList)
	return nil
}

func profileSelect(target string, profileList []*profile.Profile) error {
	i, _, err := profileListPrompt(profileList).Run()
	if err != nil {
		return fmt.Errorf("prompt run failed %w", err)
	}

	selected := profileList[i]
	if err := selected.SetProfile(target); err != nil {
		return fmt.Errorf("setting git config failed")
	}

	return nil
}

func profileEdit(profileList []*profile.Profile) error {
	i, _, err := profileListPrompt(profileList).Run()
	if err != nil {
		return fmt.Errorf("prompt run failed %w", err)
	}

	profileToEdit := profileList[i]

	label, err := (&promptui.Prompt{
		Label:       "label <can be empty - the name will be used>",
		Default:     profileToEdit.Label,
		HideEntered: true,
		AllowEdit:   true,
	}).Run()
	if err != nil {
		return err
	}

	name, err := (&promptui.Prompt{
		Label:       "name",
		Default:     profileToEdit.Name,
		HideEntered: true,
		AllowEdit:   true,
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
		Label:       "email",
		Default:     profileToEdit.Email,
		HideEntered: true,
		AllowEdit:   true,
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

	profileList[i] = &profile.Profile{
		Label: label,
		Name:  name,
		Email: email,
	}

	viper.Set("profiles", profileList)

	return nil
}

func profileAdd(profileList []*profile.Profile) error {
	label, err := (&promptui.Prompt{
		HideEntered: true,
		Label:       "label <can be empty - the name will be used>",
	}).Run()
	if err != nil {
		return err
	}

	name, err := (&promptui.Prompt{
		Label:       "name",
		HideEntered: true,
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
		Label:       "email",
		HideEntered: true,
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
	for i, config := range profileList {
		if config.Label != "" {
			result[i] = config.Label
			continue
		}
		result[i] = config.Name
	}

	return result
}

func profileListPrompt(profileList []*profile.Profile) *promptui.Select {
	return &promptui.Select{
		Label:        ">>> Select GIT profile",
		HideHelp:     true,
		HideSelected: true,
		Items:        getProfileListLabels(profileList),
		Templates: &promptui.SelectTemplates{
			Label: "{{ .Name }}",
		},
	}
}
