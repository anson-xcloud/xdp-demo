package xflag

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

type FlagSet struct {
	*pflag.FlagSet
}

func NewFlagSet() *FlagSet {
	return &FlagSet{FlagSet: &pflag.FlagSet{}}
}

func Wrap(fs *pflag.FlagSet) *FlagSet {
	return &FlagSet{FlagSet: fs}
}

func (f *FlagSet) AddNamedFlagSet(name string, fs *FlagSet) {
	if fs == nil {
		return
	}

	fs.VisitAll(func(flag *pflag.Flag) { f.AddNamedFlag(name, flag) })
}

func (f *FlagSet) AddNamedFlag(name string, flag *pflag.Flag) {
	flag.Name = fmt.Sprintf("%s.%s", name, flag.Name)
	flag.Name = strings.ToLower(flag.Name)
	if f.Lookup(flag.Name) == nil {
		f.AddFlag(flag)
	}
}

func (f *FlagSet) AllNames() []string {
	ret := []string{}
	f.VisitAll(func(pf *pflag.Flag) {
		ret = append(ret, pf.Name)
	})
	return ret
}
