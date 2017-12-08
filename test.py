#!/usr/bin/python
"""
    Project simple testing script, written due 

    Author:   Hrvoje Novosel<hrvojedotnovosel@gmail.com>
    Created:  3. Dec 2017
"""

import io
import time
import json
import struct
import urllib.request
from osgeo import ogr

SERVER = 'http://localhost:8000'


def test_server():
    """
        Method for testing main rest point of the server
    """
    print( 'Testing server url' )
    with urllib.request.urlopen( SERVER ) as fp:
        res = json.load( fp )
        if not type( res ) is list:
            raise Exception( 'Unexpected type, got:', str( type( res) ) )
        print( 'Server says: ' + res[0]['message'] )
    print( 'Server test done successfully' )


def test_wkb():
    """
        Tests wkb geometry response
    """
    print( 'Testing wkb geometry type' )
    start     = time.time()
    geocnt    = 0
    recordcnt = 0
    with urllib.request.urlopen( SERVER +  "/wkb" ) as fp:
        # TODO: create read buffer!!!
        ba = io.BytesIO( fp.read() )

        while 1:
            id = struct.unpack( 'i', ba.read( 4 ) )[0]
            n  = struct.unpack( 'i', ba.read( 4 ) )[0]
            if id * n == 0:
                break
            recordcnt += 1
            geo = ogr.CreateGeometryFromWkb( ba.read( n ) )
            if geo.IsValid() and not geo.IsEmpty():
                geocnt += 1
            # print( (id, n) )

    print( 'Got {} records from which {} geometries is valid'.format( recordcnt, geocnt ) )
    print( 'WKB test done successfully, exec time: {}'.format( round( time.time() - start, 5) * 10000 ) )


def test_wkt():
    """
        Tests wkt request
    """
    print( 'Testing wkt geometry type' )
    start = time.time()
    with urllib.request.urlopen( SERVER + "/wkt" ) as fp:
        res = json.load( fp )
        if not type( res ) is list:
            raise Exception( 'Unexpected type, got:', str( type( res) ) )
    
    geocnt    = 0
    recordcnt = 0
    for r in res:
        recordcnt += 1
        geo = ogr.CreateGeometryFromWkt( r['WKT'] )
        if geo.IsValid() and not geo.IsEmpty():
            geocnt += 1

    print( 'Got {} records from which {} geometries is valid'.format( recordcnt, geocnt ) )
    print( 'WKT test done successfully, exec time: {}'.format( round( time.time() - start, 5) * 10000 ) )


def test_info():
    """
        A simple test that will validate correct ( at least expected any ) srs definition!
    """
    print( 'Testing info on geometry' )
    with urllib.request.urlopen( SERVER + "/geo" ) as fp:
        res = json.load( fp )
        print( res[0] )

    print( 'Geo info request done!' )

if __name__ == '__main__':
    try:
        test_server()
        test_info()
        test_wkb()
        test_wkt()
        # TODO: test metadata
    except Exception as ex:
        print( 'Test failed!' )
        print( ex )