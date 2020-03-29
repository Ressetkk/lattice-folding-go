# Protein folding in Go
Requirements:
* Go version 1.13+

To run an application, simply compile it and execute the binary.
The script accepts flags:
* `-protein` character stream over the alphabet of {h, p}.
* `-output` name of the generated image [Default 'out.png']

```shell script
cd lattice-folding-go
go build
./lattice-folding-go -protein [hp/HP*]
```

The application will return the image of protein with lowest energy. For really long chains the application can do it's calculations really long.

THe algorithm used is covered in https://www.ncbi.nlm.nih.gov/pmc/articles/PMC5172541/