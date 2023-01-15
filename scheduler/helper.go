package scheduler

import (
	"fmt"
	"os"
	"io"
	"encoding/json"
	"strings"
	"path"
	"proj2/png"
)

/*
	method to read and parse json file
*/
func ReadJson(config Config) []map[string]interface{}  {

	effectsRootPath := "../data"

	// read the file
	effectsFile, _ := os.Open(path.Join(effectsRootPath, "effects.txt"))
	defer effectsFile.Close()

	
	// add a decoder
	reader 	 := json.NewDecoder(effectsFile)
	requests := make([]map[string]interface{}, 0)
	req 	 := make(map[string]interface{})

	for{
		if err := reader.Decode(&req); err != nil {
			if err == io.EOF {
				// read until eof
				break
			}
			fmt.Println(err)
			//return
			continue
		}
		// break the dirs
		req["inFiles"]  = make([]string, 0)
		req["outFiles"] = make([]string, 0)
		for _, v := range strings.Split(config.DataDirs, "+") {
			req["inFiles"] = append(req["inFiles"].([]string), path.Join(effectsRootPath, "in", v, string(req["inPath"].(string))))
			req["outFiles"] = append(req["outFiles"].([]string), path.Join(effectsRootPath, "out", v + "_" + string(req["outPath"].(string))))
		}
		 
		// append the object
		requests = append(requests, req)
		req = make(map[string]interface{})
	}

	return requests
}


/* 
	Perform affect
*/

func ApplyEffect(pngImg *png.Image, effect string)  {
	
	if effect == "B"{
		// perfoms blur
		pngImg.Blur()
	} else if effect == "E"{
		// edge detection
		pngImg.EdgeDetection()
	} else if effect == "S"{
		// performs Sharpen
		pngImg.Sharpen()
	} else if effect == "G"{
		//Performs a grayscale filtering effect on the image
		pngImg.Grayscale()
	}
}


/* 
	Create all image tasks

*/
func CreateAllImageTasks(tasks []map[string]interface{}) []*png.ImageTask  {

	imageTasks := make([]*png.ImageTask, 0)
	// loop in all the tasks
	for _, task := range tasks{
		// for each folder
		for i, filePath:= range task["inFiles"].([]string) {
			pngImg, err := png.Load(filePath)
		
			if err != nil {
				panic(err)
			}
			imageTask := &png.ImageTask{
				InPath:  filePath,
				OutPath: task["outFiles"].([]string)[i],
				Image: pngImg,
				Effects: task["effects"].([]interface{}),
			}
			// add to the list
			imageTasks = append(imageTasks, imageTask)
		}
	}
	
	return imageTasks
}