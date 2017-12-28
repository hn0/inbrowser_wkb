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
#include<byteswap.h>

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
                    // todo: push to double array
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

// for test, same lousy algorithm for comparison sake will be used
void convert( unsigned char* wkb )
{

    // AND ADDITIONAL DOUBLE IS 32 BIT ON COMPILER NOT 64!!!!
    // maybe only  solution is to write bytes back to memory?!
    // see how emscripten handels that
    // printf("x: %g \n", (float)__bswap_64( (unsigned long long)wkb[22] ) ); // expecting 147.288190079...

    double tmp = (double)wkb[22];
    print2( tmp );

    int tst = __bswap_32( 0x806640 ); // YEP, THIS IS IT!!!, BUT WHY HERE THAT ONE BIT HAS NOT BEEN SET!!?!
    printf("%g\n", (float)tst );

    // unsigned long long a = __bswap_64( ((unsigned char*)&wkb)[1] );
    // double b = (double)((int)wkb[1]);
    // double b = (double)(((unsigned char*)&wkb)[22]);
    // return b;
    // int pos = 1;

    // // TODO: again, byte order is not implemented!
    // // TODO: need an array for the coordinates
    // switch( (int)wkb[pos] ){
    //     case 6:
    //         pos += 4;
    //         int n = (int)wkb[pos];
    //         pos += 4;
    //         double ret[n]; // something like this
    //         for( int i=0; i < n; i++ ){
    //             pos += read_geom( wkb, pos );
    //             break;
    //         }
    // }
}