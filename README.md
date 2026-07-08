<h1 align="center">
Btxpack
</h1>
<h3 align="center">
Bare Texture Packer
</h3>

## Install
```
go install github.com/intervinn/btxpack/cmd/btxpack@latest
```

## Usage
Consumes folder of images, emits combined image and the metadata file
```bash
btxpack -src=/assets -out=atlas.png -meta=atlas.json
```

## Algorithms
* `line` - places each assets one by one in a single line
* `shelf` - fits all assets into a square with a size of power of two for better gpu something, uses a shelf algorithm for vertical stacking

## Changelog
### 0.1.4
* C codegen output is now an stb header
### 0.1.3
* Adjust C codegen
* Add raylib C codegen (`-raylib`)
### 0.1.2
* Added Shelf algorithm (`-alg=shelf`)
### 0.1.1
* Refactor cmd arg parsing (remove cobra, use standard go flags)
* Add C codegen
### 0.1.0
* Initial release

## License
MIT License.