package variables

// Collection name in mongoDB
var Collection = struct {
	Admin      string
	Discussion string
	Article    string
	Travel     string
	Culinary   string
	Handcraft  string
}{
	Admin:      "admin",
	Discussion: "discussion",
	Article:    "article",
	Travel:     "travel",
	Culinary:   "culinary",
	Handcraft:  "handcraft",
}
