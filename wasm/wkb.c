/*

    Simple, POC style script, with wasm platform as compile target

    Created: 22. Dec 2017
    Author:  Hrvoje Novosel<hrvojedotnovosel@gmail.com>
*/

#include<stdio.h>
#include "gdal.h"

char* type( long* wkb )
{

    OGRGeometryH hGeo;

    hGeo = OGR_G_CreateGeometry( wkbMultiPolygon );

    OGR_G_DestroyGeometry( hGeo );

    // now we need a type!
    // printf( "Length of array element is %d, byte order %04x and type %i \n", sizeof( wkb ), wkb[0]&0x1, wkb[0]&0x01 );

    return "EMPTY";
}

char* convert( char* wkb, int len )
{
    // TODO: see if usage of gdal libs is feasible?!
    printf( "Length of array element is %d and complete length of blob is %d\n", sizeof( wkb ), len );

    return "ok";
}
