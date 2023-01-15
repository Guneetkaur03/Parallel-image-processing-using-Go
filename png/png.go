// Package png allows for loading png images and applying
// image flitering effects on them
package png

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"image/draw"
)

// The Image represents a structure for working with PNG images.!
type Image struct {
	in     image.Image  
	out    *image.RGBA64  
	Bounds image.Rectangle //The size of the image
}

// subImage interface for dividing into sections
type CropImage interface {
    image.Image
    SubImage(r image.Rectangle) image.Image
}

// image -> ImageTask -> []ImageTask
type ImageTask struct {
	// the exact input path
	InPath   	   string  
	
	// the exact output path
	OutPath		   string

	// what section is the image (only sections)
	SectionCase	   string
	
	// effects to be applied
	Effects        []interface{}       

	// image
	Image          *Image    

	// starting point of section
	YStart    	   int 	        
	
	// ending point of section
	YEnd      	   int 	            	 
}


// Load returns a Image that was loaded based on the filePath parameter
// From Professor Samuels:  You are allowed to modify and update this as you wish
func Load(filePath string) (*Image, error) {

	inReader, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer inReader.Close()

	inOrig, err := png.Decode(inReader)

	if err != nil {
		return nil, err
	}

	bounds := inOrig.Bounds()

	outImg := image.NewRGBA64(bounds)
	inImg := image.NewRGBA64(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := inOrig.At(x, y).RGBA()
			inImg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}
	task := &Image{}
	task.in = inImg
	task.out = outImg
	task.Bounds = bounds
	return task, nil
}

// Save saves the image to the given file
// From Professor Samuels:  You are allowed to modify and update this as you wish
func (img *Image) Save(filePath string, save bool) error {

	outWriter, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outWriter.Close()

	if save {
		err = png.Encode(outWriter, img.in)
		if err != nil {
			return err
		}
	} else {
		err = png.Encode(outWriter, img.out)
		if err != nil {
			return err
		}
	}
	return nil
}


//clamp will clamp the comp parameter to zero if it is less than zero or to 65535 if the comp parameter
// is greater than 65535.
func clamp(comp float64) uint16 {
	return uint16(math.Min(65535, math.Max(0, comp)))
}


// savePic saves the final, filtered picture to the outPath.
func (imageTask *ImageTask) SaveImageTaskOut() {

	err := imageTask.Image.Save(imageTask.OutPath, len(imageTask.Effects) == 0)
	if err != nil {
		panic(err)
	}
}

// get a section
func (img *Image) GetSection(yMin int, yMax int) (*Image, error) {
	// because x is same
	xMax := img.in.Bounds().Max.X

	// create subimage
	var subImg image.Image
	if p, ok := img.in.(CropImage); ok {
		subImg = p.SubImage(image.Rect(0, yMin, xMax, yMax))
	}

	inBounds := subImg.Bounds()
	outImg 	 := image.NewRGBA64(inBounds)

	return &Image{subImg, outImg, image.Rect(0, yMin, xMax, yMax)}, nil
}


// splitImage splits the image into equal sections.
func (imageTask *ImageTask) SplitImage(threads int) []*ImageTask {
	// divided sections
	imageSections := make([]*ImageTask, 0)
	
	// max y 
	yMax := imageTask.Image.in.Bounds().Max.Y

	// size of each section
	sectionSize := int(math.Ceil(float64(yMax)/float64(threads)))

	sectionMin := 0
	sectionMax := 0
	var portion string

	for i := 0; i < threads; i++ {
		// next section
		sectionMin = i * sectionSize
		// last section
		if i == threads-1 {
			sectionMax = yMax
		} else {
			sectionMax = sectionMin + sectionSize
		}
		
		// add padding for the edge cases
		if threads > 1 {
			if i == 0 {
				// upper
				portion = "-1"
				sectionMax += 4
			} else if i > 0 && i < threads-1 {
				// middle
				portion = "0"
				sectionMin -= 4
				sectionMax += 4
			} else {
				// lower
				portion = "1"
				sectionMin -= 4
			}
		} else {
			portion = "2"
		}
		// Grab the section.
		sect, _ :=  imageTask.Image.GetSection(sectionMin, sectionMax)
		imageSections = append(imageSections, &ImageTask{imageTask.InPath, imageTask.OutPath, portion, imageTask.Effects, sect, sectionMin, sectionMax})
	}

	return imageSections
}


// new image
func (img *Image) NewImage() (*Image, error) {
	xMax := img.in.Bounds().Max.X
	yMax := img.in.Bounds().Max.Y
	newImage := image.NewRGBA64(image.Rect(0, 0, xMax, yMax))
	return &Image{img.in, newImage, img.Bounds}, nil
}


// Adds the seection back again to the image
func (img *Image) AddSection(section *Image, yMin int, sectionCase string) {
	// get bounds
	bounds := section.out.Bounds()

	// check to what case it belongs to
	if sectionCase == "-1" {
		// upper case
		bounds.Max.Y -= 4
	} else if sectionCase == "0" {
		// middle case
		bounds.Min.Y += 4
		yMin += 4
		bounds.Max.Y -= 4
	} else if sectionCase == "1" {
		// lower case
		bounds.Min.Y += 4
		yMin += 4
	} 
	// get the cropped version of the section calculated
	area := image.Rect(0, yMin, bounds.Max.X, bounds.Max.Y)
	subImg := section.out.SubImage(area)

	// draw the image on the decided area
	draw.Draw(img.out, bounds, subImg, image.Point{0, yMin}, draw.Src)
}
