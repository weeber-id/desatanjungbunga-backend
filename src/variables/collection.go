package variables

// Collection name in mongoDB
var Collection = struct {
	About      string
	Admin      string
	Discussion string
	Article    string
	Travel     string
	Culinary   string
	Handcraft  string

	Lodging           string
	LodgindFacilities string
}{
	About:      "about",
	Admin:      "admin",
	Discussion: "discussion",
	Article:    "article",
	Travel:     "travel",
	Culinary:   "culinary",
	Handcraft:  "handcraft",

	Lodging:           "lodging",
	LodgindFacilities: "lodging_facilities",
}
