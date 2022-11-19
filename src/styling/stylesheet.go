package styling

func AssignStyleTag(key string) string {
	tag, ok := StylesheetConfig.StylingConfig[key]
	if ok {
		return tag
	}

	// For all other characters that aren't keywords.
	all, ok := StylesheetConfig.StylingConfig["ALL"]
	if ok {
		return all
	}

	return ""	
}  