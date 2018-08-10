package main

import (
	"io"
	"net/http"
	"os"
	"github.com/labstack/echo"
	"log"
	"io/ioutil"
	"bitbucket.org/avcl/pik"
	"bytes"
	"github.com/labstack/echo/middleware"
)

func upload(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")

	if name == "" {
		return c.HTML(http.StatusOK, "<body>Please enter a valid name</body>")
		panic("error")
	}

	if email == "" {
		return c.HTML(http.StatusOK, "<body>Please enter a valid email</body>")
		panic("error")
	}

	New_Dir_Name := "name="+name+"-"+"email="+email
	os.Mkdir("public/"+New_Dir_Name,0777)

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
		//need to create all imgs in the new dir above
		dst, err := os.Create("./public/"+New_Dir_Name+"/"+file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

	}
	image_list := pickcl(New_Dir_Name)
	var str bytes.Buffer
	for _, l := range image_list{
		str.WriteString(l)
	}
	//os.RemoveAll("./public/"+New_Dir_Name)
	return c.HTML(http.StatusOK, string(str.String()))
}

func pickcl(New_Dir_Name string) []string{
	root := "./public/"+New_Dir_Name+"/"
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
			image_list = append(image_list, "image "+f.Name+":    ")
			image_list = append(image_list, "<p><img src='"+New_Dir_Name+"/"+f.Name+"' style='width:100px;height:100px;'>")
			image_list = append(image_list, "<img src='Star.jpg' style='width:100px;height:100px;'></p>")
		} else {
			image_list = append(image_list, "image "+f.Name+":    ")
			image_list = append(image_list, "<p><img src='"+New_Dir_Name+"/"+f.Name+"' style='width:100px;height:100px;'></p>")
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