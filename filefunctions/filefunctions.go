package filefunctions

import (
    "bytes"
    "encoding/binary"
    "encoding/hex"
    "log"
    "os"
    "errors"
    "fmt"
)

//ByteCount tracks the number of bytes needed to parse the entire block. Often times less than BlockLength
var ByteCount int

// ErrDetailedMagic is thrown when the detailed search for magic number searches for too long
var ErrDetailedMagic = errors.New("DetailedLookForMagic: could not find magic number")

//Possible64ByteErrorFlag tracks whether we've hit a missed byte in parsing in the main method output count
//var Possible64ByteErrorFlag bool

//ReadNextBytes reads number of bytes from file
func ReadNextBytes(file *os.File, number int) ([]byte, error) {
  bytes := make([]byte, number)
  ByteCount = ByteCount + number

  _, err := file.Read(bytes)
  if err != nil {
    return nil, err
  }
  return bytes, nil
}

//RewindAndRead64 is called when Possible64ByteErrorFlag is raised. The function moves the file pointer back and re-reads with fewer bytes included in the read
func RewindAndRead64(b []byte, file *os.File, outputValue *uint64) ([]byte, error) {
  var secondTryLen int64 = 7
    bytesTwo := make([]byte, secondTryLen)
    ByteCount = ByteCount - int(secondTryLen) - 1
    _, _ = file.Seek(-(secondTryLen + 1), 1)
    _, err := file.Read(bytesTwo)
    if err != nil {
      return nil, err
    }
    bytesTwo = append(bytesTwo[:], []byte{0} ...)
    ByteCount = ByteCount + int(secondTryLen) + 1
    ReadBinaryToUInt64(bytesTwo, outputValue)
    return bytesTwo, nil
}

//RewindAndRead32 is called when we fail validation of unsigned 32 bit integer and want to skip back a bit and restart parsing
func RewindAndRead32(b []byte, file *os.File, transactionIndex *uint32) ([]byte, error) {
  var secondTryLen int64 = 4
    bytesTwo := make([]byte, secondTryLen)
    ByteCount = ByteCount - int(secondTryLen) - 1

    _, _ = file.Seek(-(secondTryLen + 1), 1)
    _, err := file.Read(bytesTwo)
    if err != nil {
      return nil, err
    }
    ByteCount = ByteCount + int(secondTryLen) + 1
    ReadBinaryToUInt32(bytesTwo, transactionIndex)
    return bytesTwo, nil
  }

//StepBack sets the file pointer back length and updates the ByteCount field
func StepBack(length int, file *os.File) {
  _,_ = file.Seek(-int64(length), 1)
  ByteCount = ByteCount - length
}

//LookForMagic handles instance when encounter string of zeros in searching for Magic Number
func LookForMagic(file *os.File) (uint32, error) {
  var iter uint32
  for iter != 0xD9B4BEF9 {
    b, err := ReadNextBytes(file, 4)
    if err != nil {
      return 0, err
    }
    err = ReadBinaryToUInt32(b, &iter)
    if err != nil {
      log.Fatal("Read binary in LookForMagic failed: ", err)
    }
  }
  fmt.Println(iter)
  return iter, nil
}

//DetailedLookForMagic goes byte-by-byte to look for magic number
func DetailedLookForMagic(file *os.File) (uint32, error) {
  var iter uint32
  var track int
  for iter != 0xD9B4BEF9 {

    b, err := ReadNextBytes(file, 4)
    if err != nil {
      return 0, err
    }
    
    err = ReadBinaryToUInt32(b, &iter)
    if err != nil {
      log.Fatal("Read binary in DetailedLookForMagic failed: ", err)
    }

    track++

    if track > 50000 {
      return 0, ErrDetailedMagic
    }
  }
  return iter, nil
}

//ReadUInt8ByteArray reads a bytestream into an array of their values
func ReadUInt8ByteArray(b []byte, passedVariable *[]uint8) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, passedVariable)
  return err
}
// ReadUInt8ByteArrayToString reads a bytestream into a string of arbitrary length
func ReadUInt8ByteArrayToString(b []byte, passedVariable *string) {
  *passedVariable = hex.EncodeToString(b)
}

