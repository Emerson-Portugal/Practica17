package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func check4(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Se puede obtener el resultado de la para un x = 0,25, x = 0,5 y x = 0,75.
	// el porcentaje se le aplicara a la imagen 1 para ver como se gradua
	porcentaje := 0.75
	// Cargamos la prueba 1
	imgPath := "prueba1.jpeg"
	f, err := os.Open(imgPath)
	check4(err)
	defer f.Close()

	//pasamos la iamgen a img
	img, _, err := image.Decode(f)

	//cargamos la prueba 2
	imgPath2 := "prueba2.jpeg"
	f2, err := os.Open(imgPath2)
	check4(err)
	defer f2.Close()

	//pasamos la iamgen a img
	img2, _, err := image.Decode(f2)

	//dimenciones del resultado
	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	wImg := image.NewRGBA(rect)

	// se inicializa la varible para la concurencia
	wg := new(sync.WaitGroup)

	start := time.Now()

	for x := 0; x < size.X; x++ {
		// and now loop thorough all of this x's y
		wg.Add(1)
		x := x
		// se le agrega le paramero go -> para volverlo concurrente
		go func() {
			for y := 0; y < size.Y; y++ {
				pixel := img.At(x, y)
				originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
				pixel2 := img2.At(x, y)
				originalColor2 := color.RGBAModel.Convert(pixel2).(color.RGBA)

				// Q(i, j) = X ∗ P1(i, j) + (1 − X) ∗ P2(i, j) (1)
				// Donde: Q(i, j), es un pixel de la imagen resultado; P1(i, j), es un pixel de la imagen de entrada;
				// P2(i, j), es un pixel de la otra imagen de entrada.
				//X = porcentaje
				r := uint8(math.Sqrt((math.Pow(porcentaje*float64(originalColor.R), 2) + math.Pow((1-porcentaje)*float64(originalColor2.R), 2)) / 2))
				g := uint8(math.Sqrt((math.Pow(porcentaje*float64(originalColor.G), 2) + math.Pow((1-porcentaje)*float64(originalColor2.G), 2)) / 2))
				b := uint8(math.Sqrt((math.Pow(porcentaje*float64(originalColor.B), 2) + math.Pow((1-porcentaje)*float64(originalColor2.B), 2)) / 2))
				a := uint8(math.Sqrt((math.Pow(porcentaje*float64(originalColor.A), 2) + math.Pow((1-porcentaje)*float64(originalColor2.A), 2)) / 2))

				c := color.RGBA{
					R: r, G: g, B: b, A: a,
				}
				wImg.Set(x, y, c)
			}
			defer wg.Done()
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("Tiempo %s", elapsed)

	ext := filepath.Ext(imgPath)
	name := strings.TrimSuffix(filepath.Base(imgPath), ext)
	newImagePath := fmt.Sprintf("%s/%s_blending3%s", filepath.Dir(imgPath), name, ext)
	fg, err := os.Create(newImagePath)
	defer fg.Close()
	check4(err)
	err = jpeg.Encode(fg, wImg, nil)
	check4(err)
}
