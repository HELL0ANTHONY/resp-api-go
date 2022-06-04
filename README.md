# REST API without any external library

---

```go
/*
  Esta función devuelve el "pointer" de un "struct" con unos datos
  ya hardcodeados
*/
func NewCoasterHandlers() *coastersHanlders {
	return &coastersHanlders{
		store: map[string]Coaster{
			"id1": {
				Id:           "id1",
				Name:         "Fury 23",
  			Height:       89,
				InPark:       "Carowinds",
				Manufacturer: "B+M",
			},
		},
	}
}
```
