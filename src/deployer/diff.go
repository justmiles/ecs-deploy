package deployer

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Diff represents data to be added and removed
type Diff struct {
	resource string
	name     string
	changes  []string
}

func (diff Diff) String() string {
	fmt.Println(fmt.Sprintf("\n%s \"%s\" { ", diff.resource, diff.name))
	return strings.Join(diff.changes, "\n") + "\n}\n"
}

func NewDiff(resource string, name string) Diff {
	return Diff{
		resource: resource,
		name:     name,
		changes:  []string{},
	}
}

func (diff *Diff) AddChange(key, x, y string) {
	if x == y {
		diff.changes = append(diff.changes, color.WhiteString(fmt.Sprintf(" \t%s =\t%s", key, y)))
	} else if x == "" {
		diff.changes = append(diff.changes, color.GreenString(fmt.Sprintf("+\t%s =\t%s", key, y)))
	} else if y == "" {
		// Removing from x
		diff.changes = append(diff.changes, color.RedString(fmt.Sprintf("-\t%s =\t%s", key, x)))
	} else {
		diff.changes = append(diff.changes, color.YellowString(fmt.Sprintf("~\t%s =\t%s  -->  %s", key, x, y)))
	}
}
