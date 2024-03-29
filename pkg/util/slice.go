package util

/*
	Helper function
	reference: https://stackoverflow.com/questions/34070369/removing-a-string-from-a-slice-in-go
*/
func Remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
