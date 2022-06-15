package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"stash.sigma.sbrf.ru/sddevops/terraform-provider-di/translit"
)

func Reformat(name string) string {
	name = strings.TrimPrefix(name, " ")
	name = strings.ToLower(name)
	name = translit.EncodeToICAO(name)
	if string(name[0]) != `_` && !unicode.IsLetter(rune(name[0])) {
		name = fmt.Sprintf("_%s", name)
	}
	replace := []string{"-", " - ", ". ", " ", "."}
	remove := []string{":", ",", "(", ")", "[", "]"}
	for _, v := range replace {
		name = strings.ReplaceAll(name, v, "_")
	}
	for _, v := range remove {
		name = strings.ReplaceAll(name, v, "")
	}
	if len(name) > 50 {
		name = name[:50]
	}
	return name
}

func Regexp(data []byte) []byte {
	comments := []*regexp.Regexp{
		regexp.MustCompile(` id(\s)*= "(.*)"`),
		regexp.MustCompile(`_uuid(\s)*= "(.*)"`),
	}
	replaces := []*regexp.Regexp{
		regexp.MustCompile(` domain_id(\s)*= "(.*)"`),
		regexp.MustCompile(` group_id(\s)*= "(.*)"`),
		regexp.MustCompile(` app_systems_ci(\s)*= "(.*)"`),
		regexp.MustCompile(` stand_type_id(\s)*= "(.*)"`),
		regexp.MustCompile(` project_id(\s)*= "(.*)"`),
		regexp.MustCompile(` value(\s)*= "data.(.*)"`),
		regexp.MustCompile(` value(\s)*= "di_(.*)"`),
		regexp.MustCompile(` postgres_db_password(\s)*= "data.vault(.*)"`),
		regexp.MustCompile(` gg_client_password(\s)*= "data.vault(.*)"`),
		regexp.MustCompile(` ise_client_password(\s)*= "data.vault(.*)"`),
	}
	appParamsBlockToMap := regexp.MustCompile(`app_params \{`)

	parts := strings.Split(string(data), "\n")
	for k, v := range parts {
		for _, val := range comments {
			if val.MatchString(v) {
				parts[k] = fmt.Sprintf("#%s", v)
			}
		}
		for _, val := range replaces {
			if val.MatchString(v) {
				parts[k] = strings.ReplaceAll(v, `"`, "")
			}
		}
		if appParamsBlockToMap.MatchString(v) {
			parts[k] = strings.ReplaceAll(v, `{`, `= {`)
		}
	}
	out := strings.Join(parts, "\n")
	return []byte(out)
}
