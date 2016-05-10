
package wkb

import (
    "io"
    "errors"
    "encoding/binary"
)

type Header struct {
    Version     uint16
    Bands       uint16
    ScaleX      float64
    ScaleY      float64
    IpX         float64
    IpY         float64
    SkewX       float64
    SkewY       float64
    Srid        int32
    Width       uint16
    Height      uint16
}

type RasterBand struct {
    NoData          int
    IsOffline       bool
    HasNoDataValue  bool
    IsNoDataValue   bool
    Data            [][]int
}

type Raster struct {
    Version     uint16
    ScaleX      float64
    ScaleY      float64
    IpX         float64
    IpY         float64
    SkewX       float64
    SkewY       float64
    Srid        int32
    Width       uint16
    Height      uint16
    Bands       []RasterBand
}

func readIntOfType(reader io.Reader, endiannes binary.ByteOrder, valueType int) (int, error) {
    switch valueType {
    case 0, 1, 2, 4:
        var value uint8

        err := binary.Read(reader, endiannes, &value)

        if err != nil {
            return 0, err
        }

        return int(value), nil
    case 3:
        var value int8

        err := binary.Read(reader, endiannes, &value)

        if err != nil {
            return 0, err
        }

        return int(value), nil
    case 5:
        var value int16

        err := binary.Read(reader, endiannes, &value)

        if err != nil {
            return 0, err
        }

        return int(value), nil
    case 6:
        var value uint16

        err := binary.Read(reader, endiannes, &value)

        if err != nil {
            return 0, err
        }

        return int(value), nil
    case 7:
        var value int32

        err := binary.Read(reader, endiannes, &value)

        if err != nil {
            return 0, err
        }

        return int(value), nil
    case 8:
        var value uint32

        err := binary.Read(reader, endiannes, &value)

        if err != nil {
            return 0, err
        }

        return int(value), nil
    case 9:
        var value float32

        err := binary.Read(reader, endiannes, &value)

        if err != nil {
            return 0, err
        }

        return int(value), nil
    case 10:
        var value float64

        err := binary.Read(reader, endiannes, &value)

        if err != nil {
            return 0, err
        }

        return int(value), nil
    }

    return 0, errors.New("Unknown value type")
}

