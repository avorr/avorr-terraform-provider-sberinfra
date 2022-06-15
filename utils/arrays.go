package utils

func ArrInterfaceToArrStr(i interface{}) []string {
	r := []string{}
	for _, v := range i.([]interface{}) {
		r = append(r, v.(string))
	}
	return r
}

func ArrContainsStr(s []string, item string) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}
