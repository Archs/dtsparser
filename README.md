# dtsparser
Parsing typescript .d.ts file to generate JSON or Gopherjs package automatically! :)

The parser is created using [PEG for Go](https://github.com/pointlander/peg).

This work is in progress(json output is now ok), I'll finish it soon. 

# Usage

    Usage: ./dts "path/to/target/file"

Json would be output to stdout.

# ATTENTION

1. This is not a complete `.d.ts` parser, just for easy coverting.
2. Use UTF8 when parsing
3. BOM is not supported
    
    if parsing not pass, try saving the file without UTF8 BOM first

4. Anonymouse function in Object/Interface are not processed
