package variables

import (
	log "github.com/sirupsen/logrus"
	"regexp"
)

const (
	variablesRegPattern = "\\$\\{\\w+\\}"

	SystemLastSuccessVars = "${system.last_success_vars}"
)

var (
	variableReq *regexp.Regexp
)

var (
	variablesMap = map[string]VariableHandle{
		SystemLastSuccessVars: &LastSuccessVars{},
	}
)

func init() {
	var err error
	variableReq, err = regexp.Compile(variablesRegPattern)
	if err != nil {
		log.Fatalf("regexp.Compile(%s) failed. err=%s", variablesRegPattern, err)
	}
}

func ReplaceVariables(str string, handle func(string, VariableHandle) error) error {
	rs := variableReq.FindAllString(str, -1)
	if len(rs) == 0 {
		return handle(str, nil)
	}

	for i := 0; i < len(rs); i++ {
		term := rs[i]

		if h, ok := variablesMap[term]; ok {
			err := handle(term, h)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
