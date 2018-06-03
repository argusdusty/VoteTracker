package Utils

func ConvertDate(date string) string {
	return date[:4] + "-" + date[4:6] + "-" + date[6:]
}
