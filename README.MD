WKB client geometry showcase
============================


TODO: write something about it.

and most important WHY THIS WAS WRITTEN



In order to transfer geometry data in wbk format, a custom binary format was used. As geometry without any metadata in client isn't to much of use, in front of geometry a small amount of data, representing geometry id from database, was added. Due to design of current networking structure, sending a single geometry per connection would be resource wasteful approach. To address this issue and send collection of geometries down the wire, in front of wkb record a length of the same has been added.

Diagram of single wkb feature structure:
|----|----|----|----||----|----|----|----||----|----| ... |----|----|
  geom id (int32)     geometry len (int32)        wkb geometry



By packing geometries as tightly as possible a huge improvement in response time is expected. Additionally, amount of data is reduced allowing for transfer of lager datasets.
TODO: something on parsing the data in the client!

Main point of showcase is to measure perfomance in baseline to same data being delivered using more custom text protocol (not geojson + reason)
Last, a new technology web assembly ...
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