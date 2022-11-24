package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// funcion para corroborar si hay error
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Creamos los arrar para guaradar los valores recolectrados
var ArrayR []int
var ArrayG []int
var ArrayB []int

func main() {
	//carga el archivo de la primera imagen
	imgPath := "atardecer.jpeg"
	f, err := os.Open(imgPath)
	check(err)
	defer f.Close()

	//convierte el archivo en imagen y lo pasa a img
	img, _, err := image.Decode(f)

	//se crea la imagen resultado
	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	wImg := image.NewRGBA(rect)

	// se inicializa la varible para la concurencia
	wg := new(sync.WaitGroup)
	//ejecución secuencial
	start2 := time.Now()

	for x := 0; x < size.X; x++ {
		wg.Add(1)
		x := x
		go func() {
			for y := 0; y < size.Y; y++ {
				pixel := img.At(x, y)
				originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)

				// sumar los pixeles
				r := uint8(float64(originalColor.R))
				g := uint8(float64(originalColor.G))
				b := uint8(float64(originalColor.B))
				a := uint8(float64(originalColor.A))

				c := color.RGBA{
					R: r, G: g, B: b, A: a,
				}

				ArrayR = append(ArrayR, int(r))
				ArrayG = append(ArrayG, int(g))
				ArrayB = append(ArrayB, int(b))

				wImg.Set(x, y, c)
			}
			defer wg.Done()
		}()
	}
	wg.Wait()
	elapsed2 := time.Since(start2)
	log.Printf("El tiempo: %s", elapsed2)

	ext := filepath.Ext(imgPath)
	name := strings.TrimSuffix(filepath.Base(imgPath), ext)
	newImagePath := fmt.Sprintf("%s/%s_fin%s", filepath.Dir(imgPath), name, ext)
	fg, err := os.Create(newImagePath)
	defer fg.Close()
	check(err)
	err = jpeg.Encode(fg, wImg, nil)
	check(err)

	var _, e1 = os.Stat("g.txt")
	if os.IsNotExist(e1) {
		var file, e1 = os.Create("g.txt")
		if e1 != nil {
			return
		}
		defer file.Close()
	}

	var file, e2 = os.OpenFile("g.txt", os.O_RDWR, 0644)
	if e2 != nil {
		return
	}
	defer file.Close()
	for i := 0; i < len(ArrayG); i++ {
		_, e2 = file.WriteString(fmt.Sprint(ArrayG[i]) + "\n")
		if e2 != nil {
			return
		}
	}
	e2 = file.Sync()
	if e2 != nil {
		return
	}

	//-------------------------------------------------------------------------------------------------------------------------------------------

	var _, e3 = os.Stat("r.txt")
	if os.IsNotExist(e3) {
		var file1, e3 = os.Create("r.txt")
		if e3 != nil {
			return
		}
		defer file1.Close()
	}

	var file1, e4 = os.OpenFile("r.txt", os.O_RDWR, 0644)
	if e4 != nil {
		return
	}
	defer file1.Close()
	for i := 0; i < len(ArrayR); i++ {
		_, e4 = file1.WriteString(fmt.Sprint(ArrayR[i]) + "\n")
		if e4 != nil {
			return
		}
	}
	e4 = file1.Sync()
	if e4 != nil {
		return
	}

	//-------------------------------------------------------------------------------------------------------------------------------------------

	var _, e7 = os.Stat("b.txt")
	if os.IsNotExist(e7) {
		var file2, e7 = os.Create("b.txt")
		if e7 != nil {
			return
		}
		defer file2.Close()
	}

	var file2, e6 = os.OpenFile("b.txt", os.O_RDWR, 0644)
	if e6 != nil {
		return
	}
	defer file2.Close()
	for i := 0; i < len(ArrayB); i++ {
		_, e6 = file2.WriteString(fmt.Sprint(ArrayB[i]) + "\n")
		if e6 != nil {
			return
		}
	}
	e6 = file2.Sync()
	if e6 != nil {
		return
	}

}
