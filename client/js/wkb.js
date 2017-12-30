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
        case 7:
            console.warn( 'Geometry collection not supported yet!' );
        default:
            typ = 'EMPTY'
    }

    this.type = typ;
    var c = this.read( wkb );
    if( typ == 'point' ){ // not complete!
        this.coords = c[0];
    }
};

// read func EXTREAMLY BAD DESIGN, REWRITE IS NEEDED!
geom.prototype.read = function( wkb )
{
    var dw     = new DataView( wkb );
    var coords = [];
    var bo     = dw.getUint8( 0, true );
    var type   = dw.getUint32( 1, bo );
    var read   = 5;

    switch( type ){
        case 1:
            console.log( 'parse wkb pont' );
            break;
        case 2:
            console.log( 'parse linestring' );
            break;
        case 3:
            var n   = dw.getUint32( read, bo );
            read += 4;
            coords  = new Array( n );
            for( var i=0; i < n; i++ ){
                var rlen = dw.getUint32( read, bo ); // number of pts!
                // console.log( 'n rings ', n, rlen, read, 'byte order', bo );
                read += 4;
                coords[i] = new Array( rlen );
                for( var j=0; j < rlen; j++ ){
                    coords[i][j] = [
                        dw.getFloat64( read, bo ),
                        dw.getFloat64( read+8, bo )
                    ];
                    // console.log( read, coords[i][j][0] )
                    read += 16;
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
            var n   = dw.getUint32( 5, bo );
            // console.log( 'number of geoms (js):', n );
            read += 4;
            // console.log( read )
            for( var i=0; i < n; i++){
                var g = this.read( wkb.slice( read ) );
                this.coords.push( g[0] )
                read += g[1];
                // console.log( i, read );
            }
            break;
        default:
            console.log( 'unknown type ', read, type);
    }
    return [coords, read];
};