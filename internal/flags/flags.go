package flags

import (
	"fmt"
	"strconv"
	"strings"
)

type Flag struct {
	Long string
	Short string
	Value any
}

type Flags map[string]*Flag;

func (flags Flags) Var(long, short string, variable any) {
	flag := Flag{long, short, variable};
	flags[long] = &flag;
	if short != "" {
		flags[short] = &flag;
	}
}

func (flags Flags) Parse(argv []string) error {
	for i := 0; i < len(argv); i++ {
		arg := argv[i];

		var err error;

		switch {
		case strings.HasPrefix(arg, "--"):
			s := arg[2:];
			flag := flags[s];
			if flag == nil {
				return fmt.Errorf("unknown flag: %v", arg);
			}

			err = setFlag(flag, s, &i, argv);
			if err != nil {
				return err;
			}
		case strings.HasPrefix(arg, "-"):
			shorts := arg[1:];
			loop: for j := range shorts {
				s := string(shorts[j]);
				flag := flags[s];
				if flag == nil {
					return fmt.Errorf("unknown flag: %v", arg);
				}

				err = setFlag(flag, s, &i, argv);
				if err != nil {
					return err;
				}
				if _, ok := flag.Value.(*bool); !ok {
					break loop;
				}
			}
		}

	}

	return nil;
}

func setFlag(flag *Flag, s string, i* int, argv []string) error {
	switch v := flag.Value.(type) {
	case *string:
		if *i+1 >= len(argv) {
			return fmt.Errorf("flag %v needs an argument", s);
		}

		*i++
		*v = argv[*i];
	case *int:
		if *i+1 >= len(argv) {
			return fmt.Errorf("flag %v needs an argument", s);
		}
		*i++
		conv, err := strconv.Atoi(argv[*i]);
		if err != nil {
			return fmt.Errorf("argument for %v flag needs to be an integer", s);
		}
		*v = conv;
	case *bool:
		*v = true;
	}
	return nil;
}

