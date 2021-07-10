package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"
	"path/filepath"
	"github.com/nfnt/resize"
	"image/png"
	_ "image/png"
	_ "image/jpeg"
)

func main() {
	// Load rainbow filter
	reader_filter, err := os.Open("rainbow_flag.jpg")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer reader_filter.Close()
	img_filter, _, err := image.Decode(reader_filter)
	if err != nil {
		log.Fatal(err)
		return
	}

	args := os.Args
	path_avatar := args[1]
	var alpha = 0.6
	if len(args) == 3 {
		alpha, err = strconv.ParseFloat(args[2], 32)
	}
	
	fmt.Println(filepath.Abs(path_avatar))
	fmt.Printf("alpha %f\n", alpha)
	reader_avatar, err := os.Open(path_avatar)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer reader_avatar.Close()
	img_avatar, _, err := image.Decode(reader_avatar)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	bounds_avatar := img_avatar.Bounds()
	// fmt.Printf("Image avatar width: %dpix, image height: %dpix\n", bounds_avatar.Dx(), bounds_avatar.Dy())
	filter_resized := resize.Resize(uint(bounds_avatar.Dx()), uint(bounds_avatar.Dy()), img_filter, resize.NearestNeighbor)
	filtered_image := image.NewRGBA64(bounds_avatar)
	for i := 0; i < bounds_avatar.Dx(); i++ {
		for j := 0; j < bounds_avatar.Dy(); j++ {
			color_avatar := img_avatar.At(i, j)
			ra, ga, ba, aa := color_avatar.RGBA()
			color_filter := filter_resized.At(i, j)
			rf, gf, bf, af := color_filter.RGBA()
			r := float64(ra)*alpha+float64(rf)*(1.0-alpha)
			g := float64(ga)*alpha+float64(gf)*(1.0-alpha)
			b := float64(ba)*alpha+float64(bf)*(1.0-alpha)
			a := float64(aa)*alpha+float64(af)*(1.0-alpha)
			v := filtered_image.ColorModel().Convert(color.NRGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)})
			//Alpha = 0: Rainbow flag only
			rr, gg, bb, aa := v.RGBA()
			filtered_image.SetRGBA64(i, j, color.RGBA64{R: uint16(rr), G: uint16(gg), B: uint16(bb), A:uint16(aa)})
		}
	}
	img_new_file_path := path_avatar[:(len(path_avatar)-4)]+"_new.png"
	img_new, err := os.Create(img_new_file_path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer img_new.Close()
	err = png.Encode(img_new, filtered_image)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Done!\n")
	fmt.Printf("New image was saved as %s", img_new_file_path)
}

