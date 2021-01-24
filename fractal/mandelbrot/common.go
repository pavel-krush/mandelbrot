package mandelbrot

const (
	// Initial physical width and height
	// Used to calculate scale ratio
	DefaultPhysWidth = 3.0
	DefaultPhysHeight = 2.0

	// Number of iterations for each point
	DefaultIterations = 256

	// Point does belong to mandelbrot set if value is less than threshold
	DefaultThreshold = 4.0
)
