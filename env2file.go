package main

import (
	"fmt"
	"os"
	"strings"
)

const prefix = "EF_"

// target represents a filename (full path) and its contents.
type target struct {
	name string
	data string
}

// writes the contents of a target to its desired file.
func (t target) write() error {
	f, err := os.Create(t.name)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(t.data); err != nil {
		return err
	}
	return nil
}

// aggregateFromEnv collects a set of targets from the given env vars
// it only accepts pairs, so if a `name` is missing a `data` then it's ignored
func aggregateFromEnv(envs []string) ([]target, []error) {
	targ := make(map[string]struct{})
	name := make(map[string]string)
	data := make(map[string]string)
	errs := []error{}
	for _, e := range envs {
		key, value := splitEnvironmentVariable(e)
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		t, n, err := decodeKey(key)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		targ[n] = struct{}{}
		switch t {
		case "name":
			name[n] = value
		case "data":
			data[n] = value
		}
	}
	targets, joinErrors := join(targ, name, data)
	return targets, append(errs, joinErrors...)
}

// joins the aggregated sets of separate variables into a list of targets
func join(t map[string]struct{}, n, d map[string]string) ([]target, []error) {
	all := []target{}
	errs := []error{}
	for key := range t {
		var t target
		var ok bool
		if t.name, ok = n[key]; !ok {
			errs = append(errs, fmt.Errorf("missing target config for %s", key))
			continue
		}
		if t.data, ok = d[key]; !ok {
			errs = append(errs, fmt.Errorf("missing target config for %s", key))
			continue
		}
		all = append(all, t)
	}
	return all, errs
}

func splitEnvironmentVariable(keyvalue string) (key, value string) {
	v := strings.SplitN(keyvalue, "=", 2)
	return v[0], v[1]
}

// splits a key in the E2F format into a target type and a target name
// errors if the pattern is wrong such as having too few _ separators
func decodeKey(key string) (targetType string, name string, err error) {
	v := strings.SplitN(key, "_", 3)
	if len(v) != 3 {
		return "", "", fmt.Errorf("%s has invalid pattern", key)
	}
	return v[1], v[2], nil
}

func main() {
	targets, errors := aggregateFromEnv(os.Environ())
	for _, e := range errors {
		fmt.Println("Error:", e)
	}
	for _, t := range targets {
		if err := t.write(); err != nil {
			fmt.Println(err)
		}
	}
}
