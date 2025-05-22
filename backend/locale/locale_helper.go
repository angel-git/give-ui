package locale

func ConvertLocale(locale string) string {
	switch locale {
	case "English":
		return "en"
	case "Chinese":
		return "ch"
	case "Czech":
		return "cz"
	case "French":
		return "fr"
	case "German":
		return "ge"
	case "Hungarian":
		return "hu"
	case "Italian":
		return "it"
	case "Japanese":
		return "jp"
	case "Korean":
		return "kr"
	case "Polish":
		return "pl"
	case "Portuguese":
		return "po"
	case "Slovak":
		return "sk"
	case "Spanish":
		return "es"
	case "Spanish - Mexico":
		return "es-mx"
	case "Turkish":
		return "tu"
	case "Romanian":
		return "ro"
	case "Русский":
		return "ru"
	default:
		return "en"
	}
}
