(function() {
    

    mapctl   = function(){
        this.server    = 'http://localhost:8000/';
        this.requests  = { 'wkb': null, 'wkt': null };
        this.data_proj = 'EPSG:3857';
        this.map_proj  = 'EPSG:4326';
    };

    mapctl.prototype.show_map = function( target, type ){
        
        if( _.has( this.requests, type) && this.requests[type] == null ){
            var timer  = target.getElementsByClassName( 'timer' );
            var mapdiv = target.getElementsByClassName( 'map'   );

            if( timer.length * mapdiv.length ){

                timer[0].innerHTML  = '';
                mapdiv[0].innerHTML = '';

                // // this would be better suited for finally call!?
                // var map = this.init_map( mapdiv[0] );

                this.requests[type] = performance.now();
                this.get_data( this.server + type )
                    .then(function( xhr ) { 


                        // now parsing of the data
                        var features = this['parse' + type]( xhr.response );

                        console.log( features )
                        this.requests[type] = null;

                    }.bind( this ))
                    .catch(function(ex) { console.log('error call', ex); })
                    .then(function() { console.log('finally call'); })

            }
            else {
                this.log( 'unexpected structure' );
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

    mapctl.prototype.parsewkt = function( data ){
        data = JSON.parse( data );
        var ret = [];
        var format = new ol.format.WKT();

        _.forEach( data, function(record){
            var f = format.readFeature( record.WKT , {
                dataProjection: this.data_proj,
                featureProjection: this.map_proj
            } );
            f.setProperties( {id: record.ID} );
            ret.push( f );
        });

        return ret;
    };

    mapctl.prototype.parsewkb = function( data ){
        var ret       = [];
        var rreader   = new FileReader();
        var wktreader = new FileReader();
        var db        = new Blob( [data], { type: 'octet/stream' } );
        var i         = 8;

        wktreader.onload = () => {
            // wkb has structure of??!
            var a = new Int32Array( wktreader.result );
            console.log( a[0].toString(2) );
        }

        rreader.onload = () => { 
            var a = new Uint32Array( rreader.result );

            // char? 8 bytes
            console.log( a );

            wktreader.readAsArrayBuffer( db.slice( i, i + 32 ) );

            // offset?!
            // i += a[1]
            // if( i < db.size ){
            //     rreader.readAsArrayBuffer( db.slice( i, (i+8) ) );
            // }
        };

        rreader.readAsArrayBuffer( db.slice( i-8, i ) );
        // while i -->


        return ret;
    };

    mapctl.prototype.get_data = function( url ) {
        return new Promise( (success, error) => {
            
            var req = new XMLHttpRequest();

            req.onload  = () => success( req );
            req.onerror = () => error( null );

            try{
                req.open( "GET", url, true );
                req.send();
            }
            catch (ex){
                error( ex );
            }

        });
    };

    mapctl.prototype.log = function( message ) {
        console.log( message );
    };

    mapctl.prototype.clear_log = function() {
        console.log( 'clear messages!' );
    };

    window.mapctl = new mapctl();

})()