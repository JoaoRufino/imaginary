package main

// Version stores the current package semantic version
var Version = "dev"

// Versions represents the used versions for several significant dependencies
// swagger:model imaginary_Versions
type Versions struct {
	//Imaginary Server Version
	//
	//example: 0.1.28
	ImaginaryVersion string `json:"imaginary"`
	//Bimg library Version
	//
	//example: 0.1.5
	BimgVersion string `json:"bimg"`
	//Vips Library Version
	//
	//example: 8.4.1
	VipsVersion string `json:"libvips"`
}
