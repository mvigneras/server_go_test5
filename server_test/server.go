package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"io/ioutil"
	"bitbucket.org/avcl/pik"
)

func upload(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")

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

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>Uploaded successfully %d files with fields name=%s and email=%s.</p>", len(files), name, email))
}

func pickcl(c echo.Context) error{
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

	for _, f := range photos {
		b, _ := pik.IsSelected(contextId, f.relativePath)
		log.Println(f.relativePath, b)
		//fix to display HTML Image with a star next to it after it is printable
		if b {
			return c.HTML(http.StatusOK, fmt.Sprintf("<p>Image with name=%s was selected by the selection.</p>", f.Name))
		} else {
			return c.HTML(http.StatusOK, fmt.Sprintf("<p>Image with name=%s was not selected by the selection.</p>", f.Name))
		}
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())


	e.Static("/", "public")
	e.POST("/upload", upload)
	e.POST("/upload", pickcl)

	e.Logger.Fatal(e.Start(":1323"))
}