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

func check3(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Cargamos la prueba 1
	imgPath := "prueba1.jpeg"
	f, err := os.Open(imgPath)
	check3(err)
	defer f.Close()

	//pasamos la iamgen a img
	img, _, err := image.Decode(f)

	//cargamos la prueba 2
	imgPath2 := "prueba2.jpeg"
	f2, err := os.Open(imgPath2)
	check3(err)
	defer f2.Close()

	//pasamos la imagen a img
	img2, _, err := image.Decode(f2)

	//dimenciones del resultado
	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	wImg := image.NewRGBA(rect)

	// se inicializa la varible para la concurencia
	wg := new(sync.WaitGroup)

	start := time.Now()

	for x := 0; x < size.X; x++ {

		wg.Add(1)
		x := x
		// se le agrega le paramero go -> para volverlo concurrente
		go func() {
			for y := 0; y < size.Y; y++ {
				pixel := img.At(x, y)
				originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
				pixel2 := img2.At(x, y)
				originalColor2 := color.RGBAModel.Convert(pixel2).(color.RGBA)

				// Solución por adelantado: cómo promediar los colores RGB
				// Primero debe elevar al cuadrado los colores y luego sumarlos y luego sacar la raíz cuadrada.
				//Aquí hay un ejemplo con 0 y 255 que es el más alto en cada espectro.

				// 0 * 0 = 0
				// 255*255= 65025
				// 0 + 65025 = 65025
				// 65025/2 = 32512,5
				// raíz cuadrada de 32512,5 = 180,31
				// El promedio de 0 y 255 NO es 127,5, ¡sino 180!
				// Y esa es una gran diferencia, porque si lo dividieras por 2, ¡tu promedio sería mucho más oscuro!

				r := uint8(math.Sqrt((math.Pow(float64(originalColor.R), 2) + math.Pow(float64(originalColor2.R), 2)) / 2))
				g := uint8(math.Sqrt((math.Pow(float64(originalColor.G), 2) + math.Pow(float64(originalColor2.G), 2)) / 2))
				b := uint8(math.Sqrt((math.Pow(float64(originalColor.B), 2) + math.Pow(float64(originalColor2.B), 2)) / 2))
				a := uint8(math.Sqrt((math.Pow(float64(originalColor.A), 2) + math.Pow(float64(originalColor2.A), 2)) / 2))

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
	newImagePath := fmt.Sprintf("%s/%s_overflow%s", filepath.Dir(imgPath), name, ext)
	fg, err := os.Create(newImagePath)
	defer fg.Close()
	check3(err)
	err = jpeg.Encode(fg, wImg, nil)
	check3(err)
}
