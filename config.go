package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"text/tabwriter"
	"time"
)

const (
	// epsilon is a very small number that approximates 0 for cases like collision detection.
	epsilon = 1e-8
)

// Config contains the configuration from a JSON file.
type Config struct {
	// Number of CPU cores used.
	// Note that this setting isn't imported from the user settings.
	NumCPU int

	// Name of the output file.
	FileName string `json:"fileName"`

	// Random seed.
	Seed int64 `json:"seed"`

	// Dimensions of the rendered image.
	Width  int `json:"width"`
	Height int `json:"height"`

	// Bounding box of the scene.
	SceneBoundMin [3]float64 `json:"sceneBoundMin"`
	SceneBoundMax [3]float64 `json:"sceneBoundMax"`

	// Camera configurations which include its position, center of the scene, and up direction.
	// Note that its binormal vector is computed internally with the tangent and normal vectors
	// that are computed from center and up.
	CameraEye    [3]float64 `json:"cameraEye"`    // camera position
	CameraCenter [3]float64 `json:"cameraCenter"` // center of the screen
	CameraUp     [3]float64 `json:"cameraUp"`     // camera's up direction

	// Sample size for Monte Carlo integration. The sample size will specify how many rays with
	// random directions are sampled from a ray intersection.
	PixelSampleSize     int `json:"pixelSampleSize"`
	IntersectSampleSize int `json:"intersectSampleSize"`

	// Recursion depth for tracing rays.
	TraceDepth int `json:"traceDepth"`
}

// NewDefaultConfig returns a Config with default configuration.
func NewDefaultConfig() *Config {
	return &Config{
		NumCPU:              runtime.NumCPU(),
		FileName:            fmt.Sprintf("phoebe_%d.png", time.Now().UnixNano()),
		Seed:                0,
		Width:               800,
		Height:              800,
		SceneBoundMin:       [3]float64{-5.0, -5.0, -5.0},
		SceneBoundMax:       [3]float64{5.0, 5.0, 5.0},
		CameraEye:           [3]float64{0.0, 0.0, 0.0},
		CameraCenter:        [3]float64{0.0, 0.0, -1.0},
		CameraUp:            [3]float64{0.0, 1.0, 0.0},
		PixelSampleSize:     9,
		IntersectSampleSize: 4,
		TraceDepth:          3,
	}
}

// NewConfigJSON returns a new configuration, given an argument name of the JSON file.
func NewConfigJSON(fileName string) (*Config, error) {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	config := Config{NumCPU: runtime.NumCPU()}
	if err = json.Unmarshal(dat, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Summary prints the summary of the configuration.
func (c *Config) Summary() {
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintln(w, "=============== Configuration Summary ===============")
	fmt.Fprintln(w, "-----------------------------------------------------")

	// Print the number of CPU cores used by this program.
	fmt.Fprintf(w, "+ Number of CPU cores: \t%d\n", c.NumCPU)
	fmt.Fprintln(w, "-----------------------------------------------------")

	// Print the output file name.
	fmt.Fprintf(w, "+ Output file name: \t%s\n", c.FileName)
	fmt.Fprintln(w, "-----------------------------------------------------")

	// Print the random seed.
	fmt.Fprintf(w, "+ Random seed: \t%d\n", c.Seed)
	fmt.Fprintln(w, "-----------------------------------------------------")

	// Print the output image dimensions.
	fmt.Fprintf(w, "+ Image dimensions: \t(%d, %d)\n", c.Width, c.Height)
	fmt.Fprintln(w, "-----------------------------------------------------")

	// Print the scene boundary settings.
	minX, minY, minZ := c.SceneBoundMin[0], c.SceneBoundMin[1], c.SceneBoundMin[2]
	maxX, maxY, maxZ := c.SceneBoundMax[0], c.SceneBoundMax[1], c.SceneBoundMax[2]
	fmt.Fprintln(w, "+ Scene boundary settings:")
	fmt.Fprintf(w, "  Low: \t(%.3f, %.3f, %.3f)\n", minX, minY, minZ)
	fmt.Fprintf(w, "  High: \t(%.3f, %.3f, %.3f)\n", maxX, maxY, maxZ)
	fmt.Fprintln(w, "-----------------------------------------------------")

	// Print the camera settings.
	x, y, z := c.CameraEye[0], c.CameraEye[1], c.CameraEye[2]
	tx, ty, tz := c.CameraCenter[0], c.CameraCenter[1], c.CameraCenter[2]
	nx, ny, nz := c.CameraUp[0], c.CameraUp[1], c.CameraUp[2]

	fmt.Fprintln(w, "+ Camera settings:")
	fmt.Fprintf(w, "  Eye: \t(%.3f, %.3f, %.3f)\n", x, y, z)
	fmt.Fprintf(w, "  Center: \t(%.3f, %.3f, %.3f)\n", tx, ty, tz)
	fmt.Fprintf(w, "  Up: \t(%.3f, %.3f, %.3f)\n", nx, ny, nz)
	fmt.Fprintln(w, "-----------------------------------------------------")

	// Print the sample size.
	fmt.Fprintf(w, "+ Pixel sample size: \t%d\n", c.PixelSampleSize)
	fmt.Fprintf(w, "+ Intersection sample size: \t%d\n", c.IntersectSampleSize)
	fmt.Fprintln(w, "-----------------------------------------------------")

	// Print the tracing recursion depth.
	fmt.Fprintf(w, "+ Recursion depth: \t%d\n", c.TraceDepth)
	fmt.Fprintln(w, "-----------------------------------------------------")

	fmt.Fprintln(w, "=====================================================")
}

// SceneBound returns the bounding box of the scene.
func (c *Config) SceneBound() *BoundBox {
	boundMin := NewVec3(c.SceneBoundMin[0], c.SceneBoundMin[1], c.SceneBoundMin[2])
	boundMax := NewVec3(c.SceneBoundMax[0], c.SceneBoundMax[1], c.SceneBoundMax[2])
	return NewBoundBox(boundMin, boundMax)
}

// EyeCenterUp returns the position (eye), center coordinates, and up direction.
func (c *Config) EyeCenterUp() (*Vec3, *Vec3, *Vec3) {
	eye := NewVec3(c.CameraEye[0], c.CameraEye[1], c.CameraEye[2])
	center := NewVec3(c.CameraCenter[0], c.CameraCenter[1], c.CameraCenter[2])
	up := NewVec3(c.CameraUp[0], c.CameraUp[1], c.CameraUp[2])
	return eye, center, up
}
