package scheduler

import (
	//"fmt"
	//"fmt"
	"math"
	"proj2/png"
	"sync"
)


type bspWorkerContext struct {
	// count of threads
	threads 		int

	// exit flag
	shutdown 		bool

	// new Image flag
	newImage		bool

	// Requests
	requests 		[]map[string]interface{}

	// all image taskss
	imageTasks 		[]*png.ImageTask

	// list of image sections
	imageSections 	[][]int

	// image counter
	imageCounter	 int

	// effect counter
	effectsCounter   int

	// section counter
	sectionCounter 	 int

	// lock condition variable
	w			*sync.Cond

	// mutex lock
	mut			*sync.Mutex
}


func NewBSPContext(config Config) *bspWorkerContext {

	req := ReadJson(config)
	imageTasks := CreateAllImageTasks(req)
	// size of each section
	sectionSize := int(math.Ceil(float64(len(imageTasks))/float64(config.ThreadCount - 1)))
	// threads := config.ThreadCount
	// if sectionSize < 1{
	// 	sectionSize = 1
	// 	threads = len(imageTasks)
	// }
	imageSections := make([][]int, config.ThreadCount - 1)
	counter:= 0
	for i, _:= range imageSections{
	imageSections[i] = make([]int, 0)
	}

	for j:=0; j < sectionSize; j++ {
		for i, _:= range imageSections{		
				if counter < len(imageTasks){
					imageSections[i] = append(imageSections[i], counter)
					counter++
				} else {
					break
				}
			}
	}
	// fmt.Println(sectionSize)
	// fmt.Println(imageSections)

	return &bspWorkerContext{
		threads: 	config.ThreadCount,
		shutdown: 	false,
		newImage:   false,
		requests: 	req,
		imageTasks: imageTasks,
		imageSections:  imageSections,
		imageCounter: 	0,
		effectsCounter: 0,
		sectionCounter: 0,
		w: 		sync.NewCond(&sync.Mutex{}),
		mut: 	&sync.Mutex{},
	}
}

func RunBSPWorker(id int, ctx *bspWorkerContext) {

	effectsCounter := 0
	threadImageCounter := 0
	currImage := 0
	imageDone := false
	for {

		// steps
		// step sync
		// global sync
		
		if id == ctx.threads - 1 || id >= len(ctx.imageTasks){
			// manager thread

			// check if all images are done
			ctx.w.L.Lock()
			if len(ctx.imageTasks) == ctx.imageCounter {
				// fmt.Println("Shutdowwwn*******")
				ctx.shutdown = true
				ctx.w.Broadcast()
				ctx.w.L.Unlock()
				break
			}
			ctx.w.Wait()
			ctx.w.L.Unlock()
			
		
			// repeat from step 2
		} else {

			// non-manager (servant) work
			// awaiting for the manager to divide the work
			ctx.w.L.Lock()
			if len(ctx.imageTasks) == ctx.imageCounter{
				// break the code
				// fmt.Println("Shutdowwwn*******")
				ctx.w.Broadcast()
				ctx.w.L.Unlock()
				break
			}
			// barrier
			if threadImageCounter == len(ctx.imageSections[id]) &&
				imageDone {	
				ctx.w.Wait()
				ctx.w.L.Unlock()
				continue
			}
			ctx.w.L.Unlock()
			
			// first image
			if currImage == 0 && threadImageCounter == 0 {
				currImage = ctx.imageSections[id][threadImageCounter]
				threadImageCounter ++

			}
			
			if effectsCounter == len(ctx.imageTasks[currImage].Effects){
				// fmt.Printf(" Saving Thread id: %d, workign on %d\n", id, currImage)
				// fmt.Printf(" Saving image: %d\n", currImage)
				ctx.imageTasks[currImage].SaveImageTaskOut()
				imageDone = true

				if threadImageCounter + 1  < len(ctx.imageSections[id]) {
					threadImageCounter += 1
					currImage = ctx.imageSections[id][threadImageCounter]					
					effectsCounter = 0
					imageDone = false
				}

				ctx.w.L.Lock()
				ctx.imageCounter += 1
				
				//
				//ctx.w.Wait()
				ctx.w.L.Unlock()
				

			} else {
				// fmt.Printf("Thread id: %d, workign on %d\n", id, currImage)
				// fmt.Printf("Thread id: %d, workign on %d, effect: %s\n", id, currImage, ctx.imageTasks[currImage].Effects[effectsCounter].(string))
				// 1., 2. get one image-section task and perform
				ApplyEffect(ctx.imageTasks[currImage].Image, ctx.imageTasks[currImage].Effects[effectsCounter].(string))
				
				effectsCounter += 1

				if effectsCounter != len(ctx.imageTasks[currImage].Effects){
					// swap for the next effect
					ctx.imageTasks[currImage].Image.SwapImage()
				}
			}
		}
	}
}
