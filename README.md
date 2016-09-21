# Rayz
Simple raytracer written in Go, for fun and learning

![render](https://raw.githubusercontent.com/KenGorab/raytracer-go/master/examples/three-metal-balls/out.png)
Generated using this ray tracer

## Building from source & running
You'll need:

  - Go v1.7.1+ (can be installed [here](https://golang.org/dl/))
  - [Govendor](https://github.com/kardianos/govendor), for dependency management
    - You can install this with `go get -u github.com/kardianos/govendor`
  - Make (`brew install cmake` on macOS)

Once all dependencies are installed, simply run `make`. It will install and "vendor" all dependencies, and build the `rayz` binary.

The binary can then be executed with 1 argument - the filename of a yaml file representing the scene to be rendered. Once started, it will output the progress until the rendering is complete, outputting `out.ppm` in the root directory.

## Scene config options
The `rayz` binary accepts a yaml file describing the scene to be rendered. Below is an outline of all of the top-level properties:

| Property     | Example   | Type          | Notes                                                      |
|---           |---        |---            |---                                                         |
| `width`      | `1920`    | `Integer`     | The width of the rendered image, in pixels.                |
| `height`     | `1080`    | `Integer`     | The height of the rendered image, in pixels.               |
| `samples`    | `100`     | `Integer`     | The number of samples the renderer should take per pixel   |
| `camera`     | See below | `Camera`      | See the camera config properties below                     |
| `materials`  | See below | `[]Material`  | See the material config properties below                   |
| `objects`    | See below | `[]Object`    | See the object config properties below                     |

### `Camera`

| Property        | Example     | Type         | Notes                                                   |
|---              |---          |---           |---                                                      |
| `position`      | `[0, 1, 2]` | `[3]Integer` | The point (x, y, z) of the camera                       |
| `target`        | `[0, 0, 0]` | `[3]Integer` | The point (x, y, z) where the camera should be pointing |
| `field-of-view` | `75.0`      | `Float`      | The degrees measuring [field of view](https://en.wikipedia.org/wiki/Field_of_view) for the camera |
| `aperture`      | `0.01`      | `Float`      | The [aperture](https://en.wikipedia.org/wiki/Aperture) of the camera lens |

### `Material`

| Property | Example     | Type     | Notes                                                    |
|---       |---          |---       |---                                                       |
| `name`   | `'floor'`   | `String` | The material names can be referenced in object configs   |
| `type`   | `'metal'`   | `String` | Either `lambertian` or `metal` (no others supported yet) |
| `color`  | `'#ff0000'` | `String` | A hex string (with or without the #), representing the intrinsic color of the material. All materials have an intrinsic color. |
| `fuzz`   | `0.0`       | `Float`  | Applicable only when the type is `'metal'`, represents the 'fuzziness' of the material. 0 is fully reflective, 1 is fully light-scattering |


### `Object`

| Property   | Example      | Type         | Notes                                                       |
|---         |---           |---           |---                                                          |
| `name `    | `'ball'`     | `String`     | Helps keep track of objects when creating scenes, also used in error messages and internally |
| `shape`    | `'sphere'`   | `String`     | The shape of the object. Currently, only `sphere` supported |
| `radius`   | `0.5`        | `Float`      | The radius of the object, if it's a sphere                  |
| `origin`   | `[0, 0, 0]`  | `[3]Integer` | The point (x, y, z) of the center of the sphere             |
| `material` | `'floor'`    | `String`     | The name of the material of this object. It must match the name of a material defined in the `materials` top-level array. |

### Example

```yaml
width: 800
height: 400
samples: 500

camera:
  position: [0, 1, 2]
  target: [0, 0, 0]
  field-of-view: 75
  aperture: 0.01

materials:
  - name: floor
    type: lambertian
    color: "#4c4c4c"

  - name: silver
    type: metal
    color: "#cccccc"
    fuzz: 1.0

  - name: gold
    type: metal
    color: "#cc9933"
    fuzz: 0.5

  - name: obsidian
    type: metal
    color: "#191919"
    fuzz: 0.0

objects:
  - name: floor
    shape: sphere
    radius: 250
    origin: [0, -250.5, 0]
    material: floor

  - name: left
    shape: sphere
    radius: 0.5
    origin: [-1, 0, 0]
    material: silver

  - name: right
    shape: sphere
    radius: 0.5
    origin: [1, 0, 0]
    material: gold

  - name: center
    shape: sphere
    radius: 0.5
    origin: [0, 0, -1]
    material: obsidian
```

The above scene definition generated the following rendered image (the same one as at the top of this document):

![render](https://raw.githubusercontent.com/KenGorab/raytracer-go/master/examples/three-metal-balls/out.png)
