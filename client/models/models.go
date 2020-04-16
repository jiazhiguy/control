package models

import (
    "log"
    // "fmt"
    "github.com/asdine/storm/v3"
    // "github.com/boltdb/bolt"
    // "time"
    // "github.com/jinzhu/gorm"
    // _ "github.com/jinzhu/gorm/dialects/mysql"

    // "gin-simpleAgent/pkg/setting"
    // "github.com/minio/minio-go/v6"
)
var Db *storm.DB
var Datachan = make(chan interface{},20)
// var db *gorm.DB
// var Mc *minio.Client
// var BucketName string
// var InitMinioHost string
// type Model struct {
//     ID int `gorm:"primary_key" json:"id"`
//     CreatedOn int `json:"created_on"`
//     ModifiedOn int `json:"modified_on"`
// }

func init() {
    initBoltdb()
    // initMinio()
}
//初始化boltdb
func initBoltdb() error{
    Db, err := storm.Open("my.db")
    go func (){
        for data :=range Datachan {
            log.Printf("%T",data)
            if dataPicture,ok := data.(Picture) ; ok{
                if Db !=nil{
                    log.Println ("Save  ") 
                    err:=Db.Save(&dataPicture)
                    log.Println (err)
                    var pictures []Picture
                    Db.All(&pictures)
                    log.Printf ("%d",len(pictures))    
                }
                
            }
        }      
    }()
    if err !=nil{
        return err
        defer Db.Close()
    }
    return nil
}
//初始化initMinio
// func initMinio(){
//     var err error
//     sec, err := setting.Config.GetSection("minioDb")
//     if err != nil {
//         log.Fatal(2, "Fail to get section 'minioDb': %v", err)
//     }
//     endpoint:= sec.Key("endpoint").String()
//     accessKeyID := sec.Key("accessKeyID").String()
//     secretAccessKey := sec.Key("secretAccessKey").String()
//     BucketName = sec.Key("bucket").String()
//     InitMinioHost=endpoint
//     useSSL := false
//     // Initialize minio client object.
//     Mc, err = minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
//     if err != nil {
//         log.Fatalln(err)
//     }
//     //log.Printf("~~~~~~%#v\n~~~~~", Mc) // minioClient is now setup
//     location := "ap-southeast-2"
//     err = Mc.MakeBucket(BucketName, location)
//     if err != nil {
//         // Check to see if we already own this bucket (which happens if you run this twice)
//         exists, errBucketExists := Mc.BucketExists(BucketName)
//         if errBucketExists == nil && exists {
//             log.Printf("We already own %s\n", BucketName)
//         } else {
//             log.Fatalln(err)
//         }
//     } else {
//         log.Printf("Successfully created %s\n", BucketName)
//     }
// }
// func CloseDB() {
//     defer db.Close()
// }
