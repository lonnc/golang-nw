golang-nw
=========

Calling golang application in node-webkit to get a native looking application.

Please note that this is fairly ugly code that I've used to test golang & node-webkit:

1. I'm new to node so please shout if you see problems.
2. Passing a message from golang back to node is done through regex on golang's stdout (yuck).
3. Failure cases should golang fail isn't really handled.

However, it may provide a starting point for others to improve on.

Download it
-----------
go get github.com/lonnc/golang-nw/cmd/golang-nw-build

Build your app
--------------
go install .\src\github.com\lonnc\golang-nw\cmd\example

Create your nw zip file and run it with node-webkit
--------------
.\bin\golang-nw-build -app=.\bin\example.exe -name="My Application" -out="myapp.nw"
nw.exe myapp.nw
