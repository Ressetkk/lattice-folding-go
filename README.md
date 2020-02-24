# Proteing folding in Go
Requirements:
* Go version 1.13+
* A lot of RAM

To run an application, simply compile it and execute the binary.
The script accepts one flag `-protein` that accepts words over language {h,p}.
```shell script
cd lattice-folding-go
go build
./lattice-folding-go -protein [hp/HP*]
```

The application will work for 30 seconds and after the time it will return the lowest energy.

THe algorithm used is covered in https://www.ncbi.nlm.nih.gov/pmc/articles/PMC5172541/

**The application is far from perfect and might crash due to out of memory errors. It does not return the generated protein image (yet).