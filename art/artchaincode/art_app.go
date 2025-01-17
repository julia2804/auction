
/ckage main
import(
     "bufio"
     "bytes"
      "errors"
      "fmt"
     "strconv"
     "github.com/hyperledger/fabric/core/chaincode/shim"
     "runtime"
     "io"
    "os"
    "strings"
    "image"
   "image/gif"
   "image/png"
   "image/jpeg"
   "net/http"
   "time"    
"crypto/aes"
"crypto/cipher"
 "crypto/rand"
"encoding/json"

)
var recType=[]string{"ARTINV","USER","BID","AUCREQ","POSTRAN","OPENAUC","CLAUC","XFER","VERIFY","TRANS","CFER"}
var MyaucTables=[]string{"MyUserTable","MyUserCatTable","MyAssetHistoryTable","MyAssetTable","MyAssetCatTable","MyAssetAuctionTable","MyBidTable","MyCreditHistoryTable"}
type MyCreditLog struct{
     UserID string
     AuctionedBy string
     Amount string
     RecType string
     Desc string
     Date string
}
type MyAssetTransaction struct{
     AuctionID string
     RecType string
     AssetID string
     TransType string
    UserID string
    TransDate string
    HammerTime string
    HammerPrice string
   Details string
}
type MyBid struct{
   AuctionID string
   RecType string
   BideNo string
   AssetID string
   BuyerID string
   BidPrice string
   BideTime string
}
type MyAuctionRequest struct{
    AuctionID string
    RecType string
   AssetID string
    AuctionHouseID string
    SellerID string
     RequestDate string
     ReservePrice string
     BuyItNowPrice string
     Status string
     OpenDate string
    CloseDate string
}  
type MyAssetLog struct{
    AssetID string
    Status string
    RecType string
    AssetName string
    OwnerID string
   Date  string
   AuctionedBy string
}
type MyAssetObject struct{
    AssetID string
    RecType string
    OwnerID string
    AssetImageName string
    ImageType string
   AssetDate string 
  AssetKind string
    AssetPrice string
    AssetName string
    AES_Key []byte 
    AssetImage []byte
}
type MyUserObject struct{
     UserID string
     RecType string
     UserName    string
    UserType string 
   UserPhone string
     UserLevel string
     UserAmount string
     UserPassward string
}
type SimpleChainCode struct{
}
var gopath string
var ccPath string
func (t *SimpleChainCode) Init(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
      fmt.Println("Trade and Auction Application]Init")
      var err error
       for _,val:=range MyaucTables{
            err=stub.DeleteTable(val)
            if err!=nil{
                return nil,fmt.Errorf("init():delte of %s failed",val)
           }   
           err=InitLedger(stub,val)
           if err!=nil{
                return nil,fmt.Errorf("initledger of %s failed",val)
          }
       }
      err=stub.PutState("version",[]byte(strconv.Itoa(23)))
     if err!=nil{
          return nil,err
      }
    fmt.Println("init() initialization complite:",args)
     return []byte("init():initialization compliet"),nil
}
func InitLedger (stub shim.ChaincodeStubInterface,tableName string) error{
      nKeys:=GetNumberOfKeys(tableName)
      if nKeys<1{
           fmt.Println("Atleast 1 key must be provided\n")
           fmt.Println("Aucion_application:fail creating table",tableName)
          return errors.New("Auction_Application:Failed creating Table"+tableName)
          }
        var columnDefsForTbl []*shim.ColumnDefinition
        for i:=0;i<nKeys;i++{
             columnDef:=shim.ColumnDefinition{Name:"keyName"+strconv.Itoa(i),Type:shim.ColumnDefinition_STRING,Key:true}
             columnDefsForTbl=append(columnDefsForTbl,&columnDef)
        }
      columnLastTblDef:=shim.ColumnDefinition{Name:"Details",Type:shim.ColumnDefinition_BYTES,Key:false}
      columnDefsForTbl=append(columnDefsForTbl,&columnLastTblDef)
     err:=stub.CreateTable(tableName,columnDefsForTbl)
     if err!=nil{
         fmt.Println("Auction_application:fail create table",tableName)
        return errors.New("auction_appliction:fail create table"+tableName)
    }
    return err
}
func ChkReqType(args []string)bool{
       for _,rt:=range args{
               for _,val:=range recType{
                      if val==rt{
                            return true
                       }
                 }
           }
           return false
}
func (t *SimpleChainCode) Invoke(stub shim.ChaincodeStubInterface,function string ,args []string,)([]byte,error){
        var err error
        var buff []byte

      if ChkReqType(args)==true{
           InvokeRequest:=InvokeFunction(function)
           if InvokeRequest!=nil{
                  buff,err=InvokeRequest(stub,function,args)
            }  
      }else{
        fmt.Println("Invoke() Invalid recType:",args[1],"\n")
        return nil,errors.New("Invoke():Invalid recType:"+args[0])

     } 
  return buff,err
}
func main(){
 fmt.Println("hello")
   runtime.GOMAXPROCS(runtime.NumCPU())
   gopath=os.Getenv("GOPATH")
  if len(os.Args)==2 && strings.EqualFold(os.Args[1],"DEV"){
      fmt.Println("---------start in dev mode------")
      ccPath=fmt.Sprintf("%s/src/github.com/hyperledger/fabric/auction/art/myChainCode/",gopath)
    }else{
      fmt.Println("--------strat in net mode------")
       ccPath=fmt.Sprintf("%s/src/github.com/julia2804/auction/art/myChainCode/",gopath)
    }  
     err:=shim.Start(new(SimpleChainCode))
     if err!=nil{
          fmt.Printf("eror staring Simple chaincode:%s",err)
      }
}
func (t *SimpleChainCode) delete(stub shim.ChaincodeStubInterface,args []string)([]byte,error){
      if len(args)!=1{
          return nil,errors.New("incorretc numberof argument")
       }
       A:=args[0]
       err:=stub.DelState(A)
       if err!=nil{
           return nil,errors.New("fail to delete state")
        }
     return nil,nil
}
func CreateAssetObject(args []string)(MyAssetObject,error){
     var err error
     var myAsset MyAssetObject
     if len(args)!=8{
           fmt.Println("CreateAssetObject():incorrect number of argument")
           return myAsset,errors.New("CreateAssetobject():incorrect number of argument")
      }
      _,err=strconv.Atoi(args[0])
      if err!=nil{
          fmt.Println("CreateAssetObject():Id should be interger")
          return myAsset,errors.New("CreateOject():id should be integer")
      }

if err!=nil{
          fmt.Println("something wrong")
          panic(err) 
          }
           imagePath:=ccPath+args[2]
    if _,err:=os.Stat(imagePath);err==nil{
            fmt.Println(imagePath," exist")
    }else {
            fmt.Println("createAccetObject():cannot find or load image",imagePath)
            return myAsset,errors.New("createAssetOjet():ART Picture not found")
     }
     imagebytes,fileType:=imageToByteArray(imagePath)
       fmt.Println("image get succes")   
     AES_key,_:=GenAESKey()
   fmt.Println("genaeskey sucess") 
     AES_enc:=Encrypt(AES_key,imagebytes)
   fmt.Println("encrypt success") 
      myAsset=MyAssetObject{args[0],args[1],args[7],args[2],fileType,args[3],args[4],args[5],args[6],AES_key,AES_enc}
    fmt.Println("CreateAssetObject():Asset object created :ID#",myAsset.AssetID,"\n AES key:",myAsset.AES_Key)
    return myAsset,nil 
}
const(
      AESKeyLength=32
      NonceSize=24
)
func Encrypt(key []byte,ba []byte)[]byte{
     block,err:=aes.NewCipher(key)
     if err!=nil{
         panic(err)
     }
     ciphertext:=make([]byte,aes.BlockSize+len(ba))
     iv:=ciphertext[:aes.BlockSize]
     if _,err:=io.ReadFull(rand.Reader,iv);err!=nil{
        panic(err)
     }
     stream:=cipher.NewCFBEncrypter(block,iv)
     stream.XORKeyStream(ciphertext[aes.BlockSize:],ba)
     return ciphertext
} 
func Decrypt(key []byte,ciphertext []byte)[]byte{
      block,err:=aes.NewCipher(key)
      if err!=nil{
          panic(err)
       }
      if len(ciphertext)<aes.BlockSize{
          panic("text is too short")
       }
       iv:=ciphertext[:aes.BlockSize]
      ciphertext=ciphertext[aes.BlockSize:]
       stream:=cipher.NewCFBDecrypter(block,iv)
       stream.XORKeyStream(ciphertext,ciphertext)
      return ciphertext 
}
func GenAESKey()([]byte,error){
      fmt.Println("enter genskey") 
      return GetRandomBytes(AESKeyLength)
}
func GetRandomBytes(len int)([]byte,error){
   fmt.Println("random bytes success") 
     key:=make([]byte,len)
     _,err:=rand.Read(key)
     if err!=nil{
        return nil,err
     }
   fmt.Println("gen get random sucess ending") 
      return key,nil
}
func ValidateMember(stub shim.ChaincodeStubInterface,owner string)([]byte,error){
     args:=[]string{owner,"USER"}
     Avalbytes,err:=QueryLedger(stub,"MyUserTable",args)
     if err!=nil{
        fmt.Println("ValedateMember():fail -cannot find valid owner recodrd for it",owner)
        jsonResp:="{\"error\":\"fail to ger owner information"+owner+"\"}"
        return nil,errors.New(jsonResp)
}
     if Avalbytes==nil{
         fmt.Println("ValidateMember():fail-imcoplite information",owner)
         jsonResp:="{\"error\":\" fail-imcomplete information"+owner+"\"}"
         return nil,errors.New(jsonResp)
        }
       fmt.Println("validateMember():validateMamber success")
   return Avalbytes,nil
}
func UserToCreditLog(io MyUserObject) MyCreditLog{
       iLog:=MyCreditLog{}
       iLog.UserID=io.UserID
       iLog.AuctionedBy="DEFAULT"
       iLog.Amount=io.UserAmount 
        iLog.Desc="created"
       iLog.Date=time.Now().Format("2006-01-02 15:04:05")
       return iLog
}
func CreditLogtoJSON(credit MyCreditLog)([]byte,error){
       ajson,err:=json.Marshal(credit)
       if err!=nil{
             fmt.Println(err)
             return nil,err
        }
        return ajson,nil
}
func JSONtoCreditLog(ithis []byte)(MyCreditLog,error){
      credit:=MyCreditLog{}
     err:=json.Unmarshal(ithis,&credit)
     if err!=nil{
         fmt.Println("JSONtoCreditLog error:",err)
         return credit,err
     }
     return credit,err
}
func GetCreditLog(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
     if len(args)<1{
          fmt.Println("incorrect aument amont")
           return nil,errors.New("incorrect argument")
      }
      rows,err:=GetList(stub,"MyCreditHistoryTable",args)
      if err!=nil{
             return nil,fmt.Errorf("unmarshal son:%s",err)
      }
      nCol:=GetNumberOfKeys("MyCreditHistoryTable")
      tlist:=make([]MyCreditLog,len(rows))
      for i:=0;i<len(rows);i++{
           ts:=rows[i].Columns[nCol].GetBytes()
           il,err:=JSONtoCreditLog(ts)
           if err!=nil{
               fmt.Println("unmarshall error")
               return nil,fmt.Errorf("operation err:%s",err)
            }
            tlist[i]=il
       }
        jsonRows,_:=json.Marshal(tlist)
        return jsonRows,nil
}
func PostCreditLog(stub shim.ChaincodeStubInterface,user MyUserObject,amount string,ah string)([]byte,error){
      iLog:=UserToCreditLog(user)
      iLog.AuctionedBy=ah
    if ((strings.Compare(amount,"0"))!=0){
           iLog.Desc="ammented by "+ah 
           iLog.Amount=amount 
      } else {
           iLog.Desc="updated automatically"
           iLog.Amount="0" 
      } 
      buff,err:=CreditLogtoJSON(iLog)
      if err!=nil{
           fmt.Println("fail to create:",user.UserID)
           return nil,errors.New("failto create "+user.UserID)
       }else{
           keys:=[]string{iLog.UserID,iLog.AuctionedBy,time.Now().Format("2016-01-02 15:04:05")}
      err=UpdateLedger(stub,"MyCreditHistoryTable",keys,buff)
           if err!=nil{
             fmt.Println("write error")
              return buff,err
           }
      }
      return buff,nil
}
func PostAsset(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
    assetObject,err:=CreateAssetObject(args[0:])
   if err!=nil{
      fmt.Println("PostAsset():cannot create item object\n")
      return nil,err
    }
    ownerInfo,err:=ValidateMember(stub,assetObject.OwnerID)
    fmt.Println("owner information",ownerInfo,assetObject.OwnerID)
    if err!=nil{
        fmt.Println("postAsset():failed woner information not found",assetObject.OwnerID)
      }
   buff,err:=ARtoJSON(assetObject)
   if err!=nil{
       fmt.Println("PostAsset():fail cannot create object buff for write:",args[1])
      return nil,errors.New("PostAsset():fail cannot create objet buffer for write:"+ args[1])
     }else {
          keys:=[]string{args[0]}
          err=UpdateLedger(stub,"MyAssetTable",keys,buff)
          if err!=nil{
               fmt.Println("PostAsset():write error while insert\n")
               return buff,err
          }
      _,err=PostAssetLog(stub,assetObject,"INITIAL","DEFAULT")
     if err!=nil{
          fmt.Println("PostAssetLog():write error")
         return nil,err
       }
     fmt.Println("the args[5]:",args[5])
      keys=[]string{"2016",args[4],args[0]}
      err=UpdateLedger(stub,"MyAssetCatTable",keys,buff)
    if err!=nil{
         fmt.Println("PostAsset():write error")
         return buff,err
      }
  }  
     secret_key,_:=json.Marshal(assetObject.AES_Key)
    fmt.Println(string(secret_key))
    return secret_key,nil
}

func AssetToAssetLog(io MyAssetObject) MyAssetLog{
    iLog:=MyAssetLog{}
    iLog.AssetID=io.AssetID
    iLog.Status="INITIAL"
  iLog.AuctionedBy="DEFAULT" 
    iLog.RecType="ALOG"
    iLog.AssetName=io.AssetName
    iLog.OwnerID=io.OwnerID
    iLog.Date=time.Now().Format("2017-03-22 16:33:09")
    return iLog
} 

