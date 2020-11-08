# codewords
Handy, sometimes hilarious, sometimes inappropriate, codeword generator.
For when human-friendly identifiers are desired.

This library generates codewords of the form \<adjective\>-\<noun\>. Words
are randomly selected from the Princeton WordNet database. The library
is go-getable and has no dependencies. The codewords are not guaranteed
to be unique.

### Example

```go
cwf := NewFactory()
for i := 0; i < 16; i++ {
    fmt.Println(cwf.Generate())
}
```

```
monied-ixia
extensive-snit
complementary-sachet
writhen-saltworks
underpopulated-soricidae
gritty-marche
protracted-jacobi
azerbaijani-collembolan
mercurous-spalacidae
enthralled-unprofitability
emmetropic-unrelatedness
melodramatic-charlatanism
dopey-locution
comprehended-pothook
checked-hustle
nonparallel-laryngismus
```