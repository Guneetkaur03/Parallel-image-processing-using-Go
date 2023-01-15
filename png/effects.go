// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	//"fmt"
	"image"
	"image/color"
)

// Grayscale applies a grayscale filtering effect to the image
func (img *Image) Grayscale() {

	// Bounds returns defines the dimensions of the image. Always
	// use the bounds Min and Max fields to get out the width
	// and height for the image
	bounds := img.out.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Returns the pixel (i.e., RGBA) value at a (x,y) position
			// Note: These get returned as int32 so based on the math you'll
			// be performing you'll need to do a conversion to float64(..)
			r, g, b, a := img.in.At(x, y).RGBA()

			//Note: The values for r,g,b,a for this assignment will range between [0, 65535].
			//For certain computations (i.e., convolution) the values might fall outside this
			// range so you need to clamp them between those values.
			greyC := clamp(float64(r+g+b) / 3)

			//Note: The values need to be stored back as uint16 (I know weird..but there's valid reasons
			// for this that I won't get into right now).
			img.out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}

	// swap the pointer
	img.in = img.out
}

/*
	perferms convolution given kernel and input
*/
func (img *Image) Convolve(kernel [3][3]float64) {

	bounds := img.out.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b uint32
			var sumR, sumG, sumB float64

			_, _, _, a := img.in.At(x, y).RGBA()

			// sumR := uint16(0)
			// sumB := uint16(0)
			// sumG := uint16(0)

			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {

					//get the image (x + i, y + j) position
					xpos := x + i - 1
					ypos := y + j - 1

					//handle  boundary condition
					if xpos < 0 || ypos < 0 {
						r = 0
						g = 0
						b = 0
					} else {
						// get the rgba values
						r, g, b, _ = img.in.At(xpos, ypos).RGBA()
					}

					// calculate convolve
					sumR += float64(r) * kernel[i][j]
					sumG += float64(g) * kernel[i][j]
					sumB += float64(b) * kernel[i][j]

				}
			}
			img.out.Set(x, y, color.RGBA64{clamp(sumR), clamp(sumG), clamp(sumB), uint16(a)})
		}
	}

}

/*
	Sharpen Effect
*/
func (img *Image) Sharpen() {
	kernel := [3][3]float64{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}

	// apply affect
	img.Convolve(kernel)

}

/*
	Edge Detection Effect
*/
func (img *Image) EdgeDetection() {
	kernel := [3][3]float64{
		{-1, -1, -1},
		{-1, 8, -1},
		{-1, -1, -1},
	}

	// apply affect
	img.Convolve(kernel)
}

/*
	Blur Effect
*/
func (img *Image) Blur() {
	kernel := [3][3]float64{
		{1 / 9.0, 1 / 9.0, 1 / 9.0},
		{1 / 9.0, 1 / 9.0, 1 / 9.0},
		{1 / 9.0, 1 / 9.0, 1 / 9.0},
	}

	// apply affect
	img.Convolve(kernel)
}


// swap the image pointers when rquired
func (img *Image) SwapImage() {
	// swap
	img.in 	= img.out
	
	// create blank canvas
	img.out = image.NewRGBA64(img.in.Bounds())
}
