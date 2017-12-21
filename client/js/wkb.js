/*
    Web parsing class, implemented based on wkb definition form: http://edndoc.esri.com/arcsde/9.1/general_topics/wkb_representation.htm

    Created: 20. Dec 2017
    Author:  Hrvoje Novosel<hrvojedotnovosel@gmail.com>

*/

function wkb_format() 
{
    return new geom();
};


var geom = function(){
    this.type   = 'EMPTY';
    this.coords = [];
};

geom.prototype.parse = function( wkb )
{
    this.coords.length = 0;
    var dw  = new DataView( wkb );
    // first byte indicates byte order 
    // next 4 bytes denotes geometry type
    var bo  = dw.getUint8( 0, true );
    var typ = '';

    switch ( dw.getUint32( 1, bo ) ) {
        case 1:
            typ = 'point';
            break;
        case 2:
            typ = 'linestring';
            break;
        case 3: 
            typ = 'polygon';
            break;
        case 4:
            typ = 'multipoint';
            break;
        case 5:
            typ = 'multilinestring';
            break;
        case 6:
            typ = 'multipolygon';
            break;
            // var n = dw.getUint32( 5, bo );
            // this.coords there is a n polygons here
            // this.parse( wkb.slice( 9 ) )
            // console.log( 'a wkb_multypolygon here' );
        case 7:
            console.warn( 'Geometry collection not supported yet!' );
        default:
            typ = 'EMPTY'
    }

    this.type = typ;
    var c = this.read( wkb );
    if( typ == 'point' ){ // not complete!
        this.coords = c;
    }
};

geom.prototype.read = function( wkb )
{
    var dw     = new DataView( wkb );
    var coords = [];
    var bo     = dw.getUint8( 0, true );
    var type   = dw.getUint32( 1, bo );

    switch( type ){
        case 1:
            console.log( 'parse wkb pont' );
            break;
        case 2:
            console.log( 'parse linestring' );
            break;
        case 3:
            var n   = dw.getUint32( 5, bo );
            var pos = 9;
            coords  = new Array( n );
            for( var i=0; i < n; i++ ){
                var rlen = dw.getUint32( pos, bo ); // number of pts!
                pos += 4;
                coords[i] = new Array( rlen );
                for( var j=0; j < rlen; j++ ){
                    coords[i][j] = [
                        dw.getFloat64( pos, bo ),
                        dw.getFloat64( pos+8, bo )
                    ];
                    pos += 16;
                }
            }
            break;
        case 4:
            console.log( 'parse multipoint' );
            break;
        case 5:
            console.log( 'parse multilinestring' );
            break;
        case 6:
            var n = dw.getUint32( 5, bo );
            for( var i=0; i < n; i++){
                // now when first part is done, where is the pointer to the next geometry!?
                // slight rewrite of this idea is needed!
                this.coords.push( this.read( wkb.slice( 9 ) ) )
                break;
            }
            break;
    }
    return coords;
};