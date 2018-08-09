package main

import (
	"io"
	"net/http"
	"os"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"io/ioutil"
	"bitbucket.org/avcl/pik"
	"bytes"
)

func upload(c echo.Context) error {
	//name := c.FormValue("name")
	//email := c.FormValue("email")

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

	}
	image_list := pickcl()
	var str bytes.Buffer
	for _, l := range image_list{
		str.WriteString(l)
	}
	return c.HTML(http.StatusOK, string(str.String()))
}

func pickcl() []string{
	root := "./"
	pf := photoFactory{}
	contextId := pik.NewDesktopContext()
	files, err := ioutil.ReadDir(root)
	if err != nil{
		log.Fatal(err)
	}
	var photos []*Photo
	for _, f := range files {

		names := root + f.Name()
		if pf.accept(names) {
			photo := pf.CreatePhoto(names)
			pik.AddPhoto(contextId, photo.relativePath, photo.Created, photo.Lat, photo.Lon)
			photos = append(photos, photo)
		}
	}
	pik.StartSmartSelection(contextId)
	image_list := []string{}
	for _, f := range photos {
		b, _ := pik.IsSelected(contextId, f.relativePath)
		//fix to display HTML Image with a star next to it after it is printable
		if b {
			image_list = append(image_list, f.Name+"is_selected, ")
		} else {
			image_list = append(image_list, f.Name+"is_not_selected, ")
		}
	}
	return image_list
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())


	e.Static("/", "public")
	e.POST("/upload", upload)

	e.File("/", "public/index.html")

	e.Logger.Fatal(e.Start(":1323"))
}