package variables

// Collection name in mongoDB
var Collection = struct {
	Admin      string
	Discussion string
	Article    string
	Wisata     string
	Kuliner    string
	Belanja    string
}{
	Admin:      "admin",
	Discussion: "discussion",
	Article:    "article",
	Wisata:     "wisata",
	Kuliner:    "kuliner",
	Belanja:    "belanja",
}