//ReadUInt8ByteArrayLength32 reads a bytestream into a string hash with 32 characters
func ReadUInt8ByteArrayLength32(b []byte, passedVariable *string) {
  *passedVariable = hex.EncodeToString(b)
}

//ReadBinaryToUInt8 reads a binary bytestream into an unsigned integer byte
func ReadBinaryToUInt8(b []byte, passedVariable *uint8) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, passedVariable)
  return err
}

//ReadBinaryToUInt16 reads a binary bytestream into an unsigned integer of 2 bytes
func ReadBinaryToUInt16(b []byte, passedVariable *uint16) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, passedVariable)
  return err
}

//ReadBinaryToUInt16Big reads a binary bytestream in big endian into an unsigned integer of 2 bytes
func ReadBinaryToUInt16Big(b []byte, passedVariable *uint16) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.BigEndian, passedVariable)
  return err
}

//ReadBinaryToUInt32 reads a binary bytestream into an unsigned integer of 4 bytes
func ReadBinaryToUInt32(b []byte, passedVariable *uint32) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, passedVariable)
  return err
}

//ReadBinaryToUInt32Big reads a binary bytestream encoded in Big Endian into an unsigned integer of 4 bytes
func ReadBinaryToUInt32Big(b []byte, passedVariable *uint32) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.BigEndian, passedVariable)
  return err
}

//ReadBinaryToUInt64 reads a binary bytestream into an unsigned integer of 8 bytes
func ReadBinaryToUInt64(b []byte, passedVariable *uint64) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.LittleEndian, passedVariable)
  return err
}

//ReadBinaryToUInt64Big reads a binary bytestream in big endian into an unsigned integer of 8 bytes
func ReadBinaryToUInt64Big(b []byte, passedVariable *uint64) (error) {
  buf := bytes.NewReader(b)
  err := binary.Read(buf, binary.BigEndian, passedVariable)
  return err
}



// ReadVariableLengthIntegerWithBackup reads a variable length integer as described by the bitcoin protocol into an unsigned 8 byte integer
// // and sets the file pointer back after each read
func ReadVariableLengthIntegerWithBackup(file *os.File) (uint64, []byte, error) {

  var ret uint64
  var eight uint8
  var byteret []byte

  bytes := make([]byte, 1)
  _, err := file.Read(bytes)
  if err != nil {
    return ret, nil, err
  }
  err = ReadBinaryToUInt8(bytes, &eight)
  if err != nil {
    return ret, nil, err
  }
  if eight < 0xFD {       // If it's less than 0xFD use this value as the unsigned integer
    ByteCount++
    ret = uint64(eight)
    byteret = bytes
  } else {
      StepBack(1, file)
      var sixteen uint16
      bytes = make([]byte, 3)
      _, err = file.Read(bytes)
      if err != nil {
        return ret, nil, err
      }
      err = ReadBinaryToUInt16(bytes, &sixteen)
      if err != nil {
        return ret, nil, err
      }
      if sixteen < 0xFFFF {
        ByteCount += 3
        ret = uint64(sixteen)
        byteret = bytes
      } else {
          StepBack(3, file)
          var thirtytwo uint32
          bytes = make([]byte, 5)
          _, err = file.Read(bytes)
          if err != nil {
            return ret, nil, err
          }
          err = ReadBinaryToUInt32(bytes, &thirtytwo)
          if err != nil {
            return ret, nil, err
          }
          if thirtytwo < 0xFFFFFFFF {
            ByteCount += 5
            ret = uint64(thirtytwo)
            byteret = bytes
          } else {  // never expect to actually encounter a 64bit integer in the block-chain stream; it's outside of any reasonable expected value
              StepBack(5, file)
              var sixtyfour uint64
              bytes = make([]byte, 9)
              _, err = file.Read(bytes)
              if err != nil {
                return ret, nil, err
              }
              err = ReadBinaryToUInt64(bytes, &sixtyfour)
              if err != nil {
                return ret, nil, err
              }
              ByteCount += 9
              ret = uint64(sixtyfour)
              byteret = bytes
            }
          }
      }
  return ret, byteret, nil
}

