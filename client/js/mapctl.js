(function() {
    

    mapctl   = function(){
        this.server      = 'http://localhost:8000/';
        this.requests    = { 'wkb': null, 'wkt': null, 'wasm': null };
        this.data_proj   = 'EPSG:4326';
        this.map_proj    = 'EPSG:4326';
        this.log_element = null;
    };

    mapctl.prototype.show_map = function( target, type ){
        
        if( _.has( this.requests, type) && this.requests[type] == null ){
            var timer  = target.getElementsByClassName( 'timer' );
            var mapdiv = target.getElementsByClassName( 'map'   );

            console.log( 'implement timer!' );

            if( timer.length * mapdiv.length ){

                timer[0].innerHTML  = '';
                mapdiv[0].innerHTML = '';

                this.requests[type] = performance.now();
                this.log( 'Initializing ' + type + 'server requests' );
                this.get_data( this.server, type )
                    .then(function( xhr ) { 
                        
                        this.log( 'Response form server took: ' + (performance.now() - this.requests[type]) + 'ms; status ok' );
                        var parse_performance = performance.now();
                        var features = this['parse' + type]( xhr.response );
                        this.log( 'Parsing features took: ' + (performance.now() - parse_performance) + "ms; got " + features.length + " features " );

                        var map = new ol.Map({
                            target: mapdiv[0],
                            layers: [ 
                                new ol.layer.Tile({
                                    source: new ol.source.OSM({})
                                }), 
                                new ol.layer.Vector({
                                    renderMode: 'image',
                                    source: new ol.source.Vector({
                                        features: features
                                    })
                                }) ],
                            view:  new ol.View({
                                center: features.reduce( function(c, f){
                                    var extent = f.getGeometry().getExtent();
                                    c[0] = (c[0] + ((extent[0] + extent[2]) * .5) ) * .5;
                                    c[1] = (c[1] + ((extent[1] + extent[3]) * .5) ) * .5;
                                    return c;
                                }, [0, 0]),
                                zoom: 1,
                                projection: this.map_proj
                            })
                        });

                        map.on( 'click', function( evt ){ 
                            var features = map.getFeaturesAtPixel( evt.pixel );
                            if( features ){
                                features.forEach( function(f) {
                                    var id = f.getProperties().id;
                                    window.mapctl.log( 'You have clicked on feature with ID: ' + id );
                                    console.log( 'You have clicked on feature with ID: ',  id );
                                } );
                            }
                        });

                    }.bind( this ))
                    .catch(function(ex) { 
                        this.log( 'Response form server took: ' + (performance.now() - this.requests[type]) + 'ms; status error' );
                        console.log( ex )
                    }.bind( this ))
                    .then(function() { 
                        this.log( 'Whole process took: ' + (performance.now() - this.requests[type]) + 'ms' );
                        this.requests[type] = null;
                        // console.log('finally call'); 
                    }.bind( this ));

            }
            else {
                this.log( 'Unexpected response structure' );
            }

        }
        else {
            this.log( 'At the moment request cannot be performed!' );
        }

    };

    mapctl.prototype.init_map = function( container, features ){
        var map = new ol.Map({
            target: container,
            layers: [ new ol.layer.Tile({ source: new ol.source.OSM() }) ]
        });
        return map;
    };

    mapctl.prototype.parsewasm = function( data ){
        var ret = [];
        var wkb = wkb_format();
        var pos   = 0;
        while( pos < data.byteLength ){
            var buf = new Uint32Array( data.slice( pos, pos+8 ) );
            pos += 8;
            if( buf[1] ){
                
                var id = buf[0];

                // console.log( 'wasam start' );

                var arr = new Uint8Array( data.slice( pos, pos+buf[1] ) );
                var wa_buff = Module._malloc( arr.byteLength );
                Module.writeArrayToMemory( arr, wa_buff );
                var typ   = Module.ccall( 'type', 'string', ['arraybuffer'], [wa_buff] );
                var ngeom = Module.ccall( 'geomN', 'number', ['arraybuffer'], [wa_buff] );

                if( typ == 'multipolygon' && ngeom > 0){

                    // TODO: we will need array for floats as well!?
                    var polyn = Module._malloc( Uint32Array.BYTES_PER_ELEMENT * ngeom );
                    var rings = Module.ccall( 'convert', 'number', ['arraybuffer', 'arraybuffer'], [wa_buff, polyn]);
                    if( rings ){

                        var coords = new Array( ngeom );
                        for( var i=0; i < ngeom; i++ ){
                            // console.log('get method: ', Module.getValue( polyn + i, 'i8' ) );
                            // console.log( i, '.direct heap method: ', Module.HEAP8[ (polyn + i) / Uint8Array.BYTES_PER_ELEMENT ]);
                            var npts = Module.HEAPU32[ (polyn / Uint32Array.BYTES_PER_ELEMENT) + i ];
                            var ringptr = Module.HEAPU32[ (rings / Uint32Array.BYTES_PER_ELEMENT) + i ];

                            // console.log( npts ); // ok npts is still a bit off!
                            coords[i] = [ new Array( npts ) ];
                            for( var j=0; j < 2 * npts; j+=2 ){

                                coords[i][0][j / 2] = [
                                    Module.HEAPF64[ (ringptr / Float64Array.BYTES_PER_ELEMENT) + j ],
                                    Module.HEAPF64[ (ringptr / Float64Array.BYTES_PER_ELEMENT) + j + 1 ]
                                ];
                                // console.log( 'first value:', Module.HEAPF64[ (ringptr / Float64Array.BYTES_PER_ELEMENT) + j ] );
                                // console.log( 'first value:', Module.HEAPF64[ (ringptr / Float64Array.BYTES_PER_ELEMENT) + j + 1 ] );
                            }

                            Module._free( ringptr );
                        }
                        
                        // now construct the feature
                        var f = new ol.Feature({
                            id: id,
                            geometry: new ol.geom.MultiPolygon( coords )
                        });
                        ret.push( f );
                    }
                    
                    Module._free( polyn );

                }
                else {
                    console.warn( 'Experimental support is available for multipolygon geometry type only!' );
                }

                Module._free( wa_buff );
                pos  += buf[1];
            }
            // return ret; // debug
        }
        return ret;
    };

    mapctl.prototype.parsewkt = function( data ){
        data = JSON.parse( data );
        var ret = [];
        var format = new ol.format.WKT();

        _.forEach( data, function(record){
            var f = format.readFeature( record.WKT , {
                dataProjection: this.data_proj,
                featureProjection: this.map_proj
            } );
            f.setProperties( {id: record.Id} );
            ret.push( f );
        });

        return ret;
    };

    mapctl.prototype.parsewkb = function( data ){
        var ret = [];
        var wkb = wkb_format();
        var i   = 0;
        while( i < data.byteLength ){
            var buf = new Uint32Array( data.slice( i, i+8 ) );
            i += 8;
            if( buf[1] ){
                
                var id = buf[0];
                var dw = new DataView( data.slice( i, i + buf[1] ) );
                wkb.parse( data.slice( i, i + buf[1] ) );
                if( wkb.type == 'multipolygon' ){
                    var f = new ol.Feature({
                        id: id,
                        geometry: new ol.geom.MultiPolygon( wkb.coords )
                    });
                    ret.push( f );
                }
                // console.log( id, wkb.type, wkb.coords );
                i += buf[1];
            }
        }
        return ret;
    };

    mapctl.prototype.get_data = function( url, type ) {
        return new Promise( (success, error) => {
            
            var req = new XMLHttpRequest();

            req.responseType = (type == 'wkt') ? 'JSON' : 'arraybuffer';
            req.onload  = () => success( req );
            req.onerror = () => error( null );

            try{
                req.open( "GET", url + type, true );
                req.send();
            }
            catch (ex){
                error( ex );
            }

        });
    };

    mapctl.prototype.log = function( message ) {
        if( this.log_element ){            
            var a = document.createElement( 'li' );
            a.appendChild( document.createTextNode( message ) );
            if( this.log_element.firstElementChild ){
                this.log_element.insertBefore( a, this.log_element.firstElementChild );
            }
            else{ 
                this.log_element.appendChild( a );
            }
        }
        else {
            console.log( message );
        }
    };

    mapctl.prototype.clear_log = function() {
        if( this.log_element ){
            this.log_element.innerHTML = '';
        }
        else {

        }
        console.log( 'clear messages!' );
    };

    window.mapctl = new mapctl();

})()