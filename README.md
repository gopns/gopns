Go Push Notification Service (GPNS)
===

GPNS provides scalable and efficient a mass push notificaiton service built on top of Amazon's Simple Notification Service. The project provides REST end points for registering Device IDs with arbitrary metadata. These can then be used to send a push notification to all or a segment of the user base based on the metadata. This project is written in the Go Programming language and is designed to be resource efficient. In addition it uses Amazon Simple Queue Service in order to distribute work load out to multiple instances if you have them available. 

### Build from Source
* Install [Go 1.2](http://golang.org/doc/install#install)
* Setup GOPATH
* Add GOPATH/bin to PATH
* go get gpns
* go install gpns
* run gpns

```bash
mkdir ~/go
export GOPATH=~/go
export PATH=$PATH:$GOPATH/bin
go get github.com/usmanismail/gpns
go install github.com/usmanismail/gpns
gpns -awsConfig=[PATH-TO-AWS-CONFIG] -baseConfig=[PATH-TO-BASE-CONFIG]
```

### Configuration 

__Base Configuration__ ($GOPATH/src/github.com/usmanismail/gpns/config/base.conf)
```bash
[default]
#Port for the REST server to run on
port=8080
```
__AWS Configuration__ ($GOPATH/src/github.com/usmanismail/gpns/config/aws.conf)
```bash
[default]
#AWS User ID and Secret
id=ASJDKLADKNDKLNDKLSN
secret=SDBSDBSDNMBSNBDNbsndbsnbdsnBDSMNBsnbdnsdnsdnbs

#Name of Platform Applications, comma seperated list (e.g. Test1,Test2)
platform-applications=Test1,Test2

#Add one section for each Application configured above
[Test1]
#The Application ARN
arn=arn:aws:sns:us-east-1:238699486533:app/GCM/Test
#Which region is the application running in?
region=us-east-1

[Test2]
arn=arn:aws:sns:us-east-1:238699486533:app/GCM/Test
region=us-east-1
```

