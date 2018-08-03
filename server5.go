package main

import (
	"github.com/labstack/echo"
	"os"
	"io"
	"net/http"
	"io/ioutil"
)


func main() {
	func save(c echo.Context) error{
		root := "./" + "pictures" + "/"
		files, _ := ioutil.ReadDir(root)
		// Source
		for _, f range files{
			names := root + f.Name()
			src, err := names.Open()
			if err != nil{
			return err
		}
			defer src.Close()
			dst, err := os.Create(names.Filename)
			if err != nil{
			return err
		}
			defer dst.Close()
			if _, err = io.Copy(dst, src); err != nil{
			return err
		}
			return c.HTML(http.StatusOK, "<b>Image Uploaded</b>")
		}
	}
}