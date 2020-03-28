# Protein folding in Go
Requirements:
* Go version 1.13+

To run an application, simply compile it and execute the binary.
The script accepts flags:
* `-protein` character stream over the alphabet of {h, p}.
* `p1` first probability [Default 0.4]
* `p2` "second probability [Default 0.4]
* `output` name of the generated image [Default 'out.png']

```shell script
cd lattice-folding-go
go build
./lattice-folding-go -protein [hp/HP*]
```

The application will work for 30 seconds and after the time it will return the lowest folded protein chain and save the chain as an image.

THe algorithm used is covered in https://www.ncbi.nlm.nih.gov/pmc/articles/PMC5172541/