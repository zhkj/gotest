package utils

func CheckErr(errMsg error){
	if errMsg != nil{
		panic(errMsg)
	}
}
