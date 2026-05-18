package humanize

// dictPack holds the rewrite tables for one language.
// The pipeline applies AICliches (the main synonym table) first.
// Filler and Upgrades are kept as empty placeholders for future per-language
// fine-tuning (formal vs informal, register adjustments, etc.).
type dictPack struct {
	// AICliches: the main synonym replacement table. Keys and values are
	// case-insensitive at phrase-match time (see replace.go). For multi-language
	// packs, keys must be in the target language's native script.
	AICliches map[string]string

	// Filler: empty intensifiers and hedges that add no information.
	// Currently unused in single-pass synonym mode; reserved for future
	// register adjustments.
	Filler map[string]string

	// Upgrades: bureaucratese-to-plain and slang-to-formal rewrites.
	// Currently unused; reserved for future register adjustments.
	Upgrades map[string]string
}

// Per-language dictionary packs are defined in dict_<lang>.go files:
//
//   synPack    (English)   — dict_en.go    3000+ synonyms
//   synPackRu  (Russian)   — dict_ru.go    1000+ synonyms
//   synPackMo  (Moldovan)  — dict_mo.go    800+ synonyms
//
// pickPack(lang) in language.go selects which pack to use based on the
// detected or user-specified language.
