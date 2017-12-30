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

typedef unsigned long long uint;
// actually we only need a double**!
typedef struct Polygons {
    double** rings;
} Poly;

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

// need a pointer for the doubles and array len!?
uint read_geom( unsigned char* wkb, uint pos, Poly* poly )
{
    uint read = 1;
    int typ  = (int)wkb[pos+read];
    switch( typ ){
        case 3:
            read += 4;
            uint n = (uint)wkb[pos+read] | ( (uint)wkb[pos+read+1] << 8) | ( (uint)wkb[pos+read+2] << 16) | ( (uint)wkb[pos+read+3] << 24); // n rings
            read += 4;
            poly->rings = malloc( sizeof( double* ) * n );
            for( int i=0; i < n; i++ ){
                uint ncoord = (uint)wkb[pos+read] | ( (uint)wkb[pos+read+1] << 8) | ( (uint)wkb[pos+read+2] << 16) | ( (uint)wkb[pos+read+3] << 24); // n coords
                read += 4;
                double* coord = malloc( sizeof( double ) * ncoord * 2 );
                poly->rings[i] = coord;
                int cpos = -1;
                for( int j=0; j < ncoord; j++ ){
                    // todo: push to double array
                    memcpy( &coord[++cpos], &wkb[pos+read], sizeof( double ) );
                    memcpy( &coord[++cpos], &wkb[pos+read+8], sizeof( double ) );
                    // printf("x: %g   y: %g\n", coord[cpos-1], coord[cpos]);
                    read += 16;
                }
            }
            break;
        default:
            printf( "type %i and pos %llu\n", typ, pos );
            break;
    }

    return read;
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
            uint n = (int)wkb[pos];
            pos += 4;
            // printf( "n poly %llu\n", n);
            Poly* polygons = malloc( n * sizeof *polygons );
            for( int i=0; i < n; i++ ){
                // what is happening on 6 polygon, js reads the same value!?
                pos += read_geom( wkb, pos, &polygons[i] );
                // printf( "%i -> poly done pos: %d \n", i, pos );
                // TODO: now to return the value, see how to do it over emscripten
            }
    }
}