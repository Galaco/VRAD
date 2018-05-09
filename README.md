# VRADiant

### What is it?
VRADiant is an implementation of the Valve RADiosity simulator, written in Golang. For the unfamiliar, VRAD is the tool 
used by Source Engine to pre-compute light in .bsp maps. 
VRADiant is ported on the latest Source SDK base (2013 at the time of writing).

### Why? 
Source Engine is huge and unwieldy. Different components are deeply entrenched in the repository. This is an effort to 
isolate and extract a single coherent section of the codebase into its own, independent application.


### Roadmap
- [x] Load internal bsp data structures into memory
- [x] Prepare a environment from which a raytracer can be run against
- [ ] Port the Radiosity process (top level C functions are RadWorld_*())
- [ ] Investigate GPU accleration options during radiosity calculations.

#### Misc
This is (certainly at this point in time) a port of Source Engine 2013 base code, which is freely available on Github.