func PostAssetLog(stub shim.ChaincodeStubInterface,asset MyAssetObject,status string,ah string)([]byte,error){
     iLog:=AssetToAssetLog(asset)
     iLog.Status=status
     iLog.AuctionedBy=ah
     buff,err:=AssetLogtoJSON(iLog)
     if err!=nil{
        fmt.Println("PostAssetLog():failed cannotcreate object buffer "+asset.AssetID)
      return nil,errors.New("PostAssetLog():failed cannot create object"+asset.AssetID)
    }else {
         keys:=[]string{iLog.AssetID,iLog.Status,iLog.AuctionedBy,time.Now().Format("2017-03-22 16:33:09")}
        err=UpdateLedger(stub,"MyAssetHistoryTable",keys,buff)
       if err!=nil{
            fmt.Println("PostAssetLog():write error")
           return buff,err
       }
   }
 return buff,nil
}
func AssetLogtoJSON(asset MyAssetLog)([]byte,error){
     ajson,err:=json.Marshal(asset)
     if err!=nil{
        fmt.Println(err)
        return nil,err
     }
    return ajson,nil
} 
func JSONtoArgs(Avalbytes []byte)(map[string]interface{},error){
     var data map[string]interface{}
     if err:=json.Unmarshal(Avalbytes,&data);err!=nil{
             return nil,err
     }
     return data,nil
} 
func JSONtoAR(data []byte)(MyAssetObject,error){
     ar:=MyAssetObject{}
     err:=json.Unmarshal([]byte(data),&ar)
     if err!=nil{
        fmt.Println("Unmarshal failed:",err)
     }
     return ar,err
}
func ByteArrayToImage(imgByte []byte,imageFile string)error{
     img,_,_:=image.Decode(bytes.NewReader(imgByte))
     fmt.Println("processQueryResult byteArrayToImage:proceeding to create image")
     out,err:=os.Create(imageFile)
     if err!=nil{
        fmt.Println("byteArrayToImage():cannot crate image file ",err)
        return errors.New("byteArrayToImage():proced image file failed")
      }
      fmt.Println("processQueryType byteArrayToImage:proceding to encode image")
      filetype:=http.DetectContentType(imgByte)
      switch filetype{
      case "image/jpeg","image/jpg":
            var opt jpeg.Options
            opt.Quality=100
               err=jpeg.Encode(out,img,&opt)
      case "image/gif":
           var opt gif.Options
             opt.NumColors=256
           err=gif.Encode(out,img,&opt)
     case "image/png":
           err=png.Encode(out,img)
      default:
           err=errors.New("ohly pmng,jpg and gif supported")
      }
      if err!=nil{
         fmt.Println("ByteArrayToImage():cannot encode image file ",err)
         return errors.New("buteArrayToImage():cannot encode image file ")
      }
     fmt.Println("image filegenrated and saved to ",imageFile)
     return nil
}
func ARtoJSON(ar MyAssetObject)([]byte,error){
    fmt.Println("ar to json:",ar.AES_Key) 
      ajson,err:=json.Marshal(ar)
      if err!=nil{
         fmt.Println(err)
         return nil,err
      }
     return ajson,nil
}
func QueryFunction(fname string) func(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
fmt.Println("enter funtion")
fmt.Println("fanme:",fname)
    QueryFunc:=map[string]func(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
       "GetAsset": GetAsset,
       "GetUser": GetUser,
      "GetUserListByCat":GetUserListByCat,
      "GetAssetListByCat":GetAssetListByCat,
      "GetAssetLog":GetAssetLog,
    "ValidateItemOwnership":ValidateItemOwnership,  
   "GetCreditLog":GetCreditLog,  
     "ValidateUser":ValidateUser, 
    }
    return QueryFunc[fname]
}
func GetAssetLog(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
       if len(args)<1{
            fmt.Println("getAssetLog():incorect argument")
            return nil,errors.New("incorrect argumanet")
        }
       rows,err:=GetList(stub,"MyAssetHistoryTable",args)
       if err!=nil{
          return nil,fmt.Errorf("error marshal json:%s",err)
       }
       nCol:=GetNumberOfKeys("MyAssetHistoryTable")
       tlist:=make([]MyAssetLog,len(rows))
        for i:=0;i<len(rows);i++{
            ts:=rows[i].Columns[nCol].GetBytes()
            il,err:=JSONtoAssetLog(ts)
            if err!=nil{
                fmt.Println("unmarshalerror")
                return nil,fmt.Errorf("operation err:%s",err)
             }
            tlist[i]=il
        }
        jsonRows,_:=json.Marshal(tlist)
        return jsonRows,nil
} 
func ValidateItemOwnership(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
       var err error
       if len(args)<3{
            fmt.Println("item,owner,key needer")
            return nil,errors.New("request 3 arguent")
       }
       Avalbytes,err:=QueryLedger(stub,"MyAssetTable",[]string{args[0]})
      if err!=nil{
           fmt.Println("failed to query")
           jsonResp:="{\"error\" get q object data for "+args[0]+"\"}"
           return nil,errors.New(jsonResp)
       }
      if Avalbytes==nil{
           fmt.Println("fail imcoplete query")
            jsonResp:="{error imcomplete informateio for "+args[0]+"\"}"
            return nil,errors.New(jsonResp)
       }
      myItem,err:=JSONtoAR(Avalbytes)
     if err!=nil{
          fmt.Println("faile myitem")
          jsonResp:="{\"error\" get data for(item) "+args[0]+"\"}"
          return nil,errors.New(jsonResp)
        }
       myKey:=GetKeyValue(Avalbytes,"AES_Key")
     myName:=GetKeyValue(Avalbytes,"AssetName")
      myID:=GetKeyValue(Avalbytes,"AssetID") 
        fmt.Println("name string:",myName)
        fmt.Println("id string:",myID) 
        fmt.Println("key string:=",myKey)
       if myKey!=args[2]{
            fmt.Println("key not match",args[2],"-",myKey)
            jsonResp:="{\"error\" jey not match "+args[0]+"\"}"
            return nil,errors.New(jsonResp)
       }
        if myItem.OwnerID!=args[1]{
            fmt.Println("owner not march ",args[1])
            jsonResp:="{\"error\" owner id not marh"+args[0]+"\"}"
           return nil,errors.New(jsonResp) 
       } 
     fmt.Println("successful")
  return Avalbytes,nil
} 
func GetKeyValue(Avalbytes []byte,key string) string{
      var dat map[string]interface{}
      if err:=json.Unmarshal(Avalbytes,&dat);err!=nil{
              panic(err)
      }
     val:=dat[key].(string)
     return val
}
func TransferItem(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
      var err error
      if len(args)<5{
         fmt.Println("Translf():item,ownerid,key,new owner,moner")
          return nil,errors.New("dtransfer new 5 argument")
          }
        err=VerifyIfItemIsOnAuction(stub,args[0])
       if err!=nil{
            fmt.Println("faied",args[0])
            return nil,err
          }
        _,err=ValidateMember(stub,args[3])
        if err!=nil{
              fmt.Println("item not registe yet",args[3])
              return nil,err
        }
      ar,err:=ValidateItemOwnership(stub,"ValidateItemOwnership",args[:3])
      if err!=nil{
            fmt.Println("transfer r fail to authenticate:")
            return nil,err
        }
        myItem,err:=JSONtoAR(ar)
        if err!=nil{
           fmt.Println("faile create item from josn")
           return nil,err
        }
       CurrentAES_Key:=myItem.AES_Key
       image:=Decrypt(CurrentAES_Key,myItem.AssetImage)
      myItem.AES_Key,_=GenAESKey()
    myItem.AssetImage=Encrypt(myItem.AES_Key,image) 
      myItem.OwnerID=args[3]
     ar,err=ARtoJSON(myItem)
      keys:=[]string{myItem.AssetID,myItem.OwnerID}
      err=ReplaceLedgerEntry(stub,"MyAssetTable",keys,ar)
  if err!=nil{
    fmt.Println("transferasset failed to replsde")
    return nil,err
   }
   fmt.Println("transferasset sucess")
  keys=[]string{"2016",myItem.AssetKind,myItem.AssetID}
    err=ReplaceLedgerEntry(stub,"MyAssetCatTable",keys,ar)
    if err!=nil{
       fmt.Println("failed to replace asset at table")
       return nil,err
     }
     _,err=PostAssetLog(stub,myItem,"Transfer",args[1])
     if err!=nil{
         fmt.Println("write error post asset log")
        return nil,err
     }
    fmt.Println("myitem keys:",myItem.AES_Key) 
    fmt.Println("replace cat table success") 
    return myItem.AES_Key,nil
} 
func ReplaceLedgerEntry(stub shim.ChaincodeStubInterface,tableName string,keys []string,args []byte)error{
        nKey:=GetNumberOfKeys(tableName)
        if nKey<1{
             fmt.Println("at lest 1 key")
          }
          var columns []*shim.Column
         for i:=0;i<nKey;i++{
             col:=shim.Column{Value:&shim.Column_String_{String_:keys[i]}}
             columns=append(columns,&col)
         }
         lastCol:=shim.Column{Value:&shim.Column_Bytes{Bytes:[]byte(args)}}
         columns=append(columns,&lastCol)
         row:=shim.Row{columns}
          ok,err:=stub.ReplaceRow(tableName,row)
          if err!=nil{
                return fmt.Errorf("replace row into "+tableName+" able operation failed.%s",err)
          }
          if !ok{
                 return errors.New("replace row into "+tableName+" tabelf failed.Row with given key"+keys[0]+"alreay exist")
         }
        fmt.Println("replace row in "+tableName+" table success")
        return nil
} 
func VerifyIfItemIsOnAuction(stub shim.ChaincodeStubInterface,itemID string)error{
    return nil
}
func InvokeFunction(fname string) func(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
    InvokeFunc:=map[string]func(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
        "PostAsset":    PostAsset,
        "PostUser":     PostUser,
       "TransferCredit":TransferCredit, 
        "TransferItem":TransferItem,
        }
       return InvokeFunc[fname]
}
func ProcessQueryResult(stub shim.ChaincodeStubInterface,Avalbytes []byte,args []string)error{
      var dat map[string]interface{}
      if err:=json.Unmarshal(Avalbytes,&dat);err!=nil{
            panic(err)
        }
      var recType string
      recType=dat["RecType"].(string)
      switch recType{
      case "ARTINV":
          ar,err:=JSONtoAR(Avalbytes)
          if err!=nil{
                   fmt.Println("ProcessRequestType():Cannot creae assetObject \n")
                return err
            }
           image:=Decrypt(ar.AES_Key,ar.AssetImage)
           if err!=nil{
                fmt.Println("processRequestType():image decryption faied")
               return err
           }
          fmt.Println("ProcessRequestType():Image conversion sucessfull")
         err=ByteArrayToImage(image,ccPath+"copy."+ar.AssetImageName)
       if err!=nil{
            fmt.Println("ProcessRequestType():image conversion fail")
            return err
       }
      return err
     case "USER":
          ur,err:=JSONtoUser(Avalbytes)
          if err!=nil{
             return err
          }
         fmt.Println("ProcessRequestType():",ur)
        return err
     case "AUCREQ":
     case "OPENAUC":
     case "CLAUC":
         ar,err:=JSONtoAucReq(Avalbytes)
         if err!=nil{
              return err
         }
        fmt.Println("ProcessRequestType():",ar)
       return err
     case "POSTTRAN":
         atr,err:=JSONtoTran(Avalbytes)
         if err!=nil{
             return err
          }
          fmt.Println("PrcessRequestType():",atr)
     case "BID":
         bid,err:=JSONtoBid(Avalbytes)
         if err!=nil{
             return err
            }
          fmt.Println("processRequestType():",bid)
          return err
     case "DEFAULT":
          return nil
    case "XFER":
          return nil
    case "CFER":
          return nil 
      case "VERIFY":
           return nil
      default:
          return errors.New("unknown")
      }
     return nil
} 
func GetAsset(stub shim.ChaincodeStubInterface,function string,args []string)([]byte ,error){
     Avalbytes,err:=QueryLedger(stub,"MyAssetTable",args)
     if err!=nil{
        fmt.Println("gerAsser():fail to uery object")
        jsonResp:="{\"error\":\"fail to ger data for "+args[0]+"\"}"
        return nil,errors.New(jsonResp)
     }
     if Avalbytes==nil{
        fmt.Println("ger asset():incomplet query")
        jsonResp:="{\"err\":\"incomplete query" +args[0]+"\"}"
        return nil,errors.New(jsonResp)
     }
fmt.Println("get asset:response:success")
     assetObj,_:=JSONtoAR(Avalbytes)
     assetObj.AssetImage=[]byte{}
    fmt.Println("get asset:aeskdy:",assetObj.AES_Key) 
      Avalbytes,_=ARtoJSON(assetObj)
     return Avalbytes,nil
    }

func(t *SimpleChainCode) Query(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
    var err error
    var buff []byte
    fmt.Println("ID extracted and type= ",args[0])
    fmt.Println("Args supp;;",args)
    if len(args)<1{
         fmt.Println("at lest 1 arguments key")
         return nil, errors.New("query():expecting transaction type")       }
    QueryRequest:=QueryFunction(function)
    if QueryRequest!=nil{
       buff,err=QueryRequest(stub,function,args)
    }else {
        fmt.Println("query() invalid function call:",function)
        return nil,errors.New("Query():invalid functio acll:" +function)    }
if err!=nil{
    fmt.Println("query() object ot found:",args[0])
    return nil,errors.New("not found:" +args[0])
     }
    return buff,err
   }
func UpdateUserObject(stub shim.ChaincodeStubInterface,ar []byte,hammerUser string,amount string)(string,error){
      var err error
      myUser,err:=JSONtoUser(ar)
      if err!=nil{
           fmt.Println("fail to create")
          return "wrong",err
      }
     number,error:=strconv.Atoi(amount)
    if error!=nil{
          fmt.Println("tarandform failed")
     }
    amo,error:=strconv.Atoi(myUser.UserAmount)
   if error!=nil{
         fmt.Println("trander amount of user fialed")
    }
    amo=amo+number
    myUser.UserAmount=strconv.Itoa(amo)
    ar,err=UsertoJSON(myUser)
    keys:=[]string{myUser.UserID}
    err=ReplaceLedgerEntry(stub,"MyUserTable",keys,ar)
    if err!=nil{
        fmt.Println("fail to replace ledger")
        return "",err
     }
    fmt.Println("repleace user table succesfull")
    keys=[]string{"2016",myUser.UserType,myUser.UserID}
    err=ReplaceLedgerEntry(stub,"MyUserCatTable",keys,ar)
    if err!=nil{
          fmt.Println("replace cat table failed")
          return "",err
    }
     fmt.Println("succesful")
    return myUser.UserID,nil
} 
func TransferCredit(stub shim.ChaincodeStubInterface,function string,args[]string)([]byte,error){
       var err error
      if len(args)<5{
             fmt.Println("transferItem():argument wrong")
             return nil,errors.New("argument wrong")
        }
   ar,err:=ValidateMember(stub,args[0])
   if err!=nil{
      fmt.Println("gail valide", args[0])
      return nil,err
   }
   _,err2:=ValidateUser(stub,"ValidateUser",args)
   if err2!=nil{
       fmt.Println("gail authentic", args[2])
      return nil,err2
   }
  myUser,err:=JSONtoUser(ar)
   if err!=nil{
      fmt.Println("faile tao marshall")
      return nil,err
   }
  str:=myUser.UserAmount
 amo,_:=strconv.Atoi(str)
 count,_:=strconv.Atoi(args[1]) 
  amo=amo+count 
  string_amount:=strconv.Itoa(amo)
   myUser.UserAmount=string_amount 
   ar,err=UsertoJSON(myUser)
   keys:=[]string{myUser.UserID} 
   err=ReplaceLedgerEntry(stub,"MyUserTable",keys,ar)
    if err!=nil{
      fmt.Println("faile to replace user table")
      return nil,err
   }
   fmt.Println("success")
   keys=[]string{"2016",myUser.UserType,myUser.UserID}
   err=ReplaceLedgerEntry(stub,"MyUserCatTable",keys,ar)
  if err!=nil{      fmt.Println("faile to replace user cat table")
      return nil,err
   }
  _,err=PostCreditLog(stub,myUser,args[1],args[2])
   if err!=nil{ 
     fmt.Println("faile to replace user table")
      return nil,err
   }
  _,err=ValidateLevel(stub,"ValidateLevel",keys)
   if err!=nil{
     fmt.Println("faile to replace user table level")
   }
   fmt.Println("repleae success")
   return ar,nil
}
func ValidateLevel(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
{
  var err error
  arg:=[]string{args[2]}
   Avalbytes,err:=QueryLedger(stub,"MyUserTable",arg)
    if err!=nil{
           fmt.Println("fail to quey")
           jsonResp:="{\"error\":get data for"+args[2]+"\"}"
           return nil,errors.New(jsonResp)
    }
   if Avalbytes==nil{
        fmt.Println("incomplete query ojedt")
        jsonResp:="{\"error\":\"get data avalbtes err for "+args[2]+"\"}"
        return nil,errors.New(jsonResp)
   }
   myUser,err:=JSONtoUser(Avalbytes)
 if err!=nil{
           fmt.Println("fail to marshal")
           jsonResp:="{\"error\":\"get marshal for"+args[2]+"\"}"
           return nil,errors.New(jsonResp)
    }
   leve,_:=strconv.Atoi(myUser.UserLevel) 
    fmt.Println("amount:",myUser.UserAmount) 
    amo,_:=strconv.Atoi(myUser.UserAmount)
  fmt.Println("amount number",amo) 
     if (amo>1000 && leve==1 && amo<=3000){
           myUser.UserLevel="2"
            Avalbytes,err=UsertoJSON(myUser) 
            keys:=[]string{myUser.UserID,myUser.UserLevel} 
            err=ReplaceLedgerEntry(stub,"MyUserTable",keys,Avalbytes)  
            if err!=nil{
                fmt.Println("update user level:failed")
             }
            keys=[]string{"2016",myUser.UserType,myUser.UserID}
            err=ReplaceLedgerEntry(stub,"MyUserCatTable",keys,Avalbytes) 
            if err!=nil{
                fmt.Println("update user level:failed")
                return nil,err  
            }
            _,err=PostCreditLog(stub,myUser,"0",args[2])
             if err!=nil{
                fmt.Println("post creditlog error")
                 return nil,err
             }  
            fmt.Println("success update user level")
      } else if(amo>3000 && (leve==1 || leve==2)){
              myUser.UserLevel="3"
             Avalbytes,err=UsertoJSON(myUser)
              keys:=[]string{myUser.UserID,myUser.UserLevel} 
            err=ReplaceLedgerEntry(stub,"MyUserTable",keys,Avalbytes)
            if err!=nil{
                fmt.Println("update user level:failed")
             }
            keys=[]string{"2016",myUser.UserType,myUser.UserID}
            err=ReplaceLedgerEntry(stub,"MyUserCatTable",keys,Avalbytes)
            if err!=nil{
                fmt.Println("update user cat level:failed",err)
             }
            _,err=PostCreditLog(stub,myUser,"0",args[2])
             if err!=nil{
                fmt.Println("post creditlog error")
                 return nil,err
             } 
             fmt.Println("success update user level")
        }
     return Avalbytes,nil
}
}
func ValidateUser(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
     var err error
    arg:=[]string{args[2]} 
     Avalbytes,err:=QueryLedger(stub,"MyUserTable",arg)
    if err!=nil{
           fmt.Println("fail to quey")
           jsonResp:="{\"error\":get data for"+args[2]+"\"}"
           return nil,errors.New(jsonResp)
    }
   if Avalbytes==nil{
        fmt.Println("incomplete query ojedt")
        jsonResp:="{\"error\":\"get data avalbtes err for "+args[2]+"\"}"
        return nil,errors.New(jsonResp)
   }
   myUser,err:=JSONtoUser(Avalbytes)
 if err!=nil{
           fmt.Println("fail to marshal")
           jsonResp:="{\"error\":\"get marshal for"+args[2]+"\"}"
           return nil,errors.New(jsonResp)
    }
   if strings.Compare(myUser.UserType,"商家")!=0{
         fmt.Println("this user has no right tranfer redit")
         jsonResp:="{\"error\":\"not right to transfer"+args[2]+"\"}"
         return nil,errors.New(jsonResp)
   }
   fmt.Println("valide success")
    return Avalbytes,nil
}
func CreateUserObject(args []string)(MyUserObject,error){
      var err error
      var aUser MyUserObject
      if len(args)!=8{
        fmt.Println("CreateUserObject():incorrect argument number,expecting 8")
      }
      _,err=strconv.Atoi(args[0])
     if err!=nil{
         return aUser,errors.New("CreateUserObject():Incorrect number of user id")
    }
    aUser=MyUserObject{args[0],args[1],args[2],args[3],args[4],args[5],args[6],args[7]}
    fmt.Println("CreateUserObject():User Object:",aUser)
  return aUser,nil
} 
func UsertoJSON(user MyUserObject)([]byte,error){
   fmt.Println("user.phone:",user.UserPhone) 
    ajson,err:=json.Marshal(user)
    if err!=nil{
      fmt.Println("UserJSON error:",err)
      return nil,err
      }
     fmt.Println("UsertoJSON created:",ajson)
    record,_:=JSONtoUser(ajson) 
   fmt.Println("record.level,userjson:",record.UserPhone) 
     return ajson,nil
} 
func GetNumberOfKeys(tname string )int{
     TableMap:=map[string]int{
            "MyUserTable":   1,
            "MyUserCatTable": 3,
            "MyAssetCatTable": 3,
            "MyAssetTable":    1,
            "MyAssetHistoryTable":4,
           "MyBidTable":2,
           "MyTransTable":2,
          "MyAssetAuctionTable":1, 
          "MyCreditHistoryTable":3, 
          }
      return TableMap[tname]
}
func QueryLedger2(stub shim.ChaincodeStubInterface,tableName string,args []string)([]byte,error){
      var columns []shim.Column
      nCol:=GetNumberOfKeys(tableName)
      for i:=0;i<nCol;i++{
           colNext:=shim.Column{Value: &shim.Column_String_{String_:args[i]}}
           columns=append(columns,colNext)
}
      row,err:=stub.GetRow(tableName,columns)
     fmt.Println("Length or number of rows retrived ",len(row.Columns))
     if len(row.Columns)==0{
           jsonResp:="{\"error\":\" fail retrieving data"+args[0]+" .\"}"
          fmt.Println("error retriving data record for key="+args[0],"error:",jsonResp)
        return nil,errors.New(jsonResp)
        }
      Avalbytes:=row.Columns[nCol].GetBytes()
      fmt.Println("QueryLedger():successful-proceeding to process quest type")
      err=ProcessQueryResult(stub,Avalbytes,args)
     if err!=nil{
           fmt.Println("QueryLedger():cannot create object:",args[1])
    jsonResp:="{\"QueryLedger()error\":\" cannot create object for key"+args[0]+"\"}"
    return nil,errors.New(jsonResp)
    }
  return Avalbytes,nil
}

func UpdateLedger(stub shim.ChaincodeStubInterface,tableName string,keys []string,args []byte) error{
     nKeys:=GetNumberOfKeys(tableName)
     var record MyUserObject
     var err error 
     fmt.Println("the compare result:",strings.Compare(tableName,"MyUserTable")) 
     if strings.Compare(tableName,"MyUserTable")==0{
       record,err=JSONtoUser(args)
            fmt.Println("record.phone:",record.UserPhone)
        } 
     fmt.Println("iam entering updateledger")
     if nKeys<1{
        fmt.Println("At least 1 key must be porovide\n")
     }
     var columns []*shim.Column
     for i:=0;i<nKeys;i++{
        col:=shim.Column{Value:&shim.Column_String_{String_:keys[i]}}
        columns=append(columns,&col)
       }
    lastCol:=shim.Column{Value:&shim.Column_Bytes{Bytes:[]byte(args)}}
    columns=append(columns,&lastCol)
    row:=shim.Row{columns}
fmt.Println("i am inserting") 
   ok,err:=stub.InsertRow(tableName,row)
     if err!=nil{
fmt.Println(" inser row,err")
        return fmt.Errorf("UpdateLedger:InsertRow into "+tableName+" Table operateon failed,%s",err)
      }
     if !ok{
      fmt.Println("insert ok,but,existed")   
           return errors.New("UpdateLedger:Insert Row into"+tableName+" Table failed.Given keys"+keys[0]+"already existed")
        }
  fmt.Println("UpdateLedger:InsertRoew into "+tableName +" table successufull")
    return nil
}
func imageToByteArray(imageFile string)([]byte,string){
       file,err:=os.Open(imageFile)
       if err!=nil{
           fmt.Println("imageToByteAraay():cannot open image file",err)
         return nil,string("imageToByteAray():cannot open image file")
}
     defer file.Close()
     fileInfo,_:=file.Stat()
     var size int64=fileInfo.Size()
     bytes:=make([]byte,size)
     buff:=bufio.NewReader(file)
     _,err=buff.Read(bytes)
      if err!=nil{     
      return nil,string("imageToByteArray():cannot read image")
        }
     filetype:=http.DetectContentType(bytes)
     fmt.Println("imageToByteArray():",filetype)
     return bytes,filetype
}
func JSONtoAssetLog(ithis []byte)(MyAssetLog,error){
     item:=MyAssetLog{}
     err:=json.Unmarshal(ithis,&item)
     if err!=nil{
         fmt.Println("log error:",err)
         return item,err
     }
     return item,err
}
func JSONtoUser(user []byte)(MyUserObject,error){
     ur:=MyUserObject{}
     err:=json.Unmarshal(user,&ur)
     if err!=nil{
        fmt.Println("JSONtoUsr error:",err)
        return ur,err
     }
     fmt.Println("JSONtoUser created:",ur)
     return ur,err
}
func JSONtoAucReq(areq []byte)(MyAuctionRequest,error){
     ar:=MyAuctionRequest{}
     err:=json.Unmarshal(areq,&ar)
     if err!=nil{
        fmt.Println("JSONtoAucReq error:",err)
        return ar,err
      }
     return ar,err
}
func JSONtoBid(areq []byte)(MyBid,error){
      myHand:=MyBid{}
       err:=json.Unmarshal(areq,&myHand)
      if err!=nil{
          fmt.Println("JSONtoBid error:",err)
          return myHand,err
        }
       return myHand,err
}
func JSONtoTran(areq []byte)(MyAssetTransaction,error){
     at:=MyAssetTransaction{}
     err:=json.Unmarshal(areq,&at)
     if err!=nil{
         fmt.Println("JSONtoTran error:",err)
         return at,err
     }
     return at,err
}
func GetUser(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
var err error
Avalbytes,err:=QueryLedger(stub,"MyUserTable",args)
if err!=nil{
    fmt.Println("GetUser():Failed to Qeuery objet")
     jsonResp:="{\"Error\":\"failed to ger oject data faor "+args[0]+"\"}"
     return nil,errors.New(jsonResp)
}
if Avalbytes==nil{
    fmt.Println("Get user():inomplet query")
    jsonResp:="{\"error\":\"imcomplete for "+args[0]+"\"}"
    return nil,errors.New(jsonResp)
}
fmt.Println("GetUSr():successful")
return Avalbytes,nil
}
func QueryLedger(stub shim.ChaincodeStubInterface,tableName string,args []string)([]byte,error){
fmt.Println("enter query elgder")
    var columns []shim.Column
    nCol:=GetNumberOfKeys(tableName)
fmt.Println("ncol:",nCol,".tabelName:",tableName)
    for i:=0;i<nCol;i++{
        colNext:=shim.Column{Value: &shim.Column_String_{String_:args[i]}}
        columns=append(columns,colNext)
       }
     fmt.Println("append successful") 
     row,err:=stub.GetRow(tableName,columns)
    fmt.Println("lenth or number of rows retrieved",len(row.Columns))
if len(row.Columns)==0{
     jsonResp:="{\"error\":\"failed retrieving data "+args[0] +".\"}"
     fmt.Println("fail retrieving for key = ",args[0],"error",jsonResp)
     return nil,errors.New(jsonResp)
}
Avalbytes:=row.Columns[nCol].GetBytes()
fmt.Println("successful")
err=ProcessQueryResult(stub,Avalbytes,args)
if err!=nil{
      fmt.Println("queryLedger():cannot create object:",args[1])
      jsonResp:="{\"error\":\"cannot create "+args[0] +"\"}"
     return nil,errors.New(jsonResp)
}
  return Avalbytes,nil
}
func GetList(stub shim.ChaincodeStubInterface,tableName string,args []string)([]shim.Row,error){
       var columns []shim.Column
        nKeys:=GetNumberOfKeys(tableName)
        nCol:=len(args)
         if nCol<1{
               fmt.Println("at least one key\n")
               return nil,errors.New("getlist failed")
         }
               for i:=0;i<nCol;i++{
                  colNext:=shim.Column{Value:&shim.Column_String_{String_:args[i]}}
                  columns=append(columns,colNext)
                }
         rowChannel,err:=stub.GetRows(tableName,columns)
         if err!=nil{
               return nil,fmt.Errorf(" operation fail.%s",err)
         }
         var rows []shim.Row
         for{
              select{
                case row,ok:=<-rowChannel:
                      if !ok{
                            rowChannel=nil
                       }else{
                             rows=append(rows,row)
                       }
                 }
                 if rowChannel==nil{
                       break
                 }
            }
            fmt.Println("keys retrieved:",nKeys)
            fmt.Println("rows retrieved:",len(rows))
            return rows,nil
}
func GetUserListByCat(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
      if len(args)<1{
           fmt.Println("getUSerListBycat():incorect number of arguents")
           return nil,errors.New("ccreateUserObjetct():incorrect")
       }
       rows,err:=GetList(stub,"MyUserCatTable",args)
       if err!=nil{
              return nil,fmt.Errorf("ger failed.error marshaling json:%s",err)
      }
      nCol:=GetNumberOfKeys("MyUserCatTable")
      tlist:=make([]MyUserObject,len(rows))
      for i:=0;i<len(rows);i++{
          ts:=rows[i].Columns[nCol].GetBytes()
          uo,err:=JSONtoUser(ts)
          if err!=nil{
              fmt.Println("GerUserListByCat()failed:ummarshasll error")
              return nil,fmt.Errorf("operaion faile.%s",err)
           }
           tlist[i]=uo
       }
       jsonRows,_:=json.Marshal(tlist)
       return jsonRows,nil
}
func GetAssetListByCat(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
    if len(args)<1{
           fmt.Println("getAssetListByCat():ncorrect arguent")
           return nil,errors.New("incorect argument")
      }
      rows,err:=GetList(stub,"MyAssetCatTable",args)
      if err!=nil{
           return nil,fmt.Errorf("getItem List cat .errot gerlist:%s",err)
       }
       nCol:=GetNumberOfKeys("MyAssetCatTable")
       tlist:=make([]MyAssetObject,len(rows))
       for i:=0;i<len(rows);i++{
                ts:=rows[i].Columns[nCol].GetBytes()
                io,err:=JSONtoAR(ts)
                if err!=nil{
                     fmt.Println("unmarshall error")
                     return nil,fmt.Errorf("operation fial.%s",err)
                 }
                 io.AssetImage=[]byte{}
                 fmt.Println("list asset,aes-key:",io.AES_Key)  
                tlist[i]=io
        }
       jsonRows,_:=json.Marshal(tlist)
       return jsonRows,nil
}
func PostUser(stub shim.ChaincodeStubInterface,function string,args []string)([]byte,error){
       record,err:=CreateUserObject(args[0:])
        fmt.Println("args[5]",args[5])
      if err!=nil{
         return nil,err
     }
     buff,err:=UsertoJSON(record)
     if err!=nil{
        fmt.Println("PostUserObject():failed cannot create object:",args[1])
        return nil,errors.New("PostUserObject():faile cannot create object:"+args[1])
     }else {
          keys:=[]string{args[0]}
             err=UpdateLedger(stub,"MyUserTable",keys,buff)
           if err!=nil{
              fmt.Println("PostUser():write error while inserting recode")
              return nil,err
            }
          _,err=PostCreditLog(stub,record,record.UserAmount,"DEFAULT")
          if err!=nil{
               fmt.Println("Postcredit og werite error")
               return nil,err
         } 
          keys=[]string{"2016",args[3],args[0]}
     err=UpdateLedger(stub,"MyUserCatTable",keys,buff)
     if err!=nil{
         fmt.Println("PostUser():write error wihle inserting recode into usercatTable")
      }
   }
   return buff,err
}
******************************************************************
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
******************************************************************/

///////////////////////////////////////////////////////////////////////
// Author : Mohan Venkataraman
// Purpose: Explore the Hyperledger/fabric and understand
// how to write an chain code, application/chain code boundaries
// The code is not the best as it has just hammered out in a day or two
// Feedback and updates are appreciated
///////////////////////////////////////////////////////////////////////

package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/op/go-logging"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	// "github.com/errorpkg"
)

//////////////////////////////////////////////////////////////////////////////////////////////////
// The recType is a mandatory attribute. The original app was written with a single table
// in mind. The only way to know how to process a record was the 70's style 80 column punch card
// which used a record type field. The array below holds a list of valid record types.
// This could be stored on a blockchain table or an application
//////////////////////////////////////////////////////////////////////////////////////////////////
var recType = []string{"ARTINV", "USER", "BID", "AUCREQ", "POSTTRAN", "OPENAUC", "CLAUC", "XFER", "VERIFY"}

//////////////////////////////////////////////////////////////////////////////////////////////////
// The following array holds the list of tables that should be created
// The deploy/init deletes the tables and recreates them every time a deploy is invoked
//////////////////////////////////////////////////////////////////////////////////////////////////
var aucTables = []string{"UserTable", "UserCatTable", "ItemTable", "ItemCatTable", "ItemHistoryTable", "AuctionTable", "AucInitTable", "AucOpenTable", "BidTable", "TransTable"}

///////////////////////////////////////////////////////////////////////////////////////
// This creates a record of the Asset (Inventory)
// Includes Description, title, certificate of authenticity or image whatever..idea is to checkin a image and store it
// in encrypted form
// Example:
// Item { 113869, "Flower Urn on a Patio", "Liz Jardine", "10102007", "Original", "Floral", "Acrylic", "15 x 15 in", "sample_9.png","$600", "My Gallery }
///////////////////////////////////////////////////////////////////////////////////////

type ItemObject struct {
	ItemID         string
	RecType        string
	ItemDesc       string
	ItemDetail     string // Could included details such as who created the Art work if item is a Painting
	ItemDate       string
	ItemType       string
	ItemSubject    string
	ItemMedia      string
	ItemSize       string
	ItemPicFN      string
	ItemImage      []byte // This has to be generated AES encrypted using the file name
	AES_Key        []byte // This is generated by the AES Algorithms
	ItemImageType  string // should be used to regenerate the appropriate image type
	ItemBasePrice  string // Reserve Price at Auction must be greater than this price
	CurrentOwnerID string // This is validated for a user registered record
}

////////////////////////////////////////////////////////////////////////////////
// Has an item entry every time the item changes hands
////////////////////////////////////////////////////////////////////////////////
type ItemLog struct {
	ItemID       string // PRIMARY KEY
	Status       string // SECONDARY KEY - OnAuc, OnSale, NA
	AuctionedBy  string // SECONDARY KEY - Auction House ID if applicable
	RecType      string // ITEMHIS
	ItemDesc     string
	CurrentOwner string
	Date         string // Date when status changed
}

/////////////////////////////////////////////////////////////
// Create Buyer, Seller , Auction House, Authenticator
// Could establish valid UserTypes -
// AH (Auction House)
// TR (Buyer or Seller)
// AP (Appraiser)
// IN (Insurance)
// BK (bank)
// SH (Shipper)
/////////////////////////////////////////////////////////////
type UserObject struct {
	UserID    string
	RecType   string // Type = USER
	Name      string
	UserType  string // Auction House (AH), Bank (BK), Buyer or Seller (TR), Shipper (SH), Appraiser (AP)
	Address   string
	Phone     string
	Email     string
	Bank      string
	AccountNo string
	RoutingNo string
}

/////////////////////////////////////////////////////////////////////////////
// Register a request for participating in an auction
// Usually posted by a seller who owns a piece of ITEM
// The Auction house will determine when to open the item for Auction
// The Auction House may conduct an appraisal and genuineness of the item
/////////////////////////////////////////////////////////////////////////////

type AuctionRequest struct {
	AuctionID      string
	RecType        string // AUCREQ
	ItemID         string
	AuctionHouseID string // ID of the Auction House managing the auction
	SellerID       string // ID Of Seller - to verified against the Item CurrentOwnerId
	RequestDate    string // Date on which Auction Request was filed
	ReservePrice   string // reserver price > previous purchase price
	BuyItNowPrice  string // 0 (Zero) if not applicable else specify price
	Status         string // INIT, OPEN, CLOSED (To be Updated by Trgger Auction)
	OpenDate       string // Date on which auction will occur (To be Updated by Trigger Auction)
	CloseDate      string // Date and time when Auction will close (To be Updated by Trigger Auction)
}

/////////////////////////////////////////////////////////////
// POST the transaction after the Auction Completes
// Post an Auction Transaction
// Post an Updated Item Object
// Once an auction request is opened for auctions, a timer is kicked
// off and bids are accepted. When the timer expires, the highest bid
// is selected and converted into a Transaction
// This transaction is a simple view
/////////////////////////////////////////////////////////////

type ItemTransaction struct {
	AuctionID   string
	RecType     string // POSTTRAN
	ItemID      string
	TransType   string // Sale, Buy, Commission
	UserId      string // Buyer or Seller ID
	TransDate   string // Date of Settlement (Buyer or Seller)
	HammerTime  string // Time of hammer strike - SOLD
	HammerPrice string // Total Settlement price
	Details     string // Details about the Transaction
}

////////////////////////////////////////////////////////////////
//  This is a Bid. Bids are accepted only if an auction is OPEN
////////////////////////////////////////////////////////////////

type Bid struct {
	AuctionID string
	RecType   string // BID
	BidNo     string
	ItemID    string
	BuyerID   string // ID Of Buyer - to be verified against the Item CurrentOwnerId
	BidPrice  string // BidPrice > Previous Bid
	BidTime   string // Time the bid was received
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// A Map that holds TableNames and the number of Keys
// This information is used to dynamically Create, Update
// Replace , and Query the Ledger
// In this model all attributes in a table are strings
// The chain code does both validation
// A dummy key like 2016 in some cases is used for a query to get all rows
//
//              "UserTable":        1, Key: UserID
//              "ItemTable":        1, Key: ItemID
//              "UserCatTable":     3, Key: "2016", UserType, UserID
//              "ItemCatTable":     3, Key: "2016", ItemSubject, ItemID
//              "AuctionTable":     1, Key: AuctionID
//              "AucInitTable":     2, Key: Year, AuctionID
//              "AucOpenTable":     2, Key: Year, AuctionID
//              "TransTable":       2, Key: AuctionID, ItemID
//              "BidTable":         2, Key: AuctionID, BidNo
//              "ItemHistoryTable": 4, Key: ItemID, Status, AuctionHouseID(if applicable),date-time
//
/////////////////////////////////////////////////////////////////////////////////////////////////////

func GetNumberOfKeys(tname string) int {
	TableMap := map[string]int{
		"UserTable":        1,
		"ItemTable":        1,
		"UserCatTable":     3,
		"ItemCatTable":     3,
		"AuctionTable":     1,
		"AucInitTable":     2,
		"AucOpenTable":     2,
		"TransTable":       2,
		"BidTable":         2,
		"ItemHistoryTable": 4,
	}
	return TableMap[tname]
}

//////////////////////////////////////////////////////////////
// Invoke Functions based on Function name
// The function name gets resolved to one of the following calls
// during an invoke
//
//////////////////////////////////////////////////////////////
func InvokeFunction(fname string) func(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	InvokeFunc := map[string]func(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){
		"PostItem":           PostItem,
		"PostUser":           PostUser,
		"PostAuctionRequest": PostAuctionRequest,
		"PostTransaction":    PostTransaction,
		"PostBid":            PostBid,
		"OpenAuctionForBids": OpenAuctionForBids,
		"BuyItNow":           BuyItNow,
		"TransferItem":       TransferItem,
		"CloseAuction":       CloseAuction,
		"CloseOpenAuctions":  CloseOpenAuctions,
	}
	return InvokeFunc[fname]
}

//////////////////////////////////////////////////////////////
// Query Functions based on Function name
//
//////////////////////////////////////////////////////////////
func QueryFunction(fname string) func(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	QueryFunc := map[string]func(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){
		"GetItem":               GetItem,
		"GetUser":               GetUser,
		"GetAuctionRequest":     GetAuctionRequest,
		"GetTransaction":        GetTransaction,
		"GetBid":                GetBid,
		"GetLastBid":            GetLastBid,
		"GetHighestBid":         GetHighestBid,
		"GetNoOfBidsReceived":   GetNoOfBidsReceived,
		"GetListOfBids":         GetListOfBids,
		"GetItemLog":            GetItemLog,
		"GetItemListByCat":      GetItemListByCat,
		"GetUserListByCat":      GetUserListByCat,
		"GetListOfInitAucs":     GetListOfInitAucs,
		"GetListOfOpenAucs":     GetListOfOpenAucs,
		"ValidateItemOwnership": ValidateItemOwnership,
		"IsItemOnAuction":       IsItemOnAuction,
		"GetVersion":            GetVersion,
	}
	return QueryFunc[fname]
}

//var myLogger = logging.MustGetLogger("auction_trading")

type SimpleChaincode struct {
}

var gopath string
var ccPath string

////////////////////////////////////////////////////////////////////////////////
// Chain Code Kick-off Main function
////////////////////////////////////////////////////////////////////////////////
func main() {

	// maximize CPU usage for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Starting Item Auction Application chaincode BlueMix ver 0.25 Dated 2016-07-17 15.20.00 ")

	gopath = os.Getenv("GOPATH")
	if len(os.Args) == 2 && strings.EqualFold(os.Args[1], "DEV") {
		fmt.Println("----------------- STARTED IN DEV MODE -------------------- ")
		//set chaincode path for DEV MODE
		ccPath = fmt.Sprintf("%s/src/github.com/hyperledger/fabric/auction/art/artchaincode/", gopath)
	} else {
		fmt.Println("----------------- STARTED IN NET MODE -------------------- ")
		//set chaincode path for NET MODE
		ccPath = fmt.Sprintf("%s/src/github.com/ITPeople-Blockchain/auction/art/artchaincode/", gopath)
	}

	// Start the shim -- running the fabric
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("Error starting Item Fun Application chaincode: %s", err)
	}

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// SimpleChaincode - Init Chaincode implementation - The following sequence of transactions can be used to test the Chaincode
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// TODO - Include all initialization to be complete before Invoke and Query
	// Uses aucTables to delete tables if they exist and re-create them

	//myLogger.Info("[Trade and Auction Application] Init")
	fmt.Println("[Trade and Auction Application] Init")
	var err error

	for _, val := range aucTables {
		err = stub.DeleteTable(val)
		if err != nil {
			return nil, fmt.Errorf("Init(): DeleteTable of %s  Failed ", val)
		}
		err = InitLedger(stub, val)
		if err != nil {
			return nil, fmt.Errorf("Init(): InitLedger of %s  Failed ", val)
		}
	}
	// Update the ledger with the Application version
	err = stub.PutState("version", []byte(strconv.Itoa(23)))
	if err != nil {
		return nil, err
	}

	fmt.Println("Init() Initialization Complete  : ", args)
	return []byte("Init(): Initialization Complete"), nil
}

////////////////////////////////////////////////////////////////
// SimpleChaincode - INVOKE Chaincode implementation
// User Can Invoke
// - Register a user using PostUser
// - Register an item using PostItem
// - The Owner of the item (User) can request that the item be put on auction using PostAuctionRequest
// - The Auction House can request that the auction request be Opened for bids using OpenAuctionForBids
// - One the auction is OPEN, registered buyers (Buyers) can send in bids vis PostBid
// - No bid is accepted when the status of the auction request is INIT or CLOSED
// - Either manually or by OpenAuctionRequest, the auction can be closed using CloseAuction
// - The CloseAuction creates a transaction and invokes PostTransaction
////////////////////////////////////////////////////////////////

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	var buff []byte

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Check Type of Transaction and apply business rules
	// before adding record to the block chain
	// In this version, the assumption is that args[1] specifies recType for all defined structs
	// Newer structs - the recType can be positioned anywhere and ChkReqType will check for recType
	// example:
	// ./peer chaincode invoke -l golang -n mycc -c '{"Function": "PostBid", "Args":["1111", "BID", "1", "1000", "300", "1200"]}'
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	if ChkReqType(args) == true {

		InvokeRequest := InvokeFunction(function)
		if InvokeRequest != nil {
			buff, err = InvokeRequest(stub, function, args)
		}
	} else {
		fmt.Println("Invoke() Invalid recType : ", args, "\n")
		return nil, errors.New("Invoke() : Invalid recType : " + args[0])
	}

	return buff, err
}

//////////////////////////////////////////////////////////////////////////////////////////
// SimpleChaincode - QUERY Chaincode implementation
// Client Can Query
// Sample Data
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetUser", "Args": ["4000"]}'
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetItem", "Args": ["2000"]}'
//////////////////////////////////////////////////////////////////////////////////////////

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	var buff []byte
	fmt.Println("ID Extracted and Type = ", args[0])
	fmt.Println("Args supplied : ", args)

	if len(args) < 1 {
		fmt.Println("Query() : Include at least 1 arguments Key ")
		return nil, errors.New("Query() : Expecting Transation type and Key value for query")
	}

	QueryRequest := QueryFunction(function)
	if QueryRequest != nil {
		buff, err = QueryRequest(stub, function, args)
	} else {
		fmt.Println("Query() Invalid function call : ", function)
		return nil, errors.New("Query() : Invalid function call : " + function)
	}

	if err != nil {
		fmt.Println("Query() Object not found : ", args[0])
		return nil, errors.New("Query() : Object not found : " + args[0])
	}
	return buff, err
}

//////////////////////////////////////////////////////////////////////////////////////////
// Retrieve Auction applications version Information
// This API is to check whether application has been deployed successfully or not
// example:
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetVersion", "Args": ["version"]}'
//
//////////////////////////////////////////////////////////////////////////////////////////
func GetVersion(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) < 1 {
		fmt.Println("GetVersion() : Requires 1 argument 'version'")
		return nil, errors.New("GetVersion() : Requires 1 argument 'version'")
	}
	// Get version from the ledger
	version, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for version\"}"
		return nil, errors.New(jsonResp)
	}

	if version == nil {
		jsonResp := "{\"Error\":\" auction application version is invalid\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"version\":\"" + string(version) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return version, nil
}

