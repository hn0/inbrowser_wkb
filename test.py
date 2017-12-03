#!/usr/bin/python
"""
    Project simple testing script, written due 

    Author:   Hrvoje Novosel<hrvojedotnovosel@gmail.com>
    Created:  3. Dec 2017
"""

import json
import urllib.request


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

    with urllib.request.urlopen( SERVER +  "/wkb" ) as fp:
        # TODO: create read buffer!!!
        print( fp.read() )

    print( 'WKB test done successfully' )


if __name__ == '__main__':
    try:
        test_server()
        test_wkb()
    except Exception as ex:
        print( 'Test failed!' )
        print( ex )