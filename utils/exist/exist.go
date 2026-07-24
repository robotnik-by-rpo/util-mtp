package exist

import (
	"fmt"
	"os"
)

// Function for checking file in directory
func IsExistFile(file string) bool {
	_, err := os.Stat(file)

	if err == nil{
		fmt.Println("File exists")
		return true
	} else if os.IsNotExist(err){
		fmt.Println("File doesn't exist in directory")
		return false
	}
	fmt.Printf("Error checking exist file: %v\n",err)	
	return false
}

