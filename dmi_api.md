# The DMI API
You can get access to their API and documentation at
https://opendatadocs.dmi.govcloud.dk/DMIOpenData,
however I also want to try to decode their web app and use that instead. That
way everyone using this does not depend on my or their own API key.

I am going to be using the `DMI Radar map`. We are interested in when it rains,
which is the default mode it is in upon opening the radar. It could look like

<p align="center">
  <img src="/images/map.png">
</p>

To help recognize the pattern between the API calls, I have recorded some of them below.
- https://www.dmi.dk/ZoombareKort/map?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&FORMAT=image/png&TRANSPARENT=true&TIME=2024-11-08T18:35:00Z&LAYERS=radar&WIDTH=512&HEIGHT=512&SRS=EPSG:3575&STYLES=&BBOX=-406250,-4218750,250000,-3562500

It seems to me that `TIME` and `REFERENCE_TIME` are the two of which values
change. `BBOX` is the bounding box and the coordinates are in the style of
EPSG:3575, but we can determine the normal longitude and latitude from them (
https://developers.auravant.com/en/blog/2022/09/09/post-3/). The bounding box
determines which part of the map is shown.
