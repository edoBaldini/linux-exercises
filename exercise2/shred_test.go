package Shred

import (
	"os"
	"testing"
	"fmt"
)

/*	
	This test checks if an empty file is correctly removed by Shred 
	Expected result: nil
*/
func TestEmptyFile(t *testing.T) {
	file, err := os.Create("empty")
	if err != nil {
		fmt.Printf("error during the environment setup %s", err)
	}
	file.Close()
	answ := Shred("empty")
	if answ != nil {
		t.Errorf("Nil expected, instead %s", answ)
	}
}

/*	
	This test checks if a non empty file is correctly removed by Shred 
	Expected result: nil
*/
func TestNonEmptyFile(t *testing.T) {
	file, err := os.Create("no_empty")
	if err != nil {
		fmt.Printf("error during the environment setup %s", err)
	}       
	_, err = file.Write([]byte("hello"))
	if err != nil {
		fmt.Printf("error during the environment setup %s", err)
	}
	file.Close()
	answ := Shred("no_empty")
	if answ != nil {
		t.Errorf("Nil expected, instead %s", answ)
	}
}

/*	
	This test checks if a symboluc link to a file is correctly removed by Shred 
	Expected result: nil
*/
func TestSymLink(t *testing.T) {
	file, err := os.Create("no_empty")
	if err != nil {
		fmt.Printf("error during the environment setup %s", err)
	}       
	_, err = file.Write([]byte("hello"))
	if err != nil {
		fmt.Printf("error during the environment setup %s", err)
	}
	file.Close()
	
	err = os.Symlink("no_empty", "symlink-to-file")
	if err != nil {
		fmt.Printf("error during the environment setup %s", err)
	}
	answ := Shred("symlink-to-file")
	if answ != nil {
		t.Errorf("Nil expected, instead %s", answ)
	}
}

/*	
	This test checks if Shred fails providing a file that does NOT exist 
	Expected result: error
*/
func TestNonExistingFile(t *testing.T) {
	answ := Shred("no_empty")
	if answ == nil {
		t.Errorf("Error expected, instead %s", answ)
	}
}

/*	
	This test checks if Shred fails providing a file without write privilege 
	Expected result: error
*/
func TestNoWriteFile(t *testing.T) {
	file, err := os.Create("no_write")
	if err != nil {
		fmt.Printf("error during the environment setup %s", err)
	}       
	file.Close()
	if err != nil {
		fmt.Printf("error during the environment setup %s", err)
	}       
	os.Chmod("no_write", 0555)
	answ := Shred("no_write")
	if answ == nil {
		t.Errorf("Error expected, instead %s", answ)
	}
	os.Remove("no_write")
}

/*	
	This test checks if Shred fails providing a file without rw privilege 
	Expected result: error
*/
func TestNoRWFile(t *testing.T) {
	file, err := os.Create("no_rw")
	if err != nil {
		fmt.Printf("error during the environment setup %s", err)
	}       
	file.Close()
	if err != nil {
		fmt.Printf("error during the environment setup %s", err)
	}       
	os.Chmod("no_rw", 0111)
	answ := Shred("no_rw")
	if answ == nil {
		t.Errorf("Error expected, instead %s", answ)
	}
	os.Remove("no_rw")
}

