package scheduler

import (
	//"fmt"
	"proj2/png"
	"sync"
	//"proj2/lock"
)


type bspWorkerContext struct {
	// count of threads
	threads 		int

	// exit flag
	shutdown 		bool

	// new Image flag
	newImage		bool

	// effects to be applied for that image
	effects 		[]interface{}

	// Requests
	requests 		[]map[string]interface{}

	// all image taskss
	imageTasks 		[]*png.ImageTask

	// list of image sections
	imageSections 	[]*png.ImageTask

	// image counter
	imageCounter	 int

	// effect counter
	effectsCounter   int

	// section counter
	sectionCounter 	 int

	// new Image to copy the affects 
	blankImage	 *png.Image

	// lock condition variable
	w			*sync.Cond

	// mutex lock
	mut			*sync.Mutex
}


func NewBSPContext(config Config) *bspWorkerContext {

	req := ReadJson(config)
	imageTasks := CreateAllImageTasks(req)
	
	return &bspWorkerContext{
		threads: 	config.ThreadCount,
		shutdown: 	false,
		newImage:   false,
		requests: 	req,
		imageTasks: imageTasks,
		imageCounter: 	0,
		effectsCounter: 0,
		sectionCounter: 0,
		w: 		sync.NewCond(&sync.Mutex{}),
		mut: 	&sync.Mutex{},
	}
}

func RunBSPWorker(id int, ctx *bspWorkerContext) {

	for {

		// steps
		// step sync
		// global sync

		if id == ctx.threads - 1{
			// manager thread

			// 2. pick one image
			ctx.mut.Lock()
			if !ctx.newImage && 
			len(ctx.imageTasks) != ctx.imageCounter{

				//Loads the png image and returns the image or an error
				pngImg, err := png.Load(ctx.imageTasks[ctx.imageCounter].InPath)
		
				if err != nil {
					panic(err)
				}
				// add the loaded image
				ctx.imageTasks[ctx.imageCounter].Image = pngImg

				// add new canvas to work oon
				blankImage, _ := ctx.imageTasks[ctx.imageCounter].Image.NewImage()
				ctx.blankImage = blankImage

				// 3. divide the image
				ctx.imageSections = ctx.imageTasks[ctx.imageCounter].SplitImage(ctx.threads - 1)

				// get the effect
				ctx.effects = ctx.imageTasks[ctx.imageCounter].Effects

				// we have the new image
				ctx.newImage = true
			}
			ctx.mut.Unlock()
			
			
			// 5. wait till any section gets completed
			ctx.mut.Lock()
			// fmt.Printf("Section Counter %d\n", ctx.sectionCounter)
			// fmt.Printf("Sections LEN %d\n", len(ctx.imageSections))
			// wait until all the sections are compleeted
			if ctx.sectionCounter != len(ctx.imageSections) {
				ctx.mut.Unlock()

				// wait
				ctx.w.L.Lock()
				ctx.w.Wait()
				ctx.w.L.Unlock()
				continue
			}
			ctx.mut.Unlock()


			ctx.mut.Lock()
			// image is finished
			if ctx.effectsCounter == len(ctx.effects) - 1{

				//ctx.rw.RUnlock()

				//ctx.rw.Lock()
				// swap the image once all effects are finished
				ctx.imageTasks[ctx.imageCounter].Image = ctx.blankImage

				// 7. save the image
				ctx.imageTasks[ctx.imageCounter].SaveImageTaskOut()
				ctx.imageCounter += 1

				// new Image is required
				ctx.newImage = false

				// reset counters
				ctx.effectsCounter = 0


			} else {
				
				// ctx.rw.RUnlock()

				// ctx.rw.Lock()

				// one effect is finshed 
				ctx.effectsCounter += 1

				// replace the reconstructed image
				ctx.imageTasks[ctx.imageCounter].Image = ctx.blankImage

				// get a new canvas
				blankImage, _ := ctx.imageTasks[ctx.imageCounter].Image.NewImage()
				ctx.blankImage = blankImage
				
				// swap for the next effect
				ctx.imageTasks[ctx.imageCounter].Image.SwapImage()

				// split the images
				ctx.imageSections = ctx.imageTasks[ctx.imageCounter].SplitImage(ctx.threads - 1)

			}
			// reset counters
			ctx.sectionCounter = 0
			ctx.mut.Unlock()


			// check if all images are done
			ctx.mut.Lock()
			ctx.w.L.Lock()
			if len(ctx.imageTasks) == ctx.imageCounter {
				ctx.shutdown = true
				ctx.mut.Unlock()
				ctx.w.Broadcast()
				ctx.w.L.Unlock()
				break
			}
			ctx.mut.Unlock()
			ctx.w.Broadcast()
			ctx.w.L.Unlock()
		
			// repeat from step 2
		} else {

			// non-manager (servant) work
			
			// awaiting for the manager to divide the work
			ctx.mut.Lock()
			if len(ctx.imageSections) < 1 {
				ctx.mut.Unlock()
				continue
			}
			ctx.mut.Unlock()

			ctx.mut.Lock()
			if	len(ctx.imageSections) == ctx.sectionCounter{
				// update counter
				ctx.sectionCounter = 0
				ctx.mut.Unlock()
				continue
			}
			ctx.mut.Unlock()

			ctx.mut.Lock()
			// 1., 2. get one image-section task and perform
			ApplyEffect(ctx.imageSections[id].Image, ctx.effects[ctx.effectsCounter].(string))
			
			// 3. write its section to a canvas
			// TODO: Send a complete image section inside
			ctx.blankImage.AddSection(
				ctx.imageSections[id].Image, 
				ctx.imageSections[id].YStart,
				ctx.imageSections[id].SectionCase)
			
			ctx.mut.Unlock()
			

			ctx.mut.Lock()
			if len(ctx.imageSections) != ctx.sectionCounter{
				// update counter
				ctx.sectionCounter += 1
			}
			ctx.mut.Unlock()
			

			
			// 5. exit or wait when done
			ctx.w.L.Lock()
			if ctx.shutdown{
				// break the code
				ctx.w.Broadcast()
				ctx.w.L.Unlock()
				break
			}
			ctx.w.Broadcast()
			ctx.w.Wait()
			ctx.w.L.Unlock()
				
		}

		
	}
}
