# Backend for Paraply

As noted in the main README the backend is created with Go.

## How

Query DMI public radar, which forecasts precipitation for the next hour and use that.

### Alternative

We could use their API and do it the proper way. More information can be found at
[DMI Frie Data](https://www.dmi.dk/frie-data). The project could and should
support both ways, but the first one has priority because it is the hacky way,
which might be a source of learning.

## Dependencies

Beyond Go this project needs

- proj (and `proj-devel`)
