package local_pubsub
import (
	"io"
	"crypto/md5"
	// "log"
	"encoding/hex"
)
func Hash(input string) (string,error){
 	hash := md5.New()
	_,err:= hash.Write([]byte(input))
	if err != nil {
		return "",err
	}
	return hex.EncodeToString(hash.Sum(nil)),nil
}
func Md5(src io.Reader)(string,error){
	md5:= md5.New()
		_, err:= io.Copy(md5, src)
			if err != nil {
			return "",err
		}
	md5Str := hex.EncodeToString(md5.Sum(nil))
	return md5Str,nil
}
