# VRADiant

### What is it?
VRADiant is an implementation of the Valve RADiosity simulator, written in Golang. For the unfamiliar, VRAD is the tool 
used by Source Engine to pre-compute light in .bsp maps. 
VRADiant is ported on the latest Source SDK base (2013 at the time of writing).

### Why? 
Source Engine is huge and unwieldy. Different components are deeply entrenched in the repository. This is an effort to 
isolate and extract a single coherent section of the codebase into its own, independent application. VRAD is also an aging application that could see drastic improvements with newer technologies.

### Current state
In its current state, this is more or less a line-by-line port, with some effort to organise functions and collections into something that can be refactored more easily in the future. As such, the code should look quite familiar if you compare the original C++ to this.
There are plenty of bad practices visible in this repo. Largely due to the scale, and my unfamiliarity with the inner workings of vrad. The plan is to refactor this implementation entirely once it's in a working state.
Also note the lack of tests. yeah that sucks. The issue comes down to writing tests for enormous functions that rely on huge global datastructures. With refactoring should come good test coverage. In the meantime, judge me for not writing any.


### Roadmap
- [x] Load internal bsp data structures into memory
- [x] Prepare a environment from which a raytracer can be run against
- [ ] Port the Radiosity process (top level C functions are RadWorld_*())
- [ ] Refactor data import and raytracing environment preparation
- [ ] Investigate GPU accleration options during radiosity calculations.

#### Misc
This is (certainly at this point in time) a port of Source Engine 2013 base code, which is freely available on Github.
