package variables

// Collection name in mongoDB
var Collection = struct {
	Article string
	Wisata  string
	Kuliner string
	Belanja string
}{
	Article: "article",
	Wisata:  "wisata",
	Kuliner: "kuliner",
	Belanja: "belanja",
}