//////////////////////////////////////////////////////////////////////////////////////////
// Retrieve User Information
// example:
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetUser", "Args": ["100"]}'
//
//////////////////////////////////////////////////////////////////////////////////////////
func GetUser(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	// Get the Object and Display it
	Avalbytes, err := QueryLedger(stub, "UserTable", args)
	if err != nil {
		fmt.Println("GetUser() : Failed to Query Object ")
		jsonResp := "{\"Error\":\"Failed to get  Object Data for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		fmt.Println("GetUser() : Incomplete Query Object ")
		jsonResp := "{\"Error\":\"Incomplete information about the key for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Println("GetUser() : Response : Successfull -")
	return Avalbytes, nil
}

/////////////////////////////////////////////////////////////////////////////////////////
// Query callback representing the query of a chaincode
// Retrieve a Item by Item ID
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetItem", "Args": ["1000"]}'
/////////////////////////////////////////////////////////////////////////////////////////
func GetItem(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	// Get the Objects and Display it
	Avalbytes, err := QueryLedger(stub, "ItemTable", args)
	if err != nil {
		fmt.Println("GetItem() : Failed to Query Object ")
		jsonResp := "{\"Error\":\"Failed to get  Object Data for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		fmt.Println("GetItem() : Incomplete Query Object ")
		jsonResp := "{\"Error\":\"Incomplete information about the key for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Println("GetItem() : Response : Successfull ")

	// Masking ItemImage binary data
	itemObj, _ := JSONtoAR(Avalbytes)
	itemObj.ItemImage = []byte{}
	Avalbytes, _ = ARtoJSON(itemObj)

	return Avalbytes, nil
}

/////////////////////////////////////////////////////////////////////////////////////////
// Validates The Ownership of an Asset using ItemID, OwnerID, and HashKey
//
// ./peer chaincode query -l golang -n mycc -c '{"Function": "ValidateItemOwnership", "Args": ["1000", "100", "tGEBaZuKUBmwTjzNEyd+nr/fPUASuVJAZ1u7gha5fJg="]}'
//
/////////////////////////////////////////////////////////////////////////////////////////
func ValidateItemOwnership(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	if len(args) < 3 {
		fmt.Println("ValidateItemOwnership() : Requires 3 arguments Item#, Owner# and Key ")
		return nil, errors.New("ValidateItemOwnership() : Requires 3 arguments Item#, Owner# and Key")
	}

	// Get the Object Information
	Avalbytes, err := QueryLedger(stub, "ItemTable", []string{args[0]})
	if err != nil {
		fmt.Println("ValidateItemOwnership() : Failed to Query Object ")
		jsonResp := "{\"Error\":\"Failed to get  Object Data for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		fmt.Println("ValidateItemOwnership() : Incomplete Query Object ")
		jsonResp := "{\"Error\":\"Incomplete information about the key for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	myItem, err := JSONtoAR(Avalbytes)
	if err != nil {
		fmt.Println("ValidateItemOwnership() : Failed to Query Object ")
		jsonResp := "{\"Error\":\"Failed to get  Object Data for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	myKey := GetKeyValue(Avalbytes, "AES_Key")
	fmt.Println("Key String := ", myKey)

	if myKey != args[2] {
		fmt.Println("ValidateItemOwnership() : Key does not match supplied key ", args[2], " - ", myKey)
		jsonResp := "{\"Error\":\"ValidateItemOwnership() : Key does not match asset owner supplied key  " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	if myItem.CurrentOwnerID != args[1] {
		fmt.Println("ValidateItemOwnership() : ValidateItemOwnership() : Owner-Id does not match supplied ID ", args[1])
		jsonResp := "{\"Error\":\"ValidateItemOwnership() : Owner-Id does not match supplied ID " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Print("ValidateItemOwnership() : Response : Successfull - \n")
	return Avalbytes, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// Retrieve Auction Information
// This query runs against the AuctionTable
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetAuctionRequest", "Args": ["1111"]}'
// There are two other tables just for query purposes - AucInitTable, AucOpenTable
//
/////////////////////////////////////////////////////////////////////////////////////////////////////
func GetAuctionRequest(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	// Get the Objects and Display it
	Avalbytes, err := QueryLedger(stub, "AuctionTable", args)
	if err != nil {
		fmt.Println("GetAuctionRequest() : Failed to Query Object ")
		jsonResp := "{\"Error\":\"Failed to get  Object Data for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		fmt.Println("GetAuctionRequest() : Incomplete Query Object ")
		jsonResp := "{\"Error\":\"Incomplete information about the key for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Println("GetAuctionRequest() : Response : Successfull - \n")
	return Avalbytes, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// Retrieve a Bid based on two keys - AucID, BidNo
// A Bid has two Keys - The Auction Request Number and Bid Number
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetLastBid", "Args": ["1111"], "1"}'
//
///////////////////////////////////////////////////////////////////////////////////////////////////
func GetBid(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	// Check there are 2 Arguments provided as per the the struct - two are computed
	// See example
	if len(args) < 2 {
		fmt.Println("GetBid(): Incorrect number of arguments. Expecting 2 ")
		fmt.Println("GetBid(): ./peer chaincode query -l golang -n mycc -c '{\"Function\": \"GetBid\", \"Args\": [\"1111\",\"6\"]}'")
		return nil, errors.New("GetBid(): Incorrect number of arguments. Expecting 2 ")
	}

	// Get the Objects and Display it
	Avalbytes, err := QueryLedger(stub, "BidTable", args)
	if err != nil {
		fmt.Println("GetBid() : Failed to Query Object ")
		jsonResp := "{\"Error\":\"Failed to get  Object Data for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		fmt.Println("GetBid() : Incomplete Query Object ")
		jsonResp := "{\"Error\":\"Incomplete information about the key for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Println("GetBid() : Response : Successfull -")
	return Avalbytes, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// Retrieve Auction Closeout Information. When an Auction closes
// The highest bid is retrieved and converted to a Transaction
//  ./peer chaincode query -l golang -n mycc -c '{"Function": "GetTransaction", "Args": ["1111"]}'
//
///////////////////////////////////////////////////////////////////////////////////////////////////
func GetTransaction(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	//var err error

	// Get the Objects and Display it
	Avalbytes, err := QueryLedger(stub, "TransTable", args)
	if Avalbytes == nil {
		fmt.Println("GetTransaction() : Incomplete Query Object ")
		jsonResp := "{\"Error\":\"Incomplete information about the key for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	if err != nil {
		fmt.Println("GetTransaction() : Failed to Query Object ")
		jsonResp := "{\"Error\":\"Failed to get  Object Data for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Println("GetTransaction() : Response : Successfull")
	return Avalbytes, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Create a User Object. The first step is to have users
// registered
// There are different types of users - Traders (TRD), Auction Houses (AH)
// Shippers (SHP), Insurance Companies (INS), Banks (BNK)
// While this version of the chain code does not enforce strict validation
// the business process recomends validating each persona for the service
// they provide or their participation on the auction blockchain, future enhancements will do that
// ./peer chaincode invoke -l golang -n mycc -c '{"Function": "PostUser", "Args":["100", "USER", "Ashley Hart", "TRD",  "Morrisville Parkway, #216, Morrisville, NC 27560", "9198063535", "ashley@itpeople.com", "SUNTRUST", "00017102345", "0234678"]}'
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func PostUser(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	record, err := CreateUserObject(args[0:]) //
	if err != nil {
		return nil, err
	}
	buff, err := UsertoJSON(record) //

	if err != nil {
		fmt.Println("PostuserObject() : Failed Cannot create object buffer for write : ", args[1])
		return nil, errors.New("PostUser(): Failed Cannot create object buffer for write : " + args[1])
	} else {
		// Update the ledger with the Buffer Data
		// err = stub.PutState(args[0], buff)
		keys := []string{args[0]}
		err = UpdateLedger(stub, "UserTable", keys, buff)
		if err != nil {
			fmt.Println("PostUser() : write error while inserting record")
			return nil, err
		}

		// Post Entry into UserCatTable - i.e. User Category Table
		keys = []string{"2016", args[3], args[0]}
		err = UpdateLedger(stub, "UserCatTable", keys, buff)
		if err != nil {
			fmt.Println("PostUser() : write error while inserting recordinto UserCatTable \n")
			return nil, err
		}
	}

	return buff, err
}

func CreateUserObject(args []string) (UserObject, error) {

	var err error
	var aUser UserObject

	// Check there are 10 Arguments
	if len(args) != 10 {
		fmt.Println("CreateUserObject(): Incorrect number of arguments. Expecting 10 ")
		return aUser, errors.New("CreateUserObject() : Incorrect number of arguments. Expecting 10 ")
	}

	// Validate UserID is an integer

	_, err = strconv.Atoi(args[0])
	if err != nil {
		return aUser, errors.New("CreateUserObject() : User ID should be an integer")
	}

	aUser = UserObject{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9]}
	fmt.Println("CreateUserObject() : User Object : ", aUser)

	return aUser, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Create a master Object of the Item
// Since the Owner Changes hands, a record has to be written for each
// Transaction with the updated Encryption Key of the new owner
// Example
//./peer chaincode invoke -l golang -n mycc -c '{"Function": "PostItem", "Args":["1000", "ARTINV", "Shadows by Asppen", "Asppen Messer", "20140202", "Original", "Landscape" , "Canvas", "15 x 15 in", "sample_7.png","$600", "100"]}'
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

func PostItem(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	itemObject, err := CreateItemObject(args[0:])
	if err != nil {
		fmt.Println("PostItem(): Cannot create item object \n")
		return nil, err
	}

	// Check if the Owner ID specified is registered and valid
	ownerInfo, err := ValidateMember(stub, itemObject.CurrentOwnerID)
	fmt.Println("Owner information  ", ownerInfo, itemObject.CurrentOwnerID)
	if err != nil {
		fmt.Println("PostItem() : Failed Owner information not found for ", itemObject.CurrentOwnerID)
		return nil, err
	}

	// Convert Item Object to JSON
	buff, err := ARtoJSON(itemObject) //
	if err != nil {
		fmt.Println("PostItem() : Failed Cannot create object buffer for write : ", args[1])
		return nil, errors.New("PostItem(): Failed Cannot create object buffer for write : " + args[1])
	} else {
		// Update the ledger with the Buffer Data
		// err = stub.PutState(args[0], buff)
		keys := []string{args[0]}
		err = UpdateLedger(stub, "ItemTable", keys, buff)
		if err != nil {
			fmt.Println("PostItem() : write error while inserting record\n")
			return buff, err
		}

		// Put an entry into the Item History Table
		_, err = PostItemLog(stub, itemObject, "INITIAL", "DEFAULT")
		if err != nil {
			fmt.Println("PostItemLog() : write error while inserting record\n")
			return nil, err
		}

		// Post Entry into ItemCatTable - i.e. Item Category Table
		// The first key 2016 is a dummy (band aid) key to extract all values
		keys = []string{"2016", args[6], args[0]}
		err = UpdateLedger(stub, "ItemCatTable", keys, buff)
		if err != nil {
			fmt.Println("PostItem() : Write error while inserting record into ItemCatTable \n")
			return buff, err
		}

	}

	secret_key, _ := json.Marshal(itemObject.AES_Key)
	fmt.Println(string(secret_key))
	return secret_key, nil
}

func CreateItemObject(args []string) (ItemObject, error) {

	var err error
	var myItem ItemObject

	// Check there are 12 Arguments provided as per the the struct - two are computed
	if len(args) != 12 {
		fmt.Println("CreateItemObject(): Incorrect number of arguments. Expecting 12 ")
		return myItem, errors.New("CreateItemObject(): Incorrect number of arguments. Expecting 12 ")
	}

	// Validate ItemID is an integer

	_, err = strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("CreateItemObject(): ART ID should be an integer create failed! ")
		return myItem, errors.New("CreateItemObject(): ART ID should be an integer create failed!")
	}

	// Validate Picture File exists based on the name provided
	// Looks for file in current directory of application and must be fixed for other locations

	// Validate Picture File exists based on the name provided
	// Looks for file in current directory of application and must be fixed for other locations
	imagePath := ccPath + args[9]
	if _, err := os.Stat(imagePath); err == nil {
		fmt.Println(imagePath, "  exists!")
	} else {
		fmt.Println("CreateItemObject(): Cannot find or load Picture File = %s :  %s\n", imagePath, err)
		return myItem, errors.New("CreateItemObject(): ART Picture File not found " + imagePath)
	}

	// Get the Item Image and convert it to a byte array
	imagebytes, fileType := imageToByteArray(imagePath)

	// Generate a new key and encrypt the image

	AES_key, _ := GenAESKey()
	AES_enc := Encrypt(AES_key, imagebytes)

	// Append the AES Key, The Encrypted Image Byte Array and the file type
	myItem = ItemObject{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], AES_enc, AES_key, fileType, args[10], args[11]}

	fmt.Println("CreateItemObject(): Item Object created: ID# ", myItem.ItemID, "\n AES Key: ", myItem.AES_Key)

	// Code to Validate the Item Object)
	// If User presents Crypto Key then key is used to validate the picture that is stored as part of the title
	// TODO

	return myItem, nil
}

///////////////////////////////////////////////////////////////////////////////////
// Since the Owner Changes hands, a record has to be written for each
// Transaction with the updated Encryption Key of the new owner
// This function is internally invoked by PostTransaction and is not a Public API
///////////////////////////////////////////////////////////////////////////////////

func UpdateItemObject(stub shim.ChaincodeStubInterface, ar []byte, hammerPrice string, buyer string) ([]byte, error) {

	var err error
	myItem, err := JSONtoAR(ar)
	if err != nil {
		fmt.Println("U() : UpdateItemObject() : Failed to create Art Record Object from JSON ")
		return nil, err
	}

	// Insert logic to  re-encrypt image by first fetching the current Key
	CurrentAES_Key := myItem.AES_Key
	// Decrypt Image and Save Image in a file
	image := Decrypt(CurrentAES_Key, myItem.ItemImage)

	// Get a New Key & Encrypt Image with New Key
	myItem.AES_Key, _ = GenAESKey()
	myItem.ItemImage = Encrypt(myItem.AES_Key, image)

	// Update the owner to the Buyer and update price to auction hammer price
	myItem.ItemBasePrice = hammerPrice
	myItem.CurrentOwnerID = buyer

	ar, err = ARtoJSON(myItem)
	keys := []string{myItem.ItemID, myItem.CurrentOwnerID}
	err = ReplaceLedgerEntry(stub, "ItemTable", keys, ar)
	if err != nil {
		fmt.Println("UpdateItemObject() : Failed ReplaceLedgerEntry in ItemTable into Blockchain ")
		return nil, err
	}
	fmt.Println("UpdateItemObject() : ReplaceLedgerEntry in ItemTable successful ")

	// Update entry in Item Category Table as it holds the Item object as wekk
	keys = []string{"2016", myItem.ItemSubject, myItem.ItemID}
	err = ReplaceLedgerEntry(stub, "ItemCatTable", keys, ar)
	if err != nil {
		fmt.Println("UpdateItemObject() : Failed ReplaceLedgerEntry in ItemCategoryTable into Blockchain ")
		return nil, err
	}

	fmt.Println("UpdateItemObject() : ReplaceLedgerEntry in ItemCategoryTable successful ")
	return myItem.AES_Key, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Obtain Asset Details and Validate Item
// Transfer Item to new owner - no change in price  - In the example XFER is the recType
// ./peer chaincode invoke -l golang -n mycc -c '{"Function": "TransferItem", "Args": ["1000", "100", "tGEBaZuKUBmwTjzNEyd+nr/fPUASuVJAZ1u7gha5fJg=", "300", "XFER"]}'
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TransferItem(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	if len(args) < 5 {
		fmt.Println("TransferItem() : Requires 5 arguments Item#, Owner#, Key#, newOwnerID#, XFER \n")
		return nil, errors.New("TransferItem() : Requires 5 arguments Item#, Owner#, Key#, newOwnerID#, XFER")
	}

	// Let us make sure that the Item is not on Auction
	err = VerifyIfItemIsOnAuction(stub, args[0])
	if err != nil {
		fmt.Println("TransferItem() : Failed Item is either initiated or opened for Auction ", args[0])
		return nil, err
	}

	// Validate New Owner's ID
	_, err = ValidateMember(stub, args[3])
	if err != nil {
		fmt.Println("TransferItem() : Failed transferee not Registered in Blockchain ", args[3])
		return nil, err
	}

	// Validate Item or Asset Ownership
	ar, err := ValidateItemOwnership(stub, "ValidateItemOwnership", args[:3])
	if err != nil {
		fmt.Println("TransferItem() : ValidateItemOwnership() : Failed to authenticate item or asset ownership")
		return nil, err
	}

	myItem, err := JSONtoAR(ar)
	if err != nil {
		fmt.Println("TransferItem() : Failed to create item Object from JSON ")
		return nil, err
	}

	// Insert logic to  re-encrypt image by first fetching the current Key
	CurrentAES_Key := myItem.AES_Key
	// Decrypt Image and Save Image in a file
	image := Decrypt(CurrentAES_Key, myItem.ItemImage)

	// Get a New Key & Encrypt Image with New Key
	myItem.AES_Key, _ = GenAESKey()
	myItem.ItemImage = Encrypt(myItem.AES_Key, image)

	// Update the owner to the new owner transfered to
	myItem.CurrentOwnerID = args[3]

	ar, err = ARtoJSON(myItem)
	keys := []string{myItem.ItemID, myItem.CurrentOwnerID}
	err = ReplaceLedgerEntry(stub, "ItemTable", keys, ar)
	if err != nil {
		fmt.Println("TransferAsset() : Failed ReplaceLedgerEntry in ItemTable into Blockchain ")
		return nil, err
	}
	fmt.Println("TransferAsset() : ReplaceLedgerEntry in ItemTable successful ")

	// Update entry in Item Category Table as it holds the Item object as well
	keys = []string{"2016", myItem.ItemSubject, myItem.ItemID}
	err = ReplaceLedgerEntry(stub, "ItemCatTable", keys, ar)
	if err != nil {
		fmt.Println("TransferAsset() : Failed ReplaceLedgerEntry in ItemCategoryTable into Blockchain ")
		return nil, err
	}

	_, err = PostItemLog(stub, myItem, "Transfer", args[1])
	if err != nil {
		fmt.Println("TransferItem() : PostItemLog() write error while inserting record\n")
		return nil, err
	}

	fmt.Println("TransferAsset() : ReplaceLedgerEntry in ItemCategoryTable successful ")
	return myItem.AES_Key, nil
}

////////////////////////////////////////////////////////////////////////////////////
// Validate Item Status - Is it currently on Auction, if so Reject Transfer Request
// This can be written better - will do so if things work
// The function return the Auction ID and the Status = OPEN or INIT
////////////////////////////////////////////////////////////////////////////////////

func VerifyIfItemIsOnAuction(stub shim.ChaincodeStubInterface, itemID string) error {

	rows, err := GetListOfOpenAucs(stub, "AucOpenTable", []string{"2016"})
	if err != nil {
		return fmt.Errorf("VerifyIfItemIsOnAuction() operation failed. Error retrieving values from AucOpenTable: %s", err)
	}

	tlist := make([]AuctionRequest, len(rows))
	err = json.Unmarshal([]byte(rows), &tlist)
	if err != nil {
		fmt.Println("VerifyIfItemIsOnAuction: Unmarshal failed : ", err)
		return fmt.Errorf("VerifyIfItemIsOnAuction: operation failed. Error un-marshaling JSON: %s", err)
	}

	for i := 0; i < len(tlist); i++ {
		ar := tlist[i]

		// Compare Auction IDs
		if ar.ItemID == itemID {
			fmt.Println("VerifyIfItemIsOnAuction() : Item Exists")
			return fmt.Errorf("VerifyIfItemIsOnAuction() operation failed. %s", itemID)
		}
	}

	// Now Check if an Auction Has been inititiated
	// If so , it has to be removed from Auction for a Transfer

	rows, err = GetListOfInitAucs(stub, "AucInitTable", []string{"2016"})
	if err != nil {
		return fmt.Errorf("VerifyIfItemIsOnAuction() operation failed. Error retrieving values from AucInitTable: %s", err)
	}

	tlist = make([]AuctionRequest, len(rows))
	err = json.Unmarshal([]byte(rows), &tlist)
	if err != nil {
		fmt.Println("VerifyIfItemIsOnAuction() Unmarshal failed : ", err)
		return fmt.Errorf("VerifyIfItemIsOnAuction: operation failed. Error un-marshaling JSON: %s", err)
	}

	for i := 0; i < len(tlist); i++ {
		ar := tlist[i]
		if err != nil {
			fmt.Println("VerifyIfItemIsOnAuction() Failed : Ummarshall error")
			return fmt.Errorf("VerifyIfItemIsOnAuction() operation failed. %s", err)
		}

		// Compare Auction IDs
		if ar.ItemID == itemID {
			return fmt.Errorf("VerifyIfItemIsOnAuction() operation failed.")
		}
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Checks if an Item is available on auction or not
// Input ItemID # 1000
// See Sample data
// ./peer chaincode query -l golang -n mycc -c '{"Function": "IsItemOnAuction", "Args": ["1000", "VERIFY"]}'
/////////////////////////////////////////////////////////////////////////////////////////////////////////////
func IsItemOnAuction(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) < 2 {
		fmt.Println("IsItemOnAuction() : Requires 2 arguments Item#, RecordType")
		return nil, errors.New("IsItemOnAuction() : Requires 2 arguments Item#, RecordType")
	}

	itemExists := false
	err := VerifyIfItemIsOnAuction(stub, args[0])
	if err != nil {
		fmt.Println("IsItemOnAuction() : Failed Item# ", args[0], " is either initiated or opened for Auction ")
		itemExists = true
	}
	fmt.Println("Is Item# ", args[0], " on-auction : ", itemExists)
	ie, _ := json.Marshal(itemExists)
	return ie, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// POSTS A LOG ENTRY Every Time the Item is transacted
// Valid Status for ItemLog =  OnAuc, OnSale, NA, INITIAL
// Valid AuctionedBy: This value is set to "DEFAULT" but when it is put on auction Auction House ID is assigned
// PostItemLog IS NOT A PUBLIC API and is invoked every time some event happens in the Item's life
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func PostItemLog(stub shim.ChaincodeStubInterface, item ItemObject, status string, ah string) ([]byte, error) {

	iLog := ItemToItemLog(item)
	iLog.Status = status
	iLog.AuctionedBy = ah

	buff, err := ItemLogtoJSON(iLog)
	if err != nil {
		fmt.Println("PostItemLog() : Failed Cannot create object buffer for write : ", item.ItemID)
		return nil, errors.New("PostItemLog(): Failed Cannot create object buffer for write : " + item.ItemID)
	} else {
		// Update the ledger with the Buffer Data
		keys := []string{iLog.ItemID, iLog.Status, iLog.AuctionedBy, time.Now().Format("2006-01-02 15:04:05")}
		err = UpdateLedger(stub, "ItemHistoryTable", keys, buff)
		if err != nil {
			fmt.Println("PostItemLog() : write error while inserting record\n")
			return buff, err
		}
	}
	return buff, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Create an Auction Request
// The owner of an Item, when ready to put the item on an auction
// will create an auction request  and specify a  auction house.
//
// ./peer chaincode invoke -l golang -n mycc -c '{"Function": "PostAuctionRequest", "Args":["1111", "AUCREQ", "1700", "200", "400", "04012016", "1200", "INIT", "2016-05-20 11:00:00.3 +0000 UTC","2016-05-23 11:00:00.3 +0000 UTC"]}'
//
// The start and end time of the auction are actually assigned when the auction is opened  by OpenAuctionForBids()
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func PostAuctionRequest(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	ar, err := CreateAuctionRequest(args[0:])
	if err != nil {
		return nil, err
	}

	// Let us make sure that the Item is not on Auction
	err = VerifyIfItemIsOnAuction(stub, ar.ItemID)
	if err != nil {
		fmt.Println("PostAuctionRequest() : Failed Item is either initiated or opened for Auction ", args[0])
		return nil, err
	}

	// Validate Auction House to check it is a registered User
	aucHouse, err := ValidateMember(stub, ar.AuctionHouseID)
	fmt.Println("Auction House information  ", aucHouse, " ID: ", ar.AuctionHouseID)
	if err != nil {
		fmt.Println("PostAuctionRequest() : Failed Auction House not Registered in Blockchain ", ar.AuctionHouseID)
		return nil, err
	}

	// Validate Item record
	itemObject, err := ValidateItemSubmission(stub, ar.ItemID)
	if err != nil {
		fmt.Println("PostAuctionRequest() : Failed Could not Validate Item Object in Blockchain ", ar.ItemID)
		return itemObject, err
	}

	// Convert AuctionRequest to JSON
	buff, err := AucReqtoJSON(ar) // Converting the auction request struct to []byte array
	if err != nil {
		fmt.Println("PostAuctionRequest() : Failed Cannot create object buffer for write : ", args[1])
		return nil, errors.New("PostAuctionRequest(): Failed Cannot create object buffer for write : " + args[1])
	} else {
		// Update the ledger with the Buffer Data
		//err = stub.PutState(args[0], buff)
		keys := []string{args[0]}
		err = UpdateLedger(stub, "AuctionTable", keys, buff)
		if err != nil {
			fmt.Println("PostAuctionRequest() : write error while inserting record\n")
			return buff, err
		}

		// Post an Item Log and the Auction House ID is included in the log
		// Recall -- that by default that value is "DEFAULT"
		io, err := JSONtoAR(itemObject)
		_, err = PostItemLog(stub, io, "ReadyForAuc", ar.AuctionHouseID)
		if err != nil {
			fmt.Println("PostItemLog() : write error while inserting record\n")
			return buff, err
		}

		//An entry is made in the AuctionInitTable that this Item has been placed for Auction
		// The UI can pull all items available for auction and the item can be Opened for accepting bids
		// The 2016 is a dummy key and has notr value other than to get all rows

		keys = []string{"2016", args[0]}
		err = UpdateLedger(stub, "AucInitTable", keys, buff)
		if err != nil {
			fmt.Println("PostAuctionRequest() : write error while inserting record into AucInitTable \n")
			return buff, err
		}

	}

	return buff, err
}

func CreateAuctionRequest(args []string) (AuctionRequest, error) {
	var err error
	var aucReg AuctionRequest

	// Check there are 11 Arguments
	// See example -- The Open and Close Dates are Dummy, and will be set by open auction
	// '{"Function": "PostAuctionRequest", "Args":["1111", "AUCREQ", "1000", "200", "100", "04012016", "1200", "1800",
	//   "INIT", "2016-05-20 11:00:00.3 +0000 UTC","2016-05-23 11:00:00.3 +0000 UTC"]}'
	if len(args) != 11 {
		fmt.Println("CreateAuctionRegistrationObject(): Incorrect number of arguments. Expecting 11 ")
		return aucReg, errors.New("CreateAuctionRegistrationObject() : Incorrect number of arguments. Expecting 11 ")
	}

	// Validate UserID is an integer . I think this redundant and can be avoided

	err = validateID(args[0])
	if err != nil {
		return aucReg, errors.New("CreateAuctionRequest() : User ID should be an integer")
	}

	aucReg = AuctionRequest{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10]}
	fmt.Println("CreateAuctionObject() : Auction Registration : ", aucReg)

	return aucReg, nil
}

//////////////////////////////////////////////////////////
// Create an Item Transaction record to process Request
// This is invoked by the CloseAuctionRequest
//
////////////////////////////////////////////////////////////
func PostTransaction(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function != "PostTransaction" {
		return nil, errors.New("PostTransaction(): Invalid function name. Expecting \"PostTransaction\"")
	}

	ar, err := CreateTransactionRequest(args[0:]) //
	if err != nil {
		return nil, err
	}

	// Validate buyer's ID
	buyer, err := ValidateMember(stub, ar.UserId)
	if err != nil {
		fmt.Println("PostTransaction() : Failed Buyer not Registered in Blockchain ", ar.UserId)
		return nil, err
	}

	fmt.Println("PostTransaction(): Validated Buyer information successfully ", buyer, ar.UserId)

	// Validate Item record
	lastUpdatedItemOBCObject, err := ValidateItemSubmission(stub, ar.ItemID)
	if err != nil {
		fmt.Println("PostTransaction() : Failed Could not Validate Item Object in Blockchain ", ar.ItemID)
		return lastUpdatedItemOBCObject, err
	}
	fmt.Println("PostTransaction() : Validated Item Object in Blockchain successfully", ar.ItemID)

	// Update Item Object with new Owner Key
	newKey, err := UpdateItemObject(stub, lastUpdatedItemOBCObject, ar.HammerPrice, ar.UserId)
	if err != nil {
		fmt.Println("PostTransaction() : Failed to update Item Master Object in Blockchain ", ar.ItemID)
		return nil, err
	} else {
		// Write New Key to file
		fmt.Println("PostTransaction() : New encryption Key is  ", newKey)
	}
	fmt.Println("PostTransaction() : Updated Item Master Object in Blockchain successfully", ar.ItemID)

	// Post an Item Log
	itemObject, err := JSONtoAR(lastUpdatedItemOBCObject)
	if err != nil {
		fmt.Println("PostTransaction() : Conversion error JSON to ItemRecord\n")
		return lastUpdatedItemOBCObject, err
	}

	// A life cycle event is added to say that the Item is no longer on auction
	itemObject.ItemBasePrice = ar.HammerPrice
	itemObject.CurrentOwnerID = ar.UserId

	_, err = PostItemLog(stub, itemObject, "NA", "DEFAULT")
	if err != nil {
		fmt.Println("PostTransaction() : write error while inserting item log record\n")
		return lastUpdatedItemOBCObject, err
	}

	fmt.Println("PostTransaction() : Inserted item log record successfully", ar.ItemID)

	// Convert Transaction Object to JSON
	buff, err := TrantoJSON(ar) //
	if err != nil {
		fmt.Println("GetObjectBuffer() : Failed to convert Transaction Object to JSON ", args[0])
		return nil, err
	}

	// Update the ledger with the Buffer Data
	keys := []string{args[0], args[3]}
	err = UpdateLedger(stub, "TransTable", keys, buff)
	if err != nil {
		fmt.Println("PostTransaction() : write error while inserting record\n")
		return buff, err
	}

	fmt.Println("PostTransaction() : Posted Transaction Record successfully\n")

	// Returns New Key. To get Transaction Details, run GetTransaction

	secret_key, _ := json.Marshal(newKey)
	fmt.Println(string(secret_key))
	return secret_key, nil

}

func CreateTransactionRequest(args []string) (ItemTransaction, error) {

	var at ItemTransaction

	// Check there are 9 Arguments
	if len(args) != 9 {
		fmt.Println("CreateTransactionRequest(): Incorrect number of arguments. Expecting 9 ")
		return at, errors.New("CreateTransactionRequest() : Incorrect number of arguments. Expecting 9 ")
	}

	at = ItemTransaction{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8]}
	fmt.Println("CreateTransactionRequest() : Transaction Request: ", at)

	return at, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Create a Bid Object
// Once an Item has been opened for auction, bids can be submitted as long as the auction is "OPEN"
//./peer chaincode invoke -l golang -n mycc -c '{"Function": "PostBid", "Args":["1111", "BID", "1", "1000", "300", "1200"]}'
//./peer chaincode invoke -l golang -n mycc -c '{"Function": "PostBid", "Args":["1111", "BID", "2", "1000", "400", "3000"]}'
//
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func PostBid(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	bid, err := CreateBidObject(args[0:]) //
	if err != nil {
		return nil, err
	}

	// Reject the Bid if the Buyer Information Is not Valid or not registered on the Block Chain
	buyerInfo, err := ValidateMember(stub, args[4])
	fmt.Println("Buyer information  ", buyerInfo, "  ", args[4])
	if err != nil {
		fmt.Println("PostBid() : Failed Buyer not registered on the block-chain ", args[4])
		return nil, err
	}

	///////////////////////////////////////
	// Reject Bid if Auction is not "OPEN"
	///////////////////////////////////////
	RBytes, err := GetAuctionRequest(stub, "GetAuctionRequest", []string{args[0]})
	if err != nil {
		fmt.Println("PostBid() : Cannot find Auction record ", args[0])
		return nil, errors.New("PostBid(): Cannot find Auction record : " + args[0])
	}

	aucR, err := JSONtoAucReq(RBytes)
	if err != nil {
		fmt.Println("PostBid() : Cannot UnMarshall Auction record")
		return nil, errors.New("PostBid(): Cannot UnMarshall Auction record: " + args[0])
	}

	if aucR.Status != "OPEN" {
		fmt.Println("PostBid() : Cannot accept Bid as Auction is not OPEN ", args[0])
		return nil, errors.New("PostBid(): Cannot accept Bid as Auction is not OPEN : " + args[0])
	}

	///////////////////////////////////////////////////////////////////
	// Reject Bid if the time bid was received is > Auction Close Time
	///////////////////////////////////////////////////////////////////
	if tCompare(bid.BidTime, aucR.CloseDate) == false {
		fmt.Println("PostBid() Failed : BidTime past the Auction Close Time")
		return nil, fmt.Errorf("PostBid() Failed : BidTime past the Auction Close Time %s, %s", bid.BidTime, aucR.CloseDate)
	}

	//////////////////////////////////////////////////////////////////
	// Reject Bid if Item ID on Bid does not match Item ID on Auction
	//////////////////////////////////////////////////////////////////
	if aucR.ItemID != bid.ItemID {
		fmt.Println("PostBid() Failed : Item ID mismatch on bid. Bid Rejected")
		return nil, errors.New("PostBid() : Item ID mismatch on Bid. Bid Rejected")
	}

	//////////////////////////////////////////////////////////////////////
	// Reject Bid if Bid Price is less than Reserve Price
	// Convert Bid Price and Reserve Price to Integer (TODO - Float)
	//////////////////////////////////////////////////////////////////////
	bp, err := strconv.Atoi(bid.BidPrice)
	if err != nil {
		fmt.Println("PostBid() Failed : Bid price should be an integer")
		return nil, errors.New("PostBid() : Bid price should be an integer")
	}

	hp, err := strconv.Atoi(aucR.ReservePrice)
	if err != nil {
		return nil, errors.New("PostItem() : Reserve Price should be an integer")
	}

	// Check if Bid Price is > Auction Request Reserve Price
	if bp < hp {
		return nil, errors.New("PostItem() : Bid Price must be greater than Reserve Price")
	}

	////////////////////////////
	// Post or Accept the Bid
	////////////////////////////
	buff, err := BidtoJSON(bid) //

	if err != nil {
		fmt.Println("PostBid() : Failed Cannot create object buffer for write : ", args[1])
		return nil, errors.New("PostBid(): Failed Cannot create object buffer for write : " + args[1])
	} else {
		// Update the ledger with the Buffer Data
		// err = stub.PutState(args[0], buff)
		keys := []string{args[0], args[2]}
		err = UpdateLedger(stub, "BidTable", keys, buff)
		if err != nil {
			fmt.Println("PostBidTable() : write error while inserting record\n")
			return buff, err
		}
	}

	return buff, err
}

func CreateBidObject(args []string) (Bid, error) {
	var err error
	var aBid Bid

	// Check there are 11 Arguments
	// See example
	if len(args) != 6 {
		fmt.Println("CreateBidObject(): Incorrect number of arguments. Expecting 6 ")
		return aBid, errors.New("CreateBidObject() : Incorrect number of arguments. Expecting 6 ")
	}

	// Validate Bid is an integer

	_, err = strconv.Atoi(args[0])
	if err != nil {
		return aBid, errors.New("CreateBidObject() : Bid ID should be an integer")
	}

	_, err = strconv.Atoi(args[2])
	if err != nil {
		return aBid, errors.New("CreateBidObject() : Bid ID should be an integer")
	}

	bidTime := time.Now().Format("2006-01-02 15:04:05")

	aBid = Bid{args[0], args[1], args[2], args[3], args[4], args[5], bidTime}
	fmt.Println("CreateBidObject() : Bid Object : ", aBid)

	return aBid, nil
}

///////////////////////////////////////////////////////////
// Convert Image to []bytes and viceversa
// Detect Image Filetype
// Image Function to read an image and create a byte array
// Currently only PNG images are supported
///////////////////////////////////////////////////////////
func imageToByteArray(imageFile string) ([]byte, string) {

	file, err := os.Open(imageFile)

	if err != nil {
		fmt.Println("imageToByteArray() : cannot OPEN image file ", err)
		return nil, string("imageToByteArray() : cannot OPEN image file ")
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	// read file into bytes
	buff := bufio.NewReader(file)
	_, err = buff.Read(bytes)

	if err != nil {
		fmt.Println("imageToByteArray() : cannot READ image file")
		return nil, string("imageToByteArray() : cannot READ image file ")
	}

	filetype := http.DetectContentType(bytes)
	fmt.Println("imageToByteArray() : ", filetype)
	//filetype := GetImageType(bytes)

	return bytes, filetype
}

//////////////////////////////////////////////////////
// If Valid fileType, will have "image" as first word
//////////////////////////////////////////////////////
func GetImageType(buff []byte) string {
	filetype := http.DetectContentType(buff)

	switch filetype {
	case "image/jpeg", "image/jpg":
		return filetype

	case "image/gif":
		return filetype

	case "image/png":
		return filetype

	case "application/pdf": // not image, but application !
		filetype = "application/pdf"
	default:
		filetype = "Unknown"
	}
	return filetype
}

////////////////////////////////////////////////////////////
// Converts a byteArray into an image and saves it
// into an appropriate file
// It is important to get the file type before saving the
// file by call the GetImageType
////////////////////////////////////////////////////////////
func ByteArrayToImage(imgByte []byte, imageFile string) error {

	// convert []byte to image for saving to file
	img, _, _ := image.Decode(bytes.NewReader(imgByte))

	fmt.Println("ProcessQueryResult ByteArrayToImage : proceeding to create image ")

	//save the imgByte to file
	out, err := os.Create(imageFile)

	if err != nil {
		fmt.Println("ByteArrayToImage() : cannot CREATE image file ", err)
		return errors.New("ByteArrayToImage() : cannot CREATE image file ")
	}
	fmt.Println("ProcessRequestType ByteArrayToImage : proceeding to Encode image ")

	//err = png.Encode(out, img)
	filetype := http.DetectContentType(imgByte)

	switch filetype {
	case "image/jpeg", "image/jpg":
		var opt jpeg.Options
		opt.Quality = 100
		err = jpeg.Encode(out, img, &opt)

	case "image/gif":
		var opt gif.Options
		opt.NumColors = 256
		err = gif.Encode(out, img, &opt)

	case "image/png":
		err = png.Encode(out, img)

	default:
		err = errors.New("Only PMNG, JPG and GIF Supported ")
	}

	if err != nil {
		fmt.Println("ByteArrayToImage() : cannot ENCODE image file ", err)
		return errors.New("ByteArrayToImage() : cannot ENCODE image file ")
	}

	// everything ok
	fmt.Println("Image file  generated and saved to ", imageFile)
	return nil
}

///////////////////////////////////////////////////////////////////////
// Encryption and Decryption Section
// Images will be Encrypted and stored and the key will be part of the
// certificate that is provided to the Owner
///////////////////////////////////////////////////////////////////////

const (
	AESKeyLength = 32 // AESKeyLength is the default AES key length
	NonceSize    = 24 // NonceSize is the default NonceSize
)

///////////////////////////////////////////////////
// GetRandomBytes returns len random looking bytes
///////////////////////////////////////////////////
func GetRandomBytes(len int) ([]byte, error) {
	key := make([]byte, len)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

////////////////////////////////////////////////////////////
// GenAESKey returns a random AES key of length AESKeyLength
// 3 Functions to support Encryption and Decryption
// GENAESKey() - Generates AES symmetric key
// Encrypt() Encrypts a [] byte
// Decrypt() Decryts a [] byte
////////////////////////////////////////////////////////////
func GenAESKey() ([]byte, error) {
	return GetRandomBytes(AESKeyLength)
}

func PKCS5Pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	pad := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, pad...)
}

func PKCS5Unpad(src []byte) []byte {
	len := len(src)
	unpad := int(src[len-1])
	return src[:(len - unpad)]
}

func Decrypt(key []byte, ciphertext []byte) []byte {

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		panic("Text is too short")
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]

	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext
}

func Encrypt(key []byte, ba []byte) []byte {

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Empty array of 16 + ba length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(ba))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from ba to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], ba)

	return ciphertext
}

//////////////////////////////////////////////////////////
// JSON To args[] - return a map of the JSON string
//////////////////////////////////////////////////////////
func JSONtoArgs(Avalbytes []byte) (map[string]interface{}, error) {

	var data map[string]interface{}

	if err := json.Unmarshal(Avalbytes, &data); err != nil {
		return nil, err
	}

	return data, nil
}

//////////////////////////////////////////////////////////
// Variation of the above - return value from a JSON string
//////////////////////////////////////////////////////////

func GetKeyValue(Avalbytes []byte, key string) string {
	var dat map[string]interface{}
	if err := json.Unmarshal(Avalbytes, &dat); err != nil {
		panic(err)
	}

	val := dat[key].(string)
	return val
}

//////////////////////////////////////////////////////////
// Time and Date Comparison
// tCompare("2016-06-28 18:40:57", "2016-06-27 18:45:39")
//////////////////////////////////////////////////////////
func tCompare(t1 string, t2 string) bool {

	layout := "2006-01-02 15:04:05"
	bidTime, err := time.Parse(layout, t1)
	if err != nil {
		fmt.Println("tCompare() Failed : time Conversion error on t1")
		return false
	}

	aucCloseTime, err := time.Parse(layout, t2)
	if err != nil {
		fmt.Println("tCompare() Failed : time Conversion error on t2")
		return false
	}

	if bidTime.Before(aucCloseTime) {
		return true
	}

	return false
}

//////////////////////////////////////////////////////////
// Converts JSON String to an ART Object
//////////////////////////////////////////////////////////
func JSONtoAR(data []byte) (ItemObject, error) {

	ar := ItemObject{}
	err := json.Unmarshal([]byte(data), &ar)
	if err != nil {
		fmt.Println("Unmarshal failed : ", err)
	}

	return ar, err
}

//////////////////////////////////////////////////////////
// Converts an ART Object to a JSON String
//////////////////////////////////////////////////////////
func ARtoJSON(ar ItemObject) ([]byte, error) {

	ajson, err := json.Marshal(ar)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ajson, nil
}

//////////////////////////////////////////////////////////
// Converts an BID to a JSON String
//////////////////////////////////////////////////////////
func ItemLogtoJSON(item ItemLog) ([]byte, error) {

	ajson, err := json.Marshal(item)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ajson, nil
}

//////////////////////////////////////////////////////////
// Converts an User Object to a JSON String
//////////////////////////////////////////////////////////
func JSONtoItemLog(ithis []byte) (ItemLog, error) {

	item := ItemLog{}
	err := json.Unmarshal(ithis, &item)
	if err != nil {
		fmt.Println("JSONtoItemLog error: ", err)
		return item, err
	}
	return item, err
}

//////////////////////////////////////////////////////////
// Converts an Auction Request to a JSON String
//////////////////////////////////////////////////////////
func AucReqtoJSON(ar AuctionRequest) ([]byte, error) {

	ajson, err := json.Marshal(ar)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ajson, nil
}

//////////////////////////////////////////////////////////
// Converts an User Object to a JSON String
//////////////////////////////////////////////////////////
func JSONtoAucReq(areq []byte) (AuctionRequest, error) {

	ar := AuctionRequest{}
	err := json.Unmarshal(areq, &ar)
	if err != nil {
		fmt.Println("JSONtoAucReq error: ", err)
		return ar, err
	}
	return ar, err
}

//////////////////////////////////////////////////////////
// Converts BID Object to JSON String
//////////////////////////////////////////////////////////
func BidtoJSON(myHand Bid) ([]byte, error) {

	ajson, err := json.Marshal(myHand)
	if err != nil {
		fmt.Println("BidtoJSON error: ", err)
		return nil, err
	}
	return ajson, nil
}

//////////////////////////////////////////////////////////
// Converts JSON String to BID Object
//////////////////////////////////////////////////////////
func JSONtoBid(areq []byte) (Bid, error) {

	myHand := Bid{}
	err := json.Unmarshal(areq, &myHand)
	if err != nil {
		fmt.Println("JSONtoBid error: ", err)
		return myHand, err
	}
	return myHand, err
}

//////////////////////////////////////////////////////////
// Converts an User Object to a JSON String
//////////////////////////////////////////////////////////
func UsertoJSON(user UserObject) ([]byte, error) {

	ajson, err := json.Marshal(user)
	if err != nil {
		fmt.Println("UsertoJSON error: ", err)
		return nil, err
	}
	fmt.Println("UsertoJSON created: ", ajson)
	return ajson, nil
}

//////////////////////////////////////////////////////////
// Converts an User Object to a JSON String
//////////////////////////////////////////////////////////
func JSONtoUser(user []byte) (UserObject, error) {

	ur := UserObject{}
	err := json.Unmarshal(user, &ur)
	if err != nil {
		fmt.Println("JSONtoUser error: ", err)
		return ur, err
	}
	fmt.Println("JSONtoUser created: ", ur)
	return ur, err
}

//////////////////////////////////////////////////////////
// Converts an Item Transaction to a JSON String
//////////////////////////////////////////////////////////
func TrantoJSON(at ItemTransaction) ([]byte, error) {

	ajson, err := json.Marshal(at)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ajson, nil
}

//////////////////////////////////////////////////////////
// Converts an Trans Object to a JSON String
//////////////////////////////////////////////////////////
func JSONtoTran(areq []byte) (ItemTransaction, error) {

	at := ItemTransaction{}
	err := json.Unmarshal(areq, &at)
	if err != nil {
		fmt.Println("JSONtoTran error: ", err)
		return at, err
	}
	return at, err
}

//////////////////////////////////////////////
// Validates an ID for Well Formed
//////////////////////////////////////////////

func validateID(id string) error {
	// Validate UserID is an integer

	_, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("validateID(): User ID should be an integer")
	}
	return nil
}

//////////////////////////////////////////////
// Create an ItemLog from Item
//////////////////////////////////////////////

func ItemToItemLog(io ItemObject) ItemLog {

	iLog := ItemLog{}
	iLog.ItemID = io.ItemID
	iLog.Status = "INITIAL"
	iLog.AuctionedBy = "DEFAULT"
	iLog.RecType = "ILOG"
	iLog.ItemDesc = io.ItemDesc
	iLog.CurrentOwner = io.CurrentOwnerID
	iLog.Date = time.Now().Format("2006-01-02 15:04:05")

	return iLog
}

//////////////////////////////////////////////
// Convert Bid to Transaction for Posting
//////////////////////////////////////////////

func BidtoTransaction(bid Bid) ItemTransaction {

	var t ItemTransaction
	t.AuctionID = bid.AuctionID
	t.RecType = "POSTTRAN"
	t.ItemID = bid.ItemID
	t.TransType = "SALE"
	t.UserId = bid.BuyerID
	t.TransDate = time.Now().Format("2006-01-02 15:04:05")
	t.HammerTime = bid.BidTime
	t.HammerPrice = bid.BidPrice
	t.Details = "The Highest Bidder does not always win"

	return t
}

////////////////////////////////////////////////////////////////////////////
// Validate if the User Information Exists
// in the block-chain
////////////////////////////////////////////////////////////////////////////
func ValidateMember(stub shim.ChaincodeStubInterface, owner string) ([]byte, error) {

	// Get the Item Objects and Display it
	// Avalbytes, err := stub.GetState(owner)
	args := []string{owner, "USER"}
	Avalbytes, err := QueryLedger(stub, "UserTable", args)

	if err != nil {
		fmt.Println("ValidateMember() : Failed - Cannot find valid owner record for ART  ", owner)
		jsonResp := "{\"Error\":\"Failed to get Owner Object Data for " + owner + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		fmt.Println("ValidateMember() : Failed - Incomplete owner record for ART  ", owner)
		jsonResp := "{\"Error\":\"Failed - Incomplete information about the owner for " + owner + "\"}"
		return nil, errors.New(jsonResp)
	}

	fmt.Println("ValidateMember() : Validated Item Owner:\n", owner)
	return Avalbytes, nil
}

////////////////////////////////////////////////////////////////////////////
// Validate if the User Information Exists
// in the block-chain
////////////////////////////////////////////////////////////////////////////
func ValidateItemSubmission(stub shim.ChaincodeStubInterface, artId string) ([]byte, error) {

	// Get the Item Objects and Display it
	args := []string{artId, "ARTINV"}
	Avalbytes, err := QueryLedger(stub, "ItemTable", args)
	if err != nil {
		fmt.Println("ValidateItemSubmission() : Failed - Cannot find valid owner record for ART  ", artId)
		jsonResp := "{\"Error\":\"Failed to get Owner Object Data for " + artId + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		fmt.Println("ValidateItemSubmission() : Failed - Incomplete owner record for ART  ", artId)
		jsonResp := "{\"Error\":\"Failed - Incomplete information about the owner for " + artId + "\"}"
		return nil, errors.New(jsonResp)
	}

	//fmt.Println("ValidateItemSubmission() : Validated Item Owner:", Avalbytes)
	return Avalbytes, nil
}

////////////////////////////////////////////////////////////////////////////
// Open a Ledgers if one does not exist
// These ledgers will be used to write /  read data
// Use names are listed in aucTables {}
// THIS FUNCTION REPLACES ALL THE INIT Functions below
//  - InitUserReg()
//  - InitAucReg()
//  - InitBidReg()
//  - InitItemReg()
//  - InitItemMaster()
//  - InitTransReg()
//  - InitAuctionTriggerReg()
//  - etc. etc.
////////////////////////////////////////////////////////////////////////////
func InitLedger(stub shim.ChaincodeStubInterface, tableName string) error {

	// Generic Table Creation Function - requires Table Name and Table Key Entry
	// Create Table - Get number of Keys the tables supports
	// This version assumes all Keys are String and the Data is Bytes
	// This Function can replace all other InitLedger function in this app such as InitItemLedger()

	nKeys := GetNumberOfKeys(tableName)
	if nKeys < 1 {
		fmt.Println("Atleast 1 Key must be provided \n")
		fmt.Println("Auction_Application: Failed creating Table ", tableName)
		return errors.New("Auction_Application: Failed creating Table " + tableName)
	}

	var columnDefsForTbl []*shim.ColumnDefinition

	for i := 0; i < nKeys; i++ {
		columnDef := shim.ColumnDefinition{Name: "keyName" + strconv.Itoa(i), Type: shim.ColumnDefinition_STRING, Key: true}
		columnDefsForTbl = append(columnDefsForTbl, &columnDef)
	}

	columnLastTblDef := shim.ColumnDefinition{Name: "Details", Type: shim.ColumnDefinition_BYTES, Key: false}
	columnDefsForTbl = append(columnDefsForTbl, &columnLastTblDef)

	// Create the Table (Nil is returned if the Table exists or if the table is created successfully
	err := stub.CreateTable(tableName, columnDefsForTbl)

	if err != nil {
		fmt.Println("Auction_Application: Failed creating Table ", tableName)
		return errors.New("Auction_Application: Failed creating Table " + tableName)
	}

	return err
}

////////////////////////////////////////////////////////////////////////////
// Open a User Registration Table if one does not exist
// Register users into this table
////////////////////////////////////////////////////////////////////////////
func UpdateLedger(stub shim.ChaincodeStubInterface, tableName string, keys []string, args []byte) error {

	nKeys := GetNumberOfKeys(tableName)
	if nKeys < 1 {
		fmt.Println("Atleast 1 Key must be provided \n")
	}

	var columns []*shim.Column

	for i := 0; i < nKeys; i++ {
		col := shim.Column{Value: &shim.Column_String_{String_: keys[i]}}
		columns = append(columns, &col)
	}

	lastCol := shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(args)}}
	columns = append(columns, &lastCol)

	row := shim.Row{columns}
	ok, err := stub.InsertRow(tableName, row)
	if err != nil {
		return fmt.Errorf("UpdateLedger: InsertRow into "+tableName+" Table operation failed. %s", err)
	}
	if !ok {
		return errors.New("UpdateLedger: InsertRow into " + tableName + " Table failed. Row with given key " + keys[0] + " already exists")
	}

	fmt.Println("UpdateLedger: InsertRow into ", tableName, " Table operation Successful. ")
	return nil
}

////////////////////////////////////////////////////////////////////////////
// Open a User Registration Table if one does not exist
// Register users into this table
////////////////////////////////////////////////////////////////////////////
func DeleteFromLedger(stub shim.ChaincodeStubInterface, tableName string, keys []string) error {
	var columns []shim.Column

	//nKeys := GetNumberOfKeys(tableName)
	nCol := len(keys)
	if nCol < 1 {
		fmt.Println("Atleast 1 Key must be provided \n")
		return errors.New("DeleteFromLedger failed. Must include at least key values")
	}

	for i := 0; i < nCol; i++ {
		colNext := shim.Column{Value: &shim.Column_String_{String_: keys[i]}}
		columns = append(columns, colNext)
	}

	err := stub.DeleteRow(tableName, columns)
	if err != nil {
		return fmt.Errorf("DeleteFromLedger operation failed. %s", err)
	}

	fmt.Println("DeleteFromLedger: DeleteRow from ", tableName, " Table operation Successful. ")
	return nil
}

////////////////////////////////////////////////////////////////////////////
// Replaces the Entry in the Ledger
//
////////////////////////////////////////////////////////////////////////////
func ReplaceLedgerEntry(stub shim.ChaincodeStubInterface, tableName string, keys []string, args []byte) error {

	nKeys := GetNumberOfKeys(tableName)
	if nKeys < 1 {
		fmt.Println("Atleast 1 Key must be provided \n")
	}

	var columns []*shim.Column

	for i := 0; i < nKeys; i++ {
		col := shim.Column{Value: &shim.Column_String_{String_: keys[i]}}
		columns = append(columns, &col)
	}

	lastCol := shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(args)}}
	columns = append(columns, &lastCol)

	row := shim.Row{columns}
	ok, err := stub.ReplaceRow(tableName, row)
	if err != nil {
		return fmt.Errorf("ReplaceLedgerEntry: Replace Row into "+tableName+" Table operation failed. %s", err)
	}
	if !ok {
		return errors.New("ReplaceLedgerEntry: Replace Row into " + tableName + " Table failed. Row with given key " + keys[0] + " already exists")
	}

	fmt.Println("ReplaceLedgerEntry: Replace Row in ", tableName, " Table operation Successful. ")
	return nil
}

////////////////////////////////////////////////////////////////////////////
// Query a User Object by Table Name and Key
////////////////////////////////////////////////////////////////////////////
func QueryLedger(stub shim.ChaincodeStubInterface, tableName string, args []string) ([]byte, error) {

	var columns []shim.Column
	nCol := GetNumberOfKeys(tableName)
	for i := 0; i < nCol; i++ {
		colNext := shim.Column{Value: &shim.Column_String_{String_: args[i]}}
		columns = append(columns, colNext)
	}

	row, err := stub.GetRow(tableName, columns)
	fmt.Println("Length or number of rows retrieved ", len(row.Columns))

	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed retrieving data " + args[0] + ". \"}"
		fmt.Println("Error retrieving data record for Key = ", args[0], "Error : ", jsonResp)
		return nil, errors.New(jsonResp)
	}

	//fmt.Println("User Query Response:", row)
	//jsonResp := "{\"Owner\":\"" + string(row.Columns[nCol].GetBytes()) + "\"}"
	//fmt.Println("User Query Response:%s\n", jsonResp)
	Avalbytes := row.Columns[nCol].GetBytes()

	// Perform Any additional processing of data
	fmt.Println("QueryLedger() : Successful - Proceeding to ProcessRequestType ")
	err = ProcessQueryResult(stub, Avalbytes, args)
	if err != nil {
		fmt.Println("QueryLedger() : Cannot create object  : ", args[1])
		jsonResp := "{\"QueryLedger() Error\":\" Cannot create Object for key " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}
	return Avalbytes, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// Get List of Bids for an Auction
// in the block-chain --
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetListOfBids", "Args": ["1111"]}'
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetLastBid", "Args": ["1111"]}'
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetHighestBid", "Args": ["1111"]}'
/////////////////////////////////////////////////////////////////////////////////////////////////////
func GetListOfBids(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	rows, err := GetList(stub, "BidTable", args)
	if err != nil {
		return nil, fmt.Errorf("GetListOfBids operation failed. Error marshaling JSON: %s", err)
	}

	nCol := GetNumberOfKeys("BidTable")

	tlist := make([]Bid, len(rows))
	for i := 0; i < len(rows); i++ {
		ts := rows[i].Columns[nCol].GetBytes()
		bid, err := JSONtoBid(ts)
		if err != nil {
			fmt.Println("GetListOfBids() Failed : Ummarshall error")
			return nil, fmt.Errorf("GetListOfBids() operation failed. %s", err)
		}
		tlist[i] = bid
	}

	jsonRows, _ := json.Marshal(tlist)

	fmt.Println("List of Bids Requested : ", jsonRows)
	return jsonRows, nil

}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Get List of Auctions that have been initiated
// in the block-chain
// This is a fixed Query to be issued as below
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetListOfInitAucs", "Args": ["2016"]}'
////////////////////////////////////////////////////////////////////////////////////////////////////////
func GetListOfInitAucs(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	rows, err := GetList(stub, "AucInitTable", args)
	if err != nil {
		return nil, fmt.Errorf("GetListOfInitAucs operation failed. Error marshaling JSON: %s", err)
	}

	nCol := GetNumberOfKeys("AucInitTable")

	tlist := make([]AuctionRequest, len(rows))
	for i := 0; i < len(rows); i++ {
		ts := rows[i].Columns[nCol].GetBytes()
		ar, err := JSONtoAucReq(ts)
		if err != nil {
			fmt.Println("GetListOfInitAucs() Failed : Ummarshall error")
			return nil, fmt.Errorf("getBillForMonth() operation failed. %s", err)
		}
		tlist[i] = ar
	}

	jsonRows, _ := json.Marshal(tlist)

	//fmt.Println("List of Auctions Requested : ", jsonRows)
	return jsonRows, nil

}

////////////////////////////////////////////////////////////////////////////
// Get List of Open Auctions  for which bids can be supplied
// in the block-chain
// This is a fixed Query to be issued as below
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetListOfOpenAucs", "Args": ["2016"]}'
////////////////////////////////////////////////////////////////////////////
func GetListOfOpenAucs(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	rows, err := GetList(stub, "AucOpenTable", args)
	if err != nil {
		return nil, fmt.Errorf("GetListOfOpenAucs operation failed. Error marshaling JSON: %s", err)
	}

	nCol := GetNumberOfKeys("AucOpenTable")

	tlist := make([]AuctionRequest, len(rows))
	for i := 0; i < len(rows); i++ {
		ts := rows[i].Columns[nCol].GetBytes()
		ar, err := JSONtoAucReq(ts)
		if err != nil {
			fmt.Println("GetListOfOpenAucs() Failed : Ummarshall error")
			return nil, fmt.Errorf("GetListOfOpenAucs() operation failed. %s", err)
		}
		tlist[i] = ar
	}

	jsonRows, _ := json.Marshal(tlist)

	//fmt.Println("List of Open Auctions : ", jsonRows)
	return jsonRows, nil

}

////////////////////////////////////////////////////////////////////////////
// Get the Item History for an Item
// in the block-chain .. Pass the Item ID
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetItemLog", "Args": ["1000"]}'
////////////////////////////////////////////////////////////////////////////
func GetItemLog(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Check there are 1 Arguments provided as per the the struct - two are computed
	// See example
	if len(args) < 1 {
		fmt.Println("GetItemLog(): Incorrect number of arguments. Expecting 1 ")
		fmt.Println("GetItemLog(): ./peer chaincode query -l golang -n mycc -c '{\"Function\": \"GetItem\", \"Args\": [\"1111\"]}'")
		return nil, errors.New("CreateItemObject(): Incorrect number of arguments. Expecting 12 ")
	}

	rows, err := GetList(stub, "ItemHistoryTable", args)
	if err != nil {
		return nil, fmt.Errorf("GetItemLog() operation failed. Error marshaling JSON: %s", err)
	}
	nCol := GetNumberOfKeys("ItemHistoryTable")

	tlist := make([]ItemLog, len(rows))
	for i := 0; i < len(rows); i++ {
		ts := rows[i].Columns[nCol].GetBytes()
		il, err := JSONtoItemLog(ts)
		if err != nil {
			fmt.Println("() Failed : Ummarshall error")
			return nil, fmt.Errorf("GetItemLog() operation failed. %s", err)
		}
		tlist[i] = il
	}

	jsonRows, _ := json.Marshal(tlist)

	//fmt.Println("All History : ", jsonRows)
	return jsonRows, nil

}

////////////////////////////////////////////////////////////////////////////
// Get a List of Items by Category
// in the block-chain
// Input is 2016 + Category
// Categories include whatever has been defined in the Item Tables - Landscape, Modern, ...
// See Sample data
// ./peer chaincode query -l golang -n mycc -c '{"Function": "GetItemListByCat", "Args": ["2016", "Modern"]}'
////////////////////////////////////////////////////////////////////////////
func GetItemListByCat(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Check there are 1 Arguments provided as per the the struct - two are computed
	// See example
	if len(args) < 1 {
		fmt.Println("GetItemListByCat(): Incorrect number of arguments. Expecting 1 ")
		fmt.Println("GetItemListByCat(): ./peer chaincode query -l golang -n mycc -c '{\"Function\": \"GetItemListByCat\", \"Args\": [\"Modern\"]}'")
		return nil, errors.New("CreateItemObject(): Incorrect number of arguments. Expecting 1 ")
	}

	rows, err := GetList(stub, "ItemCatTable", args)
	if err != nil {
		return nil, fmt.Errorf("GetItemListByCat() operation failed. Error GetList: %s", err)
	}

	nCol := GetNumberOfKeys("ItemCatTable")

	tlist := make([]ItemObject, len(rows))
	for i := 0; i < len(rows); i++ {
		ts := rows[i].Columns[nCol].GetBytes()
		io, err := JSONtoAR(ts)
		if err != nil {
			fmt.Println("() Failed : Ummarshall error")
			return nil, fmt.Errorf("GetItemListByCat() operation failed. %s", err)
		}
		//TODO: Masking Image binary data, Need a clean solution ?
		io.ItemImage = []byte{}
		tlist[i] = io
	}

	jsonRows, _ := json.Marshal(tlist)

	//fmt.Println("All Items : ", jsonRows)
	return jsonRows, nil

}

////////////////////////////////////////////////////////////////////////////
// Get a List of Users by Category
// in the block-chain
////////////////////////////////////////////////////////////////////////////
func GetUserListByCat(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Check there are 1 Arguments provided as per the the struct - two are computed
	// See example
	if len(args) < 1 {
		fmt.Println("GetUserListByCat(): Incorrect number of arguments. Expecting 1 ")
		fmt.Println("GetUserListByCat(): ./peer chaincode query -l golang -n mycc -c '{\"Function\": \"GetUserListByCat\", \"Args\": [\"AH\"]}'")
		return nil, errors.New("CreateUserObject(): Incorrect number of arguments. Expecting 1 ")
	}

	rows, err := GetList(stub, "UserCatTable", args)
	if err != nil {
		return nil, fmt.Errorf("GetUserListByCat() operation failed. Error marshaling JSON: %s", err)
	}

	nCol := GetNumberOfKeys("UserCatTable")

	tlist := make([]UserObject, len(rows))
	for i := 0; i < len(rows); i++ {
		ts := rows[i].Columns[nCol].GetBytes()
		uo, err := JSONtoUser(ts)
		if err != nil {
			fmt.Println("GetUserListByCat() Failed : Ummarshall error")
			return nil, fmt.Errorf("GetUserListByCat() operation failed. %s", err)
		}
		tlist[i] = uo
	}

	jsonRows, _ := json.Marshal(tlist)

	//fmt.Println("All Users : ", jsonRows)
	return jsonRows, nil

}

////////////////////////////////////////////////////////////////////////////
// Get a List of Rows based on query criteria from the OBC
//
////////////////////////////////////////////////////////////////////////////
func GetList(stub shim.ChaincodeStubInterface, tableName string, args []string) ([]shim.Row, error) {
	var columns []shim.Column

	nKeys := GetNumberOfKeys(tableName)
	nCol := len(args)
	if nCol < 1 {
		fmt.Println("Atleast 1 Key must be provided \n")
		return nil, errors.New("GetList failed. Must include at least key values")
	}

	for i := 0; i < nCol; i++ {
		colNext := shim.Column{Value: &shim.Column_String_{String_: args[i]}}
		columns = append(columns, colNext)
	}

	rowChannel, err := stub.GetRows(tableName, columns)
	if err != nil {
		return nil, fmt.Errorf("GetList operation failed. %s", err)
	}
	var rows []shim.Row
	for {
		select {
		case row, ok := <-rowChannel:
			if !ok {
				rowChannel = nil
			} else {
				rows = append(rows, row)
				//If required enable for debugging
				//fmt.Println(row)
			}
		}
		if rowChannel == nil {
			break
		}
	}

	fmt.Println("Number of Keys retrieved : ", nKeys)
	fmt.Println("Number of rows retrieved : ", len(rows))
	return rows, nil
}

////////////////////////////////////////////////////////////////////////////
// Get The Highest Bid Received so far for an Auction
// in the block-chain
////////////////////////////////////////////////////////////////////////////
func GetLastBid(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	tn := "BidTable"
	rows, err := GetList(stub, tn, args)
	if err != nil {
		return nil, fmt.Errorf("GetLastBid operation failed. %s", err)
	}
	nCol := GetNumberOfKeys(tn)
	var Avalbytes []byte
	var dat map[string]interface{}
	layout := "2006-01-02 15:04:05"
	highestTime, err := time.Parse(layout, layout)

	for i := 0; i < len(rows); i++ {
		currentBid := rows[i].Columns[nCol].GetBytes()
		if err := json.Unmarshal(currentBid, &dat); err != nil {
			fmt.Println("GetHighestBid() Failed : Ummarshall error")
			return nil, fmt.Errorf("GetHighestBid(0 operation failed. %s", err)
		}
		bidTime, err := time.Parse(layout, dat["BidTime"].(string))
		if err != nil {
			fmt.Println("GetLastBid() Failed : time Conversion error on BidTime")
			return nil, fmt.Errorf("GetHighestBid() Int Conversion error on BidPrice! failed. %s", err)
		}

		if bidTime.Sub(highestTime) > 0 {
			highestTime = bidTime
			Avalbytes = currentBid
		}
	}

	return Avalbytes, nil

}

////////////////////////////////////////////////////////////////////////////
// Get The Highest Bid Received so far for an Auction
// in the block-chain
////////////////////////////////////////////////////////////////////////////
func GetNoOfBidsReceived(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	tn := "BidTable"
	rows, err := GetList(stub, tn, args)
	if err != nil {
		return nil, fmt.Errorf("GetLastBid operation failed. %s", err)
	}
	nBids := len(rows)
	return []byte(strconv.Itoa(nBids)), nil
}

////////////////////////////////////////////////////////////////////////////
// Get the Highest Bid in the List
//
////////////////////////////////////////////////////////////////////////////
func GetHighestBid(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	tn := "BidTable"
	rows, err := GetList(stub, tn, args)
	if err != nil {
		return nil, fmt.Errorf("GetLastBid operation failed. %s", err)
	}
	nCol := GetNumberOfKeys(tn)
	var Avalbytes []byte
	var dat map[string]interface{}
	var bidPrice, highestBid int
	highestBid = 0

	for i := 0; i < len(rows); i++ {
		currentBid := rows[i].Columns[nCol].GetBytes()
		if err := json.Unmarshal(currentBid, &dat); err != nil {
			fmt.Println("GetHighestBid() Failed : Ummarshall error")
			return nil, fmt.Errorf("GetHighestBid(0 operation failed. %s", err)
		}
		bidPrice, err = strconv.Atoi(dat["BidPrice"].(string))
		if err != nil {
			fmt.Println("GetHighestBid() Failed : Int Conversion error on BidPrice")
			return nil, fmt.Errorf("GetHighestBid() Int Conversion error on BidPrice! failed. %s", err)
		}

		if bidPrice >= highestBid {
			highestBid = bidPrice
			Avalbytes = currentBid
		}
	}

	return Avalbytes, nil
}

/////////////////////////////////////////////////////////////////
// This function checks the incoming args stuff for a valid record
// type entry as per the declared array recType[]
// The assumption is that rectType can be anywhere in the args or struct
// not necessarily in args[1] as per my old logic
// The Request type is used to process the record accordingly
/////////////////////////////////////////////////////////////////
func IdentifyReqType(args []string) string {
	for _, rt := range args {
		for _, val := range recType {
			if val == rt {
				return rt
			}
		}
	}
	return "DEFAULT"
}

/////////////////////////////////////////////////////////////////
// This function checks the incoming args stuff for a valid record
// type entry as per the declared array recType[]
// The assumption is that rectType can be anywhere in the args or struct
// not necessarily in args[1] as per my old logic
// The Request type is used to process the record accordingly
/////////////////////////////////////////////////////////////////
func ChkReqType(args []string) bool {
	for _, rt := range args {
		for _, val := range recType {
			if val == rt {
				return true
			}
		}
	}
	return false
}

/////////////////////////////////////////////////////////////////
// Checks if the incoming invoke has a valid requesType
// The Request type is used to process the record accordingly
// Old Logic (see new logic up)
/////////////////////////////////////////////////////////////////
func CheckRequestType(rt string) bool {
	for _, val := range recType {
		if val == rt {
			fmt.Println("CheckRequestType() : Valid Request Type , val : ", val, rt, "\n")
			return true
		}
	}
	fmt.Println("CheckRequestType() : Invalid Request Type , val : ", rt, "\n")
	return false
}

/////////////////////////////////////////////////////////////////////////////////////////////
// Return the right Object Buffer after validation to write to the ledger
// var recType = []string{"ARTINV", "USER", "BID", "AUCREQ", "POSTTRAN", "OPENAUC", "CLAUC"}
/////////////////////////////////////////////////////////////////////////////////////////////

func ProcessQueryResult(stub shim.ChaincodeStubInterface, Avalbytes []byte, args []string) error {

	// Identify Record Type by scanning the args for one of the recTypes
	// This is kind of a post-processor once the query fetches the results
	// RecType is the style of programming in the punch card days ..
	// ... well

	var dat map[string]interface{}

	if err := json.Unmarshal(Avalbytes, &dat); err != nil {
		panic(err)
	}

	var recType string
	recType = dat["RecType"].(string)
	switch recType {

	case "ARTINV":

		ar, err := JSONtoAR(Avalbytes) //
		if err != nil {
			fmt.Println("ProcessRequestType(): Cannot create itemObject \n")
			return err
		}
		// Decrypt Image and Save Image in a file
		image := Decrypt(ar.AES_Key, ar.ItemImage)
		if err != nil {
			fmt.Println("ProcessRequestType() : Image decrytion failed ")
			return err
		}
		fmt.Println("ProcessRequestType() : Image conversion from byte[] to file successfull ")
		err = ByteArrayToImage(image, ccPath+"copy."+ar.ItemPicFN)
		if err != nil {

			fmt.Println("ProcessRequestType() : Image conversion from byte[] to file failed ")
			return err
		}
		return err

	case "USER":
		ur, err := JSONtoUser(Avalbytes) //
		if err != nil {
			return err
		}
		fmt.Println("ProcessRequestType() : ", ur)
		return err

	case "AUCREQ":
	case "OPENAUC":
	case "CLAUC":
		ar, err := JSONtoAucReq(Avalbytes) //
		if err != nil {
			return err
		}
		fmt.Println("ProcessRequestType() : ", ar)
		return err
	case "POSTTRAN":
		atr, err := JSONtoTran(Avalbytes) //
		if err != nil {
			return err
		}
		fmt.Println("ProcessRequestType() : ", atr)
		return err
	case "BID":
		bid, err := JSONtoBid(Avalbytes) //
		if err != nil {
			return err
		}
		fmt.Println("ProcessRequestType() : ", bid)
		return err
	case "DEFAULT":
		return nil
	case "XFER":
		return nil
	case "VERIFY":
		return nil
	default:

		return errors.New("Unknown")
	}
	return nil

}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Trigger the Auction
// Structure of args auctionReqID, RecType, Duration in Minutes ( 3 = 3 minutes)
// ./peer chaincode invoke -l golang -n mycc -c '{"Function": "OpenAuctionForBids", "Args":["1111", "OPENAUC", "3"]}'
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func OpenAuctionForBids(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Fetch Auction Object and check its Status
	Avalbytes, err := QueryLedger(stub, "AuctionTable", args)
	if err != nil {
		fmt.Println("OpenAuctionForBids(): Auction Object Retrieval Failed ")
		return nil, errors.New("OpenAuctionForBids(): Auction Object Retrieval Failed ")
	}

	aucR, err := JSONtoAucReq(Avalbytes)
	if err != nil {
		fmt.Println("OpenAuctionForBids(): Auction Object Unmarshalling Failed ")
		return nil, errors.New("OpenAuctionForBids(): Auction Object UnMarshalling Failed ")
	}

	if aucR.Status == "CLOSED" {
		fmt.Println("OpenAuctionForBids(): Auction is Closed - Cannot Open for new bids ")
		return nil, errors.New("OpenAuctionForBids(): is Closed - Cannot Open for new bids Failed ")
	}

	// Calculate Time Now and Duration of Auction

	// Validate arg[1]  is an integer as it represents Duration in Minutes
	aucDuration, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("OpenAuctionForBids(): Auction Duration is an integer that represents minute! OpenAuctionForBids() Failed ")
		return nil, errors.New("OpenAuctionForBids(): Auction Duration is an integer that represents minute! OpenAuctionForBids() Failed ")
	}

	aucStartDate := time.Now()
	aucEndDate := aucStartDate.Add(time.Duration(aucDuration) * time.Minute)
	//sleepTime := time.Duration(aucDuration * 60 * 1000 * 1000 * 1000)

	//  Update Auction Object
	aucR.OpenDate = aucStartDate.Format("2006-01-02 15:04:05")
	aucR.CloseDate = aucEndDate.Format("2006-01-02 15:04:05")
	aucR.Status = "OPEN"

	buff, err := UpdateAuctionStatus(stub, "AuctionTable", aucR)
	if err != nil {
		fmt.Println("OpenAuctionForBids(): UpdateAuctionStatus() Failed ")
		return nil, errors.New("OpenAuctionForBids(): UpdateAuctionStatus() Failed ")
	}

	// Remove the Auction from INIT Bucket and move to OPEN bucket
	// This was designed primarily to help the UI

	keys := []string{"2016", aucR.AuctionID}
	err = DeleteFromLedger(stub, "AucInitTable", keys)
	if err != nil {
		fmt.Println("OpenAuctionForBids(): DeleteFromLedger() Failed ")
		return nil, errors.New("OpenAuctionForBids(): DeleteFromLedger() Failed ")
	}

	// Add the Auction to Open Bucket
	err = UpdateLedger(stub, "AucOpenTable", keys, buff)
	if err != nil {
		fmt.Println("OpenAuctionForBids() : write error while inserting record into AucInitTable \n")
		return buff, err
	}

	// Initiate Timer for the duration of the Auction
	// Bids are accepted as long as the timer is alive
	/*go func(aucR AuctionRequest, sleeptime time.Duration) ([]byte, error) {
		fmt.Println("OpenAuctionForBids(): Sleeping for ", sleeptime)
		time.Sleep(sleeptime)

		// Exec The following Command from the shell
		ShellCmdToCloseAuction(aucR.AuctionID)
		return nil, err
	}(aucR, sleepTime)*/
	return buff, err
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Create a Command to execute Close Auction From the Command line
// cloaseauction.sh is created and then executed as seen below
// The file contains just one line
// /opt/gopath/src/github.com/hyperledger/fabric/peer chaincode invoke -l golang -n mycc -c '{"Function": "CloseAuction", "Args": ["1111","AUCREQ"]}'
// This approach has been used as opposed to exec.Command... because additional logic to gather environment variables etc. is required
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ShellCmdToCloseAuction(aucID string) error {
	gopath := os.Getenv("GOPATH")
	cdir := fmt.Sprintf("cd %s/src/github.com/hyperledger/fabric/", gopath)
	argStr := "'{\"Function\": \"CloseAuction\", \"Args\": [\"" + aucID + "\"," + "\"AUCREQ\"" + "]}'"
	argStr = fmt.Sprintf("%s/src/github.com/hyperledger/fabric/peer/peer chaincode invoke -l golang -n mycc -c %s", gopath, argStr)

	fileHandle, _ := os.Create(fmt.Sprintf("%s/src/github.com/hyperledger/fabric/peer/closeauction.sh", gopath))
	writer := bufio.NewWriter(fileHandle)
	defer fileHandle.Close()

	fmt.Fprintln(writer, cdir)
	fmt.Fprintln(writer, argStr)
	writer.Flush()

	x := "sh /opt/gopath/src/github.com/hyperledger/fabric/peer/closeauction.sh"
	err := exe_cmd(x)
	if err != nil {
		fmt.Println("%s", err)
	}

	err = exe_cmd("rm /opt/gopath/src/github.com/hyperledger/fabric/peer/closeauction.sh")
	if err != nil {
		fmt.Println("%s", err)
	}

	fmt.Println("Kicking off CloseAuction", argStr)
	return nil
}

func exe_cmd(cmd string) error {

	fmt.Println("command :  ", cmd)
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	_, err := exec.Command(head, parts...).CombinedOutput()
	if err != nil {
		fmt.Println("%s", err)
	}
	return err
}

//////////////////////////////////////////////////////////////////////////
// Close Open Auctions
// 1. Read OpenAucTable
// 2. Compare now with expiry time with now
// 3. If now is > expiry time call CloseAuction
//////////////////////////////////////////////////////////////////////////

func CloseOpenAuctions(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	rows, err := GetListOfOpenAucs(stub, "AucOpenTable", []string{"2016"})
	if err != nil {
		return nil, fmt.Errorf("GetListOfOpenAucs operation failed. Error marshaling JSON: %s", err)
	}

	tlist := make([]AuctionRequest, len(rows))
	err = json.Unmarshal([]byte(rows), &tlist)
	if err != nil {
		fmt.Println("Unmarshal failed : ", err)
	}

	for i := 0; i < len(tlist); i++ {
		ar := tlist[i]
		if err != nil {
			fmt.Println("CloseOpenAuctions() Failed : Ummarshall error")
			return nil, fmt.Errorf("GetListOfOpenAucs() operation failed. %s", err)
		}

		fmt.Println("CloseOpenAuctions() ", ar)

		// Compare Auction Times
		if tCompare(time.Now().Format("2006-01-02 15:04:05"), ar.CloseDate) == false {

			// Request Closing Auction
			_, err := CloseAuction(stub, "CloseAuction", []string{ar.AuctionID})
			if err != nil {
				fmt.Println("CloseOpenAuctions() Failed : Ummarshall error")
				return nil, fmt.Errorf("GetListOfOpenAucs() operation failed. %s", err)
			}
		}
	}

	return rows, nil
}

//////////////////////////////////////////////////////////////////////////
// Close the Auction
// This is invoked by OpenAuctionForBids
// which kicks-off a go-routine timer for the duration of the auction
// When the timer expires, it creates a shell script to CloseAuction() and triggers it
// This function can also be invoked via CLI - the intent was to close as and when I implement BuyItNow()
// CloseAuction
// - Sets the status of the Auction to "CLOSED"
// - Removes the Auction from the Open Auction list (AucOpenTable)
// - Retrieves the Highest Bid and creates a Transaction
// - Posts The Transaction
//
// To invoke from Command Line via CLI or REST API
// ./peer chaincode invoke -l golang -n mycc -c '{"Function": "CloseAuction", "Args": ["1111", "AUCREQ"]}'
// ./peer chaincode invoke -l golang -n mycc -c '{"Function": "CloseAuction", "Args": ["1111", "AUCREQ"]}'
//
//////////////////////////////////////////////////////////////////////////

func CloseAuction(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Close The Auction -  Fetch Auction Object
	Avalbytes, err := QueryLedger(stub, "AuctionTable", []string{args[0], "AUCREQ"})
	if err != nil {
		fmt.Println("CloseAuction(): Auction Object Retrieval Failed ")
		return nil, errors.New("CloseAuction(): Auction Object Retrieval Failed ")
	}

	aucR, err := JSONtoAucReq(Avalbytes)
	if err != nil {
		fmt.Println("CloseAuction(): Auction Object Unmarshalling Failed ")
		return nil, errors.New("CloseAuction(): Auction Object UnMarshalling Failed ")
	}

	//  Update Auction Status
	aucR.Status = "CLOSED"
	fmt.Println("CloseAuction(): UpdateAuctionStatus() successful ", aucR)

	Avalbytes, err = UpdateAuctionStatus(stub, "AuctionTable", aucR)
	if err != nil {
		fmt.Println("CloseAuction(): UpdateAuctionStatus() Failed ")
		return nil, errors.New("CloseAuction(): UpdateAuctionStatus() Failed ")
	}

	// Remove the Auction from Open Bucket
	keys := []string{"2016", aucR.AuctionID}
	err = DeleteFromLedger(stub, "AucOpenTable", keys)
	if err != nil {
		fmt.Println("CloseAuction(): DeleteFromLedger(AucOpenTable) Failed ")
		return nil, errors.New("CloseAuction(): DeleteFromLedger(AucOpenTable) Failed ")
	}

	fmt.Println("CloseAuction(): Proceeding to process the highest bid ")

	// Process Final Bid - Turn it into a Transaction
	Avalbytes, err = GetHighestBid(stub, "GetHighestBid", []string{args[0]})
	if Avalbytes == nil {
		fmt.Println("CloseAuction(): No bids available, no change in Item Status - PostTransaction() Completed Successfully ")
		return Avalbytes, nil
	}

	if err != nil {
		fmt.Println("CloseAuction(): No bids available, error encountered - PostTransaction() failed ")
		return nil, err
	}

	bid, _ := JSONtoBid(Avalbytes)
	fmt.Println("CloseAuction(): Proceeding to process the highest bid ", bid)
	tran := BidtoTransaction(bid)
	fmt.Println("CloseAuction(): Converting Bid to tran ", tran)

	// Process the last bid once Time Expires
	tranArgs := []string{tran.AuctionID, tran.RecType, tran.ItemID, tran.TransType, tran.UserId, tran.TransDate, tran.HammerTime, tran.HammerPrice, tran.Details}
	fmt.Println("CloseAuction(): Proceeding to process the  Transaction ", tranArgs)

	Avalbytes, err = PostTransaction(stub, "PostTransaction", tranArgs)
	if err != nil {
		fmt.Println("CloseAuction(): PostTransaction() Failed ")
		return nil, errors.New("CloseAuction(): PostTransaction() Failed ")
	}
	fmt.Println("CloseAuction(): PostTransaction() Completed Successfully ")
	return Avalbytes, nil
}

////////////////////////////////////////////////////////////////////////////////////////////
// Buy It Now
// Rules:
// If Buy IT Now Option is available then a Buyer has the option to buy the ITEM
// before the bids exceed BuyITNow Price . Normally, The application should take of this
// at the UI level and this chain-code assumes application has validated that
////////////////////////////////////////////////////////////////////////////////////////////

func BuyItNow(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Process Final Bid - Turn it into a Transaction
	Avalbytes, err := GetHighestBid(stub, "GetHighestBid", []string{args[0]})
	hBidFlag := true
	if Avalbytes == nil {
		fmt.Println("BuyItNow(): No bids available, no change in Item Status - PostTransaction() Completed Successfully ")
		hBidFlag = false
	}

	if err != nil {
		fmt.Println("BuyItNow(): No bids available, error encountered - PostTransaction() failed ")
		hBidFlag = false
	}

	//If there are some bids then do validations
	if hBidFlag == true {
		bid, err := JSONtoBid(Avalbytes)
		if err != nil {
			return nil, errors.New("BuyItNow() : JSONtoBid Error")
		}

		// Check if BuyItNow Price > Highest Bid so far
		binP, err := strconv.Atoi(args[5])
		if err != nil {
			return nil, errors.New("BuyItNow() : Invalid BuyItNow Price")
		}

		hbP, err := strconv.Atoi(bid.BidPrice)
		if err != nil {
			return nil, errors.New("BuyItNow() : Invalid Highest Bid Price")
		}

		if hbP > binP {
			return nil, errors.New("BuyItNow() : Highest Bid Price > BuyItNow Price - BuyItNow Rejected")
		}
	}

	// Close The Auction -  Fetch Auction Object
	Avalbytes, err = QueryLedger(stub, "AuctionTable", []string{args[0], "AUCREQ"})
	if err != nil {
		fmt.Println("BuyItNow(): Auction Object Retrieval Failed ")
		return nil, errors.New("BuyItNow(): Auction Object Retrieval Failed ")
	}

	aucR, err := JSONtoAucReq(Avalbytes)
	if err != nil {
		fmt.Println("BuyItNow(): Auction Object Unmarshalling Failed ")
		return nil, errors.New("BuyItNow(): Auction Object UnMarshalling Failed ")
	}

	//  Update Auction Status
	aucR.Status = "CLOSED"
	fmt.Println("BuyItNow(): UpdateAuctionStatus() successful ", aucR)

	Avalbytes, err = UpdateAuctionStatus(stub, "AuctionTable", aucR)
	if err != nil {
		fmt.Println("BuyItNow(): UpdateAuctionStatus() Failed ")
		return nil, errors.New("BuyItNow(): UpdateAuctionStatus() Failed ")
	}

	// Remove the Auction from Open Bucket
	keys := []string{"2016", aucR.AuctionID}
	err = DeleteFromLedger(stub, "AucOpenTable", keys)
	if err != nil {
		fmt.Println("BuyItNow(): DeleteFromLedger(AucOpenTable) Failed ")
		return nil, errors.New("BuyItNow(): DeleteFromLedger(AucOpenTable) Failed ")
	}

	fmt.Println("BuyItNow(): Proceeding to process the highest bid ")

	// Convert the BuyITNow to a Bid type struct
	buyItNowBid, err := CreateBidObject(args[0:])
	if err != nil {
		return nil, err
	}

	// Reject the offer if the Buyer Information Is not Valid or not registered on the Block Chain
	buyerInfo, err := ValidateMember(stub, args[4])
	fmt.Println("Buyer information  ", buyerInfo, args[4])
	if err != nil {
		fmt.Println("BuyItNow() : Failed Buyer not registered on the block-chain ", args[4])
		return nil, err
	}

	tran := BidtoTransaction(buyItNowBid)
	fmt.Println("BuyItNow(): Converting Bid to tran ", tran)

	// Process the buy-it-now offer
	tranArgs := []string{tran.AuctionID, tran.RecType, tran.ItemID, tran.TransType, tran.UserId, tran.TransDate, tran.HammerTime, tran.HammerPrice, tran.Details}
	fmt.Println("BuyItNow(): Proceeding to process the  Transaction ", tranArgs)

	Avalbytes, err = PostTransaction(stub, "PostTransaction", tranArgs)
	if err != nil {
		fmt.Println("BuyItNow(): PostTransaction() Failed ")
		return nil, errors.New("CloseAuction(): PostTransaction() Failed ")
	}
	fmt.Println("BuyItNow(): PostTransaction() Completed Successfully ")
	return Avalbytes, nil
}

//////////////////////////////////////////////////////////////////////////
// Update the Auction Object
// This function updates the status of the auction
// from INIT to OPEN to CLOSED
//////////////////////////////////////////////////////////////////////////

func UpdateAuctionStatus(stub shim.ChaincodeStubInterface, tableName string, ar AuctionRequest) ([]byte, error) {

	buff, err := AucReqtoJSON(ar)
	if err != nil {
		fmt.Println("UpdateAuctionStatus() : Failed Cannot create object buffer for write : ", ar.AuctionID)
		return nil, errors.New("UpdateAuctionStatus(): Failed Cannot create object buffer for write : " + ar.AuctionID)
	}

	// Update the ledger with the Buffer Data
	keys := []string{ar.AuctionID, ar.ItemID}
	err = ReplaceLedgerEntry(stub, "AuctionTable", keys, buff)
	if err != nil {
		fmt.Println("UpdateAuctionStatus() : write error while inserting record\n")
		return buff, err
	}
	return buff, err
}
