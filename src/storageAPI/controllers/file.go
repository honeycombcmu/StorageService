package controllers

import (
	//"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"os"
	"storageAPI/models"
)

// Operations about spark tasks
type FileController struct {
	beego.Controller
}

// @Title upload
// @Description upload file to hdfs
// @Param	body		body 	models.File	true		"file description"
// @Success 200 {string} models.File.Id
// @Failure 403 body is empty
// @router / [post]
func (c *FileController) Post() {
	_, header, err := c.GetFile("files")
	if err != nil {
		fmt.Println("getfile err ", err)
	}
	f, err := header.Open()
	defer f.Close()
	if err != nil {
		return
	}

	//create destination file making sure the path is writeable.
	dir, err := os.Getwd()
	fmt.Println(dir)
	dst, err := os.Create(dir + "/upload/" + header.Filename)
	defer dst.Close()
	if err != nil {
		fmt.Println("here????")
		fmt.Println(err)
		return
	}
	//copy the uploaded file to the destination file
	if _, err := io.Copy(dst, f); err != nil {
		fmt.Println("here?")
		fmt.Println(err)
		return
	}
	file := models.File{}
	fileid := models.AddFile(file)
	c.Data["json"] = map[string]string{"FileName": header.Filename, "FileId": fileid}
	c.ServeJSON()
}
