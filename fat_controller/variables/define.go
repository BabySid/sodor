package variables

import (
	"fmt"
	"regexp"
)

const (
	LastSystemSuccessInstance = "system.last_success_instance"
)

func findVariables(content string) {
	req, err := regexp.Compile("\\$\\{\\w+\\}")
	if err != nil {
		panic(err)
	}

	ls := req.FindStringSubmatch("echo ${dt}")

	if len(ls) > 0 {
		n := len(ls)
		for i := 0; i < n; i++ {
			fmt.Println(ls[i])
		}
	}
}
