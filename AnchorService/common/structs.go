package common

import "github.com/compasses/GOProjects/AnchorService/util"

var Log = util.CommonLooger


//the first 2 bytes are ASCII Fa
//the next 6 bytes are the directory block height.
//the next 32 bytes are the KeyMR  of the directory block at that height


type DirectoryBlockAnchorInfo struct {
    KeyMR       *Hash
    DBHeight    uint32
}