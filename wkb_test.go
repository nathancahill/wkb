
package wkb

import (
    "os"
    "testing"
)

func TestWkb(t *testing.T) {
    file, err := os.Open("wkb.bin")

    if err != nil {
        panic(err)
    }

    out, err := ReadWKBRaster(file)

    if err != nil {
        panic(err)
    }

    if out.Version != 0 {
        t.Errorf("Expected version to be 0, got %d instead.", out.Version)
    }

    if out.ScaleX != 12.041454266862495 {
        t.Errorf("Expected scaleX to be 12.041454266862495, got %d instead.", out.ScaleX)
    }

    if out.ScaleY != -12.041454266862495 {
        t.Errorf("Expected scaleY to be -12.041454266862495, got %d instead.", out.ScaleY)
    }

    if out.IpX != -1.1799927868252808e+07 {
        t.Errorf("Expected ipX to be -1.1799927868252808e+07, got %d instead.", out.IpX)
    }

    if out.IpY != 5.012423608334331e+06 {
        t.Errorf("Expected ipY to be 5.012423608334331e+06, got %d instead.", out.IpY)
    }

    if out.SkewX != 0 {
        t.Errorf("Expected skewX to be 0, got %d instead.", out.SkewX)
    }

    if out.SkewY != 0 {
        t.Errorf("Expected skewY to be 0, got %d instead.", out.SkewY)
    }

    if out.Srid != 3857 {
        t.Errorf("Expected srid to be 3857, got %d instead.", out.Srid)
    }

    if out.Width != 100 {
        t.Errorf("Expected width to be 100, got %d instead.", out.Width)
    }

    if out.Height != 100 {
        t.Errorf("Expected height to be 100, got %d instead.", out.Height)
    }

    if len(out.Bands) != 1 {
        t.Errorf("Expected number of bands to be 1, got %d instead.", len(out.Bands))
    }

    if out.Bands[0].IsOffline {
        t.Errorf("Expected band 1 isOffline to be false, got %v instead", out.Bands[0].IsOffline)
    }

    if !out.Bands[0].HasNoDataValue {
        t.Errorf("Expected band 1 hasNoDataValue to be true, got %v instead", out.Bands[0].HasNoDataValue)
    }

    if out.Bands[0].IsNoDataValue {
        t.Errorf("Expected band 1 isNoDataValue to be false, got %v instead", out.Bands[0].IsNoDataValue)
    }

    if out.Bands[0].NoData != -9223372036854775808 {
        t.Errorf("Expected band 1 nodata to be -9223372036854775808, got %d instead", out.Bands[0].NoData)
    }

    if len(out.Bands[0].Data) != 100 {
        t.Errorf("Expected band 1 rows to be 100, got %d instead.", len(out.Bands[0].Data))
    }

    if len(out.Bands[0].Data[0]) != 100 {
        t.Errorf("Expected band 1 columns to be 100, got %d instead.", len(out.Bands[0].Data))
    }

    if out.Bands[0].Data[0][0] != 2402 {
        t.Errorf("Expected band 1 0x0 to be 2402, got %d instead.", out.Bands[0].Data[0][0])
    }

    if out.Bands[0].Data[50][50] != 2406 {
        t.Errorf("Expected band 1 0x0 to be 2406, got %d instead.", out.Bands[0].Data[50][50])
    }

    if out.Bands[0].Data[99][99] != 2411 {
        t.Errorf("Expected band 1 0x0 to be 2411, got %d instead.", out.Bands[0].Data[99][99])
    }
}
