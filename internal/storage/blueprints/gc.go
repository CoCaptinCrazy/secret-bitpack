package blueprints

func Unreferenced(paths, referenced map[string]bool) []string {
	var result []string
	for path := range paths {
		if !referenced[path] {
			result = append(result, path)
		}
	}
	return result
}