func ReadWKBRaster(wkb io.Reader) (Raster, error) {
    raster := Raster{}

    // Determine the endiannes of the raster
    //
    // +---------------+-------------+------------------------------+
    // | endiannes     | byte        | 1:ndr/little endian          |
    // |               |             | 0:xdr/big endian             |
    // +---------------+-------------+------------------------------+

    var endiannesValue uint8
    err := binary.Read(wkb, binary.LittleEndian, &endiannesValue)

    if err != nil {
        return raster, err
    }

    var endiannes binary.ByteOrder

    if endiannesValue == 0 {
        endiannes = binary.BigEndian
    } else if endiannesValue == 1 {
        endiannes = binary.LittleEndian
    }

    // Read the raster header data.
    //
    // +---------------+-------------+------------------------------+
    // | version       | uint16      | format version (0 for this   |
    // |               |             | structure)                   |
    // +---------------+-------------+------------------------------+
    // | nBands        | uint16      | Number of bands              |
    // +---------------+-------------+------------------------------+
    // | scaleX        | float64     | pixel width                  |
    // |               |             | in geographical units        |
    // +---------------+-------------+------------------------------+
    // | scaleY        | float64     | pixel height                 |
    // |               |             | in geographical units        |
    // +---------------+-------------+------------------------------+
    // | ipX           | float64     | X ordinate of upper-left     |
    // |               |             | pixel's upper-left corner    |
    // |               |             | in geographical units        |
    // +---------------+-------------+------------------------------+
    // | ipY           | float64     | Y ordinate of upper-left     |
    // |               |             | pixel's upper-left corner    |
    // |               |             | in geographical units        |
    // +---------------+-------------+------------------------------+
    // | skewX         | float64     | rotation about Y-axis        |
    // +---------------+-------------+------------------------------+
    // | skewY         | float64     | rotation about X-axis        |
    // +---------------+-------------+------------------------------+
    // | srid          | int32       | Spatial reference id         |
    // +---------------+-------------+------------------------------+
    // | width         | uint16      | number of pixel columns      |
    // +---------------+-------------+------------------------------+
    // | height        | uint16      | number of pixel rows         |
    // +---------------+-------------+------------------------------+

    var header Header
    err = binary.Read(wkb, endiannes, &header)

    if err != nil {
        return raster, err
    }

    raster.Version = header.Version
    raster.ScaleX = header.ScaleX
    raster.ScaleY = header.ScaleY
    raster.IpX = header.IpX
    raster.IpY = header.IpY
    raster.SkewX = header.SkewX
    raster.SkewY = header.SkewY
    raster.Srid = header.Srid
    raster.Width = header.Width
    raster.Height = header.Height

    for i := 0; i < int(header.Bands); i++ {
        band := RasterBand{}

        // Read band header data
        //
        // +---------------+--------------+-----------------------------------+
        // | isOffline     | 1bit         | If true, data is to be found      |
        // |               |              | on the filesystem, trought the    |
        // |               |              | path specified in RASTERDATA      |
        // +---------------+--------------+-----------------------------------+
        // | hasNodataValue| 1bit         | If true, stored nodata value is   |
        // |               |              | a true nodata value. Otherwise    |
        // |               |              | the value stored as a nodata      |
        // |               |              | value should be ignored.          |
        // +---------------+--------------+-----------------------------------+
        // | isNodataValue | 1bit         | If true, all the values of the    |
        // |               |              | band are expected to be nodata    |
        // |               |              | values. This is a dirty flag.     |
        // |               |              | To set the flag to its real value |
        // |               |              | the function st_bandisnodata must |
        // |               |              | must be called for the band with  |
        // |               |              | 'TRUE' as last argument.          |
        // +---------------+--------------+-----------------------------------+
        // | reserved      | 1bit         | unused in this version            |
        // +---------------+--------------+-----------------------------------+
        // | pixtype       | 4bits        | 0: 1-bit boolean                  |
        // |               |              | 1: 2-bit unsigned integer         |
        // |               |              | 2: 4-bit unsigned integer         |
        // |               |              | 3: 8-bit signed integer           |
        // |               |              | 4: 8-bit unsigned integer         |
        // |               |              | 5: 16-bit signed integer          |
        // |               |              | 6: 16-bit unsigned signed integer |
        // |               |              | 7: 32-bit signed integer          |
        // |               |              | 8: 32-bit unsigned signed integer |
        // |               |              | 9: 32-bit float                   |
        // |               |              | 10: 64-bit float                  |
        // +---------------+--------------+-----------------------------------+
        //
        // Requires reading a single byte, and splitting the bits into the
        // header attributes
        var bandheader uint8
        err := binary.Read(wkb, endiannes, &bandheader)

        if err != nil {
            return raster, err
        }

        band.IsOffline = (int(bandheader) & 128) != 0
        band.HasNoDataValue = (int(bandheader) & 64) != 0
        band.IsNoDataValue = (int(bandheader) & 32) != 0

        // Read the pixel type
        pixType := (int(bandheader) & 15) - 1

        // +---------------+--------------+-----------------------------------+
        // | nodata        | 1 to 8 bytes | Nodata value                      |
        // |               | depending on |                                   |
        // |               | pixtype [1]  |                                   |
        // +---------------+--------------+-----------------------------------+

        // Read the nodata value
        noData, err := readIntOfType(wkb, endiannes, pixType)

        if err != nil {
            return raster, err
        }

        band.NoData = noData

        // Read the pixel values: width * height * size
        //
        // +---------------+--------------+-----------------------------------+
        // | pix[w*h]      | 1 to 8 bytes | Pixels values, row after row,     |
        // |               | depending on | so pix[0] is upper-left, pix[w-1] |
        // |               | pixtype [1]  | is upper-right.                   |
        // |               |              |                                   |
        // |               |              | As for endiannes, it is specified |
        // |               |              | at the start of WKB, and implicit |
        // |               |              | up to 8bits (bit-order is most    |
        // |               |              | significant first)                |
        // |               |              |                                   |
        // +---------------+--------------+-----------------------------------+

        for i := 0; i < int(header.Height); i++ {
            row := []int{}

            for i := 0; i < int(header.Width); i++ {
                value, err := readIntOfType(wkb, endiannes, pixType)

                if err != nil {
                    return raster, err
                }

                row = append(row, value)
            }

            band.Data = append(band.Data, row)
        }

        raster.Bands = append(raster.Bands, band)
    }

    return raster, nil
}
