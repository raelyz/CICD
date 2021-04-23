package main





func main() {



	// Possible Location to run goroutines, would require passing the address of the waitgroup through but didn't have time to test for concurrency stability
	// wg.Add(2)
	// go user.MenuOptions()
	// go user.MenuOptions()
	// wg.Wait()
	//


}

func Sum(a,b int) int{

	return a+b
}
