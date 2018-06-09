package ports

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

// Alias aliases ports by using a list of aliases as input. An example of this
// would be aliasing `core/openssl` to `6c37/libressl`.
// TODO: Instead of aliasing a single port, alias a whole bunch of ports at
// once.
func (p *Port) Alias(aliases [][]Location) {
	for _, a := range aliases {
		if a[0] == p.Location {
			p.Location = a[1]
		}
	}
}

// All lists all ports found in a specified root directory.
func All(root string) ([]Port, error) {
	var pl []Port
	err := filepath.Walk(root, func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !i.IsDir() && i.Name() == "Pkgfile" {
			pl = append(pl, New(path.Dir(p)))
		}

		return nil
	})
	if err != nil {
		return pl, err
	}

	return pl, nil
}

// Locate tries to locate a port using a given list of Ports. It returns a list
// with possible ports, sorted according to the order parameter.
func Locate(port string, order []string, all []Port) ([]Port, error) {
	// Find matching port names in the `all` list.
	var pl []Port
	for _, p := range all {
		if p.Location.Port == port {
			pl = append(pl, p)
		}
	}

	// If there have been zero matches return with an error.
	if len(pl) == 0 {
		return []Port{}, fmt.Errorf("could not find `%s` in the ports tree",
			port)
	}

	// If there are multiple matches, sort according to the order parameter.
	if len(pl) > 1 {
		var i int
		for _, r := range order {
			if stringInPorts(path.Join(pl[i].Location.Root, r, pl[i].Location.
				Port), all) {
				pl[i].Location.Repo = r
				i++
			}

			// Break once everything has been ordered.
			if i == len(pl) {
				break
			}
		}
	}

	return pl, nil
}
