package variables

// Collection name in mongoDB
var Collection = struct {
	Discussion string
	Article    string
	Wisata     string
	Kuliner    string
	Belanja    string
}{
	Discussion: "discussion",
	Article:    "article",
	Wisata:     "wisata",
	Kuliner:    "kuliner",
	Belanja:    "belanja",
}
