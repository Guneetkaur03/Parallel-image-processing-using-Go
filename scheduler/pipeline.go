package scheduler


import (
	//"fmt"
	"proj2/png"
)


/*
	Method for each worker
	each worker divides the image into mutiple sections and 
	processes each section by a goroutine
*/
func worker(threads int, doneWorker chan<- bool, imageTasks <-chan *png.ImageTask, imageResults chan<- *png.ImageTask) {

	// channel to be fed for section generation
	SectionsGenerator := func(
		done <- chan interface{}, 
		imageSections []*png.ImageTask,
	) chan *png.ImageTask {
		sectionChannel := make(chan *png.ImageTask)

		// go method to feed to section channel
		go func() {
			defer close(sectionChannel)
			for _, sec := range imageSections {
				select {
					case <-done: return
					case sectionChannel <- sec:
				}
			}
		}()

		return sectionChannel
	}

	// Effects channel, for each section apply the requested effect
	effects := func(
		threads int, 
		effect string, 
		effectNum int, 
		done <- chan interface{}, 
		sectionChannel chan *png.ImageTask,
	) chan *png.ImageTask {

		// get sections equal to num_of_threads
		effectsChannel 	:= make(chan *png.ImageTask, threads) 
		doneSection 	:= make(chan bool)

		// loop through all the sections
		for section := range sectionChannel {

			if effectNum > 0 {
				// swap output image
				section.Image.SwapImage() 
			}

			// go routine for each section
			go func(section *png.ImageTask) { 
				select {
				case <-done: doneSection <- true
					return
				case effectsChannel <- section:
					ApplyEffect(section.Image, effect)
					
					// done with this section
					doneSection <- true 
				}
			}(section)
		}

		// wait for all the sections to complete
		for i := 0; i < threads; i ++ {
			<-doneSection
		}

		return effectsChannel
	}



	for {
		// get tasks
		imageTask, more := <-imageTasks

		// not any more tasks remain
		if !more{
			doneWorker <- true 
			return
		}

		imageChunks := imageTask.SplitImage(threads)
		
		// Iterate through filters on image chunks.
		// Will spawn a goroutine for each chunk in each filter (n goroutines per filter)
		done := make(chan interface{})
		defer close(done)
		sections := SectionsGenerator(done, imageChunks)
		for i := 0; i < len(imageTask.Effects); i++ {
			effect := imageTask.Effects[i].(string)
			sections = effects(threads, effect, i, done, sections)
			close(sections)
		}

		// Put the image back together.
		reconstructedImage, _ := imageTask.Image.NewImage()
		for imageSections := range sections {
			reconstructedImage.AddSection(
				imageSections.Image, 
				imageSections.YStart, 
				imageSections.SectionCase)
		}
		imageTask.Image = reconstructedImage
		imageResults <- imageTask // Send image to results channel to be saved.
	}
}

func RunPipeline(config Config) {

	// thread count
	threads 	 := config.ThreadCount
	// tasks and results channel
	imageTasks 	 := make(chan *png.ImageTask)
	imageResults := make(chan *png.ImageTask)
	// done channels
	doneWorker := make(chan bool)
	doneSave   := make(chan bool)

	// parse the tasks
	go func() {
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

				imageTask := &png.ImageTask{
					InPath:  filePath,
					OutPath: task["outFiles"].([]string)[i],
					Image: 	 pngImg,
					Effects: task["effects"].([]interface{}),
				}
				// feed to generator channel
				imageTasks <- imageTask
			}
		}
		close(imageTasks)
	}()

	// spawn workers to handle each image
	for i := 0; i < threads; i++ {
		go worker(threads, doneWorker, imageTasks, imageResults)
	}

		// Save results.
		go func() {
			for { // Do while there are more images to save.
				imgTask, more := <- imageResults // Reads from the image results channel.
				if more {
					imgTask.SaveImageTaskOut()
				} else {
					doneSave <- true 
					return
				}
			}
		}()

		// Wait for all workers to return.
		for i := 0; i < threads; i++ {
			<-doneWorker
		}

		// Wait for all images to be saved.
		close(imageResults)
		<- doneSave

	

	//fmt.Println(tasks)
}


