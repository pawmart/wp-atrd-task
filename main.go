package main

func main() {
	newApp := NewApp()
	newApp.AddHandlers()
	newApp.ListenAndServe()
}
