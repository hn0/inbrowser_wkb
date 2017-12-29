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
#include<string.h>

void print2(double x)//print double x in binary
{
    union {
        double x;
        char c[sizeof(double)];
    } u;

    u.x = x;

    for (unsigned ofs = 0; ofs < sizeof(double); ofs++) {
        for(int i = 7; i >= 0; i--) {
            printf(((1 << i) & u.c[ofs]) ? "1" : "0");
        }
    }
    printf("\n");
}

int read_geom( unsigned char* wkb, int pos )
{
    int typ = (int)wkb[++pos];

    switch( typ ){
        case 3:
            pos += 4;
            int n = (int)wkb[pos];
            pos += 4;
            double* a = malloc( sizeof( double ) );
            double* b = malloc( sizeof( double ) );
            for( int i=0; i < n; i++){
                int ncoord = (int)wkb[pos];
                pos += 4;
                for( int j=0; j < ncoord; j++){
                    // todo: push to double array
                    memcpy( a, &wkb[pos], sizeof( double ) );
                    memcpy( b, &wkb[pos+8], sizeof( double ) );
                    // printf("x: %g   y: %g\n", *a, *b);
                    pos += 16;
                }
                printf("ring done, next type:%i \n", (long)wkb[pos+1] );
            }
            free(a);
            free(b);
            break;
        default:
            printf( "type %i and pos %i\n", typ, pos );
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
void convert( unsigned char* wkb )
{

    // double* a = (double*)malloc( sizeof(double) );
    // memcpy( a, &wkb[22], sizeof( double ) );
    // print2( *a );
    // printf("value: %g\n", *a );
    // free(a);

    // TODO: again, byte order is not implemented!
    // TODO: need an array for the coordinates
    int pos = 1;
    switch( (int)wkb[pos] ){
        case 6:
            pos += 4;
            int n = (int)wkb[pos];
            pos += 4;
            printf( "n poly %i\n", n);
            for( int i=0; i < n; i++ ){
                pos = read_geom( wkb, pos );
                printf( "%i -> poly done pos: %i\n", i, pos );
            }
    }
}