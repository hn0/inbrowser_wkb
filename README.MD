WKB client geometry showcase
============================

This is a simple, visual and hopefully dynamic enough, showcase of benefits of transferring geometry data (not really limited to) to modern browsers in novel old binary way. At the moment of writing this showcase a common practice of transferring GIS data to the browser is by using textual files in JSON defined structured way (often using geoJSON format). As we all know, transferring data in textual representation, all thought convenient for development, adds a considerable amount of overhead.

For reduction in total response size we already have a variety of the methods developed, from compressing data to images using server side compression to techniques like scale dependent visibility where number of features is controlled depending on the scale of the map. Each of current solutions has its drawbacks. 

Server side rendering, even if offering supreme level of data compression isn't suitable when user map interactivity and responsiveness is requirement. Additionally, need for fetching new image on each extent change can, after longterm users map usage, lead to overall higher network utilization together with possible high loads on the server side often present in the case of simultaneously multiple users.


As solution to those issues, preloading of all of the geometries in vector representation is used. With current CONTINUE ...
Partially this can be solved (depending on dataset size) GEOJSON and initial load times



In order to transfer geometry data in wbk format, a custom binary format was used. As geometry without any metadata in client isn't to much of use, in front of geometry a small amount of data, representing geometry id from database, was added. Due to design of current networking structure, sending a single geometry per connection would be resource wasteful approach. To address this issue and send collection of geometries down the wire, in front of wkb record a length of the same has been added.

Diagram of single wkb feature structure:
|----|----|----|----||----|----|----|----||----|----| ... |----|----|
  geom id (int32)     geometry len (int32)        wkb geometry



By packing geometries as tightly as possible a huge improvement in response time is expected. Additionally, amount of data is reduced allowing for transfer of lager datasets.
TODO: something on parsing the data in the client!

Main point of showcase is to measure performance in baseline to same data being delivered using more custom text protocol (not geojson, this showcase places emphasis on amount of geometry data) that transfers exact same amount of information (id and geometry in WKT format).


Last, a currently experimental technology, web assembly was used for attempt of geometry parsing time reduction. All thought, this technology could potentially provide huge benefits in NETWORING?!
current focus was on geometry format, so same was used only for parsing wkb geometries and determination of improvements in respect to same algorithm implemented in javascript.

Live demo available at: stocar.blogsite.org#wkb



Dependencies
------------

Backend is written using go framework (https://golang.org/), while currently supported datasource is limited to sqlite (http://sqlite.org/) database (not spatialite variety). For correct execution of the backend additionally GDAL (http://gdal.org/) libraries must be present and properly configured on the system.

Client is written using emca-6 scripting language and for proper execution requires one of the modern browsers is required (eg. Chromium ver. > 63). Experimental webassembly wkb parser is developed using Firefox (ver. 58) browser and was not tested in other browsers.

In order to use preprepared server runtime environment installation of docker is necessary as well. Current support for environment startup script is limited to the linux platform.


Running the test on local machine
---------------------------------

After cloning project from git pages, prepared environment (requires docker) can be started on local machine using start_server.sh script

```code
    ./start-server.sh
```

after which test web pages can be accessed using following uri: localhost:8000


Using custom datasets
---------------------

Currently, not supported.


Sample data
-----------

Sample data used in and distributed with this showcase is due to courtesy of Federal Geographic Data Committee [url:'http://maps.tnc.org/gis_data.html']. 
Sample data describes Terrestrial ecosystems (tnc_terr_ecoregions) published by The Nature Conservancy (Arlington, VA)
ref: Olson, D. M. and E. Dinerstein. 2002. The Global 200: Priority ecoregions for global conservation. (PDF file) Annals of the Missouri Botanical Garden 89:125-126. -The Nature Conservancy, USDA Forest Service and U.S. Geological Survey, based on Bailey, Robert G. 1995. Description of the ecoregions of the United States (2nd ed.). Misc. Pub. No. 1391, Map scale 1:7,500,000. USDA Forest Service. 108pp. -The Nature Conservancy (2003), based on Wiken, E.B. (compiler). 1986. Terrestrial ecozones of Canada. Ecological Land Classification Series No. 19. Environment Canada, Hull, Que. 26 pp. + map.


After-notes
-----------

Results of the showcase script might differ depending on sample dataset used, quality of client-server connection and environment on which client page is executed.