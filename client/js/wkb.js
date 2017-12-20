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
    var dw  = new DataView( wkb );

    // first byte indicates byte order 
    // next 4 bytes denotes geometry type
    var bo  = dw.getUint8( 0, true );

    switch ( dw.getUint32( 1, bo ) ) {
        case 1:
            console.log( 'a wkb_point here' );
            break;
        case 2:
            console.log( 'a wkb_linestring here' );
            break;
        case 3: 
            console.log( 'a wkb_polygon here' );
            break;
        case 4:
        console.log( 'a wkb_multipoint here' );
            break;
        case 5:
            console.log( 'a wkb_multilinestring here' );
            break;
        case 6:
            var n = dw.getUint32( 5, bo );
            // this.coords there is a n polygons here
            this.parse( wkb.slice( 9 ) )
            console.log( 'a wkb_multypolygon here' );
            break;
        case 7:
            console.warn( 'Geometry collection not supported yet!' );
    }

    console.log( 'parsing wkb', bo );
};