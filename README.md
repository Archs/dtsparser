# dtsparser
Parsing typescript .d.ts file to generate JSON or Gopherjs package automatically! :)

The parser is created using [PEG for Go](https://github.com/pointlander/peg).

This work is in progress, I'll finish it soon. 

# ATTENTION

1. Use UTF8 when parsing

2. BOM is not supported
    
    if parsing not pass, try saving the file without UTF8 BOM first