// ReadVariableLengthInteger reads a variable length integer as described by the bitcoin protocol into an unsigned 8 byte integer
func ReadVariableLengthInteger(file *os.File) (uint64, []byte, error) {

  var ret uint64
  var eight uint8
  var byteret []byte

  bytes := make([]byte, 1)
  _, err := file.Read(bytes)
  if err != nil {
    return ret, nil, err
  }
  err = ReadBinaryToUInt8(bytes, &eight)
  if err != nil {
    return ret, nil, err
  }
  if eight < 0xFD {       // If it's less than 0xFD use this value as the unsigned integer
    ByteCount++
    ret = uint64(eight)
    byteret = bytes
  } else {
      var sixteen uint16
      bytes = make([]byte, 3)
      _, err = file.Read(bytes)
      if err != nil {
        return ret, nil, err
      }
      err = ReadBinaryToUInt16(bytes, &sixteen)
      if err != nil {
        return ret, nil, err
      }
      if sixteen < 0xFFFF {
        ByteCount += 3
        ret = uint64(sixteen)
        byteret = bytes
      } else {
          var thirtytwo uint32
          bytes = make([]byte, 5)
          _, err = file.Read(bytes)
          if err != nil {
            return ret, nil, err
          }
          err = ReadBinaryToUInt32(bytes, &thirtytwo)
          if err != nil {
            return ret, nil, err
          }
          if thirtytwo < 0xFFFFFFFF {
            ByteCount += 5
            ret = uint64(thirtytwo)
            byteret = bytes
          } else {      // never expect to actually encounter a 64bit integer in the block-chain stream; it's outside of any reasonable expected value
              var sixtyfour uint64
              bytes = make([]byte, 9)
              _, err = file.Read(bytes)
              if err != nil {
                return ret, nil, err
              }
              err = ReadBinaryToUInt64(bytes, &sixtyfour)
              if err != nil {
                return ret, nil, err
              }
              ByteCount += 9
              ret = uint64(sixtyfour)
              byteret = bytes
            }
          }
      }
  return ret, byteret, nil
}

// ReadVarIntFromBytes reads a variable length integer as described by the bitcoin protocol into an unsigned 8 byte integer
func ReadVarIntFromBytes(b []byte, start int) (uint64, int, error) {

  var ret uint64
  var eight uint8
  index := start


  err := ReadBinaryToUInt8(b[start:start+1], &eight)
  if err != nil {
    return ret, start+1, err
  }
  if eight < 0xFD {       // If it's less than 0xFD use this value as the unsigned integer
    ByteCount++
    ret = uint64(eight)
  } else {
      index = start + 1
      var sixteen uint16
      err = ReadBinaryToUInt16(b[index:index+2], &sixteen)
      if err != nil {
        return ret, start+2, err
      }
      if sixteen < 0xFFFF {
        ByteCount += 3
        ret = uint64(sixteen)
        index = start + 2
      } else {
          var thirtytwo uint32
          err = ReadBinaryToUInt32(b[index:index+4], &thirtytwo)
          if err != nil {
            return ret, start+4, err
          }
          if thirtytwo < 0xFFFFFFFF {
            ByteCount += 5
            ret = uint64(thirtytwo)
            index = start + 4
          } else {      // never expect to actually encounter a 64bit integer in the block-chain stream; it's outside of any reasonable expected value
              var sixtyfour uint64
              err = ReadBinaryToUInt64(b[index:index+8], &sixtyfour)
              if err != nil {
                return ret, start+8, err
              }
              ByteCount += 9
              ret = uint64(sixtyfour)
              index = start + 8
            }
          }
      }
  return ret, index + 1, nil
}

//ResetBlockHeadPointer points the byte-reader to the next block in the chain
func ResetBlockHeadPointer(blockLength uint32, file *os.File) ([]byte, error) {
  if ByteCount <= int(blockLength) {
    bytes := make([]byte, int(blockLength) - ByteCount)
    _, err := file.Read(bytes)
    if err != nil {
        return nil, err
    }
    ByteCount = int(blockLength)
    return bytes, nil
  }
  return nil, errors.New("used more bytes than listed in blocklength")
}

//GetByteCount returns the global ByteCount variable in filefunctions class
func GetByteCount() (int) {
  return ByteCount
}

//SetByteCount sets the global ByteCount variable in filefunctions class
func SetByteCount(newVal int) {
  ByteCount = newVal
}

//IncrementByteCount increments by incrementVal the global ByteCount variable in filefunctions class
func IncrementByteCount(incrementVal int) {
  ByteCount += incrementVal
}
