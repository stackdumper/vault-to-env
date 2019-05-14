package command

import "strings"

// FlagList represents a list of flag arguments
type FlagList []string

// Set is used to add a flag to a list
func (envs *FlagList) Set(value string) error {
	*envs = append(*envs, value)

	return nil
}

// String is used to convert FlagList to string
func (envs *FlagList) String() string {
	return strings.Join(*envs, ", ")
}
