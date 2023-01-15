package scheduler

import (
	"proj2/png"
)

func RunSequential(config Config) {
	
	// get the tasks parsed
	tasks := ReadJson(config)

	// loop in all the tasks
	for _, task := range tasks{
		// for each folder
		for i, filePath:= range task["inFiles"].([]string) {

			//Loads the png image and returns the image or an error
			pngImg, err := png.Load(filePath)
	
			if err != nil {
				panic(err)
			}
			
			for i, v := range task["effects"].([]interface{}){

				if i > 0{
					// swap image pointers for all the next effects
					pngImg.SwapImage()
				}

				ApplyEffect(pngImg, v.(string))

			}

			//Saves the image to a new file
			err = pngImg.Save(task["outFiles"].([]string)[i], true)
	
			//Checks to see if there were any errors when saving.
			if err != nil {
				panic(err)
			}
		}
	}
}
