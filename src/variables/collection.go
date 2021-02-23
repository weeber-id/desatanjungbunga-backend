package variables

// Collection name in mongoDB
var Collection = struct {
	Admin      string
	Discussion string
	Article    string
	Travel     string
	Kuliner    string
	Belanja    string
}{
	Admin:      "admin",
	Discussion: "discussion",
	Article:    "article",
	Travel:     "travel",
	Kuliner:    "kuliner",
	Belanja:    "belanja",
}
