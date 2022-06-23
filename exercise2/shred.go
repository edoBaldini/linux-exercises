package Shred

import (
	"os"
	"io/fs"
)

/*	Min return the minimum between x and y 	*/
func Min(x, y int64) int64 {
 if x < y {
   return x
 }
 return y
}

/*	Overwrite function overwrites the content of the file with random values generated through urandom device.
	The random values have a size of 1024 and the function overwrites the entire size of the file.
	Input	: file that must be overwritten
	Output	: erorr if occurs or nil
*/
func Overwrite (file *os.File) (error) {
	var index int64
	var size_left int64 = 0
	var b_size int64 = 1024	
	
	f_stat, err := file.Stat()
	if err != nil {
		return err
	}

	for index = 0; index < f_stat.Size(); index += size_left {
                size_left = Min(f_stat.Size() - index, b_size)
		
		b1 := make([]byte, b_size)
		f_rand, _ := os.OpenFile("/dev/random", os.O_RDONLY, 0555)
		f_rand.Read(b1)
		f_rand.Close()
                
		_, err := file.Write(b1[:size_left])
                if err != nil {
                        return err
                }
	}
	return nil
}

/*	The function overwrite the given file through Overwrites function three times and then Shred remove the file.
	
	ATTENTION! 
	In case of symbolic link, Shred overwrites the conted of the file linked and then remove the link and the file.
	Use Shred carefully!
	
	Input	: filename that must be shred
	Output	: error if occurs or nil
*/
func Shred(filename string) error {
	
	fi, err := os.Lstat(filename)
	
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {	
		file, err := os.OpenFile(filename, os.O_WRONLY, 0777)
		if err != nil {
			return err
		}
		
		err2 := Overwrite(file)
		if err2 != nil {
			return err
		}
		
		file.Sync()
		file.Close()
	}

	if (fi.Mode() & fs.ModeSymlink != 0) {
		origin_file, err := os.Readlink(filename)
		if err != nil {  
			return err
		}
		err = os.Remove(origin_file)
		if err != nil {
			return err
		}
	}
	
	err = os.Remove(filename)	
	if err != nil {
		return err
	}
	
	return nil
}
