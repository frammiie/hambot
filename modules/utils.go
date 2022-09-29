package modules

func ConcatArgs(args ...string) string {
	query := args[0]
	for i := 1; i < len(args); i++ {
		query += " " + args[i]
	}

	return query
}
