package utils

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"math/rand"
)

const (
	CHUNK_SIZE = 32
)

var (
	COLORS = [][]uint8{
		{ 0xD3, 0x2F, 0x2F },  // Red 700          #D32F2F
		{ 0xC2, 0x18, 0x5B },  // Pink 700         #C2185B
		{ 0x7B, 0x1F, 0xA2 },  // Purple 700       #7B1FA2
		{ 0x51, 0x2D, 0xA8 },  // Deep Purple 700  #512DA8
		{ 0x19, 0x76, 0xD2 },  // Blue 700         #1976D2
		{ 0x00, 0x97, 0xA7 },  // Cyan 700         #0097A7
		{ 0x00, 0x79, 0x6B },  // Teal 700         #00796B
		{ 0x38, 0x8E, 0x3C },  // Green 700        #388E3C
		{ 0xAF, 0xB4, 0x2B },  // Lime 700         #AFB42B
		{ 0xFB, 0xC0, 0x2D },  // Yellow 700       #FBC02D
		{ 0xFF, 0xA0, 0x00 },  // Amber 700        #FFA000
		{ 0xF5, 0x7C, 0x00 },  // Orange 700       #F57C00
	}
)

/*
 * Generates a random base64 image
 */
func RandomImage(size int) ([]byte, error) {
	// Must be a base 2 sized image //
	if (size > CHUNK_SIZE && math.Floor(math.Log2(float64(size))) != math.Log2(float64(size))) {
		return nil, errors.New(fmt.Sprintf("Image size must be base 2 and > %d", CHUNK_SIZE))
	}

	img := image.NewRGBA(image.Rect(0, 0, size, size))

	numChunks := size / CHUNK_SIZE
	pickedColors := []int{rand.Intn(len(COLORS)), rand.Intn(len(COLORS)), rand.Intn(len(COLORS))}

	for chunkY := 0; chunkY < numChunks; chunkY++ {
		for chunkX := 0; chunkX < numChunks; chunkX++ {
			index := rand.Intn(len(pickedColors))
			color := color.RGBA{
				COLORS[pickedColors[index]][0],
				COLORS[pickedColors[index]][1],
				COLORS[pickedColors[index]][2],
				0xFF,
			}

			// For each pixel //
			for x := 0; x < CHUNK_SIZE; x++ {
				for y := 0; y < CHUNK_SIZE; y++ {
					img.Set(
						x + (chunkX * CHUNK_SIZE), 
						y + (chunkY * CHUNK_SIZE), 
						color,
					)
				}
			}
		}
	}

	var buff bytes.Buffer
	jpeg.Encode(&buff, img, nil)

	return buff.Bytes(), nil 
}
