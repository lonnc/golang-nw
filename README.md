golang-nw
=========

Calling golang application from node-webkit to get a native looking application.

Please note that this is fairly ugly code that I've used to test golang & node-webkit:

1. I'm new to node so please shout if you see problems.
2. Passing a message from golang back to node is done through regex on golang's stdout (yuck).
3. Failure cases should golang hit a problem aren't really handled.

However, it may provide a starting point for others to improve on.

Dependencies
------------
Download node-webkit from https://github.com/rogerwang/node-webkit/#downloads.

Instructions
------------

Go get golang-nw:

    go get github.com/lonnc/golang-nw/cmd/golang-nw-build


Build your app:

    go install .\src\github.com\lonnc\golang-nw\cmd\example


Wrap it in node-webkit:

    .\bin\golang-nw-build.exe -app=.\bin\example.exe -name="My Application" -out="myapp.nw"
    

Finally execute node-webkit with the myapp.nw generated above as a parameter:

    nw.exe myapp.nw


You will probably want to create you own build script based on 
https://github.com/lonnc/golang-nw/blob/master/cmd/golang-nw-build/build.go
so you can control toolbar visibility, window dimensions etc.
