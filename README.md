# Storage serive
# Zhoucheng Li (zhouchel)

###### You will need to install Beego and the Bee dev tool:
```
$ go get github.com/astaxie/beego
$ go get github.com/beego/bee
```
For convenience, you should add $GOPATH/bin to your$PATH environment variable.

###### Start web service
```
$ cd src/storage
$ bee run
```

###### Open Web Broswer
visit localhost:8080

###### Register and Login!



###### API usage
####### upload
```
import requests

r = requests.post('http://localhost:8080/v1/file', files={'files': open('test2.py', 'rb')})
print r.text
```
