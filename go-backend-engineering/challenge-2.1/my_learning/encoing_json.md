Source: https://pkg.go.dev/encoding/json

## 1. A way to implement enums:
type Animal int

const (
	Unknown Animal = iota (iota means, first 0, then other 1, 2 ,,,)
	Gopher
	Zebra
)

if not specify iota, but "1" for example, all untyped will have the same value as the position above itself.

- Note: If not defining Animal or something similar, the compiler won't catch misuse --> Best Practice: define a type in advance.

## 2. JSON in MarshalJSON

When add JSON after the method name, example, MarshalJSON, it will trigger a built-in feature and allows us to run this code block:
```
blob := `["gopher","armadillo","zebra","unknown","gopher","bee","gopher","zebra"]`
var zoo []Animal
json.Unmarshal([]bytes(blob), &zoo)
```



func (a *Animal) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		*a = Unknown
	case "gopher":
		*a = Gopher
	case "zebra":
		*a = Zebra
	}

	return nil
}

blob := `["gopher","armadillo","zebra","unknown","gopher","bee","gopher","zebra"]`
var zoo []Animal
if err := json.Unmarshal([]byte(blob), &zoo); err != nil {
	log.Fatal(err)
}