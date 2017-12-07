(function() {
    

    mapctl   = function(){
        this.server   = 'http://localhost:8000/';
        this.requests = { 'wkb': null, 'wkt': null };
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
                        // var features = this['parse' + type]( xhr.response );
                        console.log( xhr.response )

                        this.requests[type] = null;

                    })
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

    mapctl.prototype.init_map = function( container ){
        var map = new ol.Map({
            target: container,
            layers: [ new ol.layer.Tile({ source: new ol.source.OSM() }) ]
        });

        return map;
    };

    mapctl.prototype.parsewkt = function( data ){
        var ret = [];

        console.log( data )

        return ret;
    };

    mapctl.prototype.parsewkb = function( data ){
        var ret = [];

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