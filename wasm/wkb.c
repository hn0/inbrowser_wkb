/*

    Simple, POC style script, with wasm platform as compile target

    TODO: currently support for byte order is missing
    for byte swapping refer to:
        https://stackoverflow.com/questions/2182002/convert-big-endian-to-little-endian-in-c-without-using-provided-func

    Created: 22. Dec 2017
    Author:  Hrvoje Novosel<hrvojedotnovosel@gmail.com>
*/

#include<stdio.h>
#include<stdlib.h>

int read_geom( unsigned char* wkb, int pos )
{
    int typ = (int)wkb[++pos];

    switch( typ ){
        case 3:
            pos += 4;
            int n = (int)wkb[pos];
            pos += 4;
            for( int i=0; i < n; i++){
                int ncoord = (int)wkb[pos];
                pos += 4;
                for( int j=0; j < ncoord; j++){
                    // maybe right but byte order is wrong
                    // well mem copy?!
                    double* a = (double*) malloc( sizeof(double) );
                    // *a = (double)wkb[pos];
                    memcpy(a, &wkb[pos], sizeof(double));
                    printf("x: %g \n", *a );
                    pos += 8;
                    break;
                }
                break;
            }
            break;
    }
    return pos;
}

char* type( unsigned char* wkb )
{
    switch( (int)wkb[1] ){
        case 1:
            return "point";
        case 2:
            return "linestring";
        case 3:
            return "polygon";
        case 4:
            return "multipoint";
        case 5:
            return "multilinestring";
        case 6:
            return "multipolygon";
        default:
            printf( "Unsupported geometry type detected, falling back to default empty geometry!\n" );
    }
    return "EMPTY";
}

// for test, same lousy algorithm for comparison sake will be used
void convert( unsigned char* wkb, int len )
{
    int pos = 1;

    // TODO: again, byte order is not implemented!
    // TODO: need an array for the coordinates
    switch( (int)wkb[pos] ){
        case 6:
            pos += 4;
            int n = (int)wkb[pos];
            pos += 4;
            double ret[n]; // something like this
            for( int i=0; i < n; i++ ){
                pos += read_geom( wkb, pos );
                break;
            }
    }
}
