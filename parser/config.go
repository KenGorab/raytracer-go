package parser

import (
	"gopkg.in/yaml.v2"
	"kenrg.co/rayz/primitives"
	"strings"
	"errors"
	"fmt"
	"strconv"
)

type MaterialConfig struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Color string `yaml:"color"`
	Fuzz  string `yaml:"fuzz"`
}

type ObjectConfig struct {
	Name     string `yaml:"name"`
	Shape    string `yaml:"shape"`
	Radius   float64 `yaml:"radius,omitempty"`
	Origin   []float64 `yaml:"origin,omitempty"`
	Material string `yaml:"material"`
}

type CameraConfig struct {
	Position    []float64 `yaml:"position,flow"`
	Target      []float64 `yaml:"target,flow"`
	FieldOfView float64 `yaml:"field-of-view"`
	Aperture    float64 `yaml:"aperture"`
}

type SceneConfig struct {
	Width     int `yaml:"width,omitempty"`
	Height    int `yaml:"height,omitempty"`
	Samples   int `yaml:"samples,omitempty"`
	Materials []MaterialConfig `yaml:"materials,flow"`
	Objects   []ObjectConfig `yaml:"objects,flow"`
	Camera    CameraConfig `yaml:"camera,flow"`
}

func ReadConfigFromFile(fileBytes []byte) SceneConfig {
	sceneConfig := SceneConfig{}
	yaml.Unmarshal(fileBytes, &sceneConfig)
	return sceneConfig
}

func (sc SceneConfig) CameraFromSceneConfig() primitives.Camera {
	cameraConfig := sc.Camera
	position := primitives.Vector{X: cameraConfig.Position[0], Y: cameraConfig.Position[1], Z: cameraConfig.Position[2]}
	target := primitives.Vector{X: cameraConfig.Target[0], Y: cameraConfig.Target[1], Z: cameraConfig.Target[2]}
	fieldOfView := cameraConfig.FieldOfView
	aspectRatio := float64(sc.Width) / float64(sc.Height)
	aperture := cameraConfig.Aperture

	return primitives.NewCamera(position, target, fieldOfView, aspectRatio, aperture)
}

func hexColorToVector(hexStr string) (primitives.Vector, error) {
	hexStr = strings.Replace(hexStr, "#", "", 1)
	if len(hexStr) == 6 {
		color, err := strconv.ParseUint(hexStr, 16, 32)
		if err != nil {
			return primitives.Vector{}, errors.New(fmt.Sprintf("Could not parse color from string: \"%v\"", hexStr))
		}

		r := float64(color >> 16 & 0xff) / 255
		g := float64(color >> 8 & 0xff) / 255
		b := float64(color & 0xff) / 255
		return primitives.Vector{r, g, b}, nil
	} else {
		return primitives.Vector{}, errors.New(fmt.Sprintf("Could not parse color from string: \"%v\"", hexStr))
	}
}

func (conf MaterialConfig) materialFromConfig() (string, primitives.Material, error) {
	name := conf.Name
	color, err := hexColorToVector(conf.Color)
	if err != nil {
		return "", nil, err
	}

	switch conf.Type {
	case "lambertian":
		lambertian := primitives.Lambertian{C: color}
		return name, lambertian, nil
	case "metal":
		fuzz, err := strconv.ParseFloat(conf.Fuzz, 32)
		if err != nil {
			return name, nil, errors.New(fmt.Sprintf("Could not create material with name: \"%v\", invalid fuzz", name))
		}
		metal := primitives.Metal{C: color, Fuzz: fuzz}
		return name, metal, nil
	default:
		return "", nil, errors.New(fmt.Sprintf("Could not create material with type: \"%v\"", conf.Type))
	}
}

func (conf ObjectConfig) objectFromConfig(materials map[string]primitives.Material) (string, primitives.Sphere, error) {
	name := conf.Name
	switch conf.Shape {
	case "sphere":
		if material, ok := materials[conf.Material]; !ok {
			errMsg := fmt.Sprintf("Could not create object with name: \"%v\", material \"%v\" has not been defined", name, conf.Material)
			return "", primitives.Sphere{}, errors.New(errMsg)
		} else {
			center := primitives.Vector{conf.Origin[0], conf.Origin[1], conf.Origin[2]}
			sphere := primitives.Sphere{Center: center, Radius: conf.Radius, Material: material}
			return name, sphere, nil
		}
	default:
		errMsg := fmt.Sprintf("Could not create object with name: \"%v\", shape \"%v\" is not supported", name, conf.Shape)
		return "", primitives.Sphere{}, errors.New(errMsg)
	}
}

func (sc SceneConfig) WorldFromSceneConfig() (primitives.World, error) {
	world := primitives.World{}

	materials := make(map[string]primitives.Material)
	for _, config := range sc.Materials {
		name, material, err := config.materialFromConfig()
		if err != nil {
			return world, err
		}

		materials[name] = material
	}

	for _, obj := range sc.Objects {
		_, object, err := obj.objectFromConfig(materials)
		if err != nil {
			return world, err
		}

		world.Add(&object)
	}

	return world, nil
}
