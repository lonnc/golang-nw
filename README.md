golang-nw
=========

Calling golang application in node-webkit to get a native looking application.

Please note that this is fairly ugly code that I've used to test golang & node-webkit:

1. I never used node prior to yesterday - so there are likely to be problems with it.
2. Passing a message from golang back to node is done through regex on golang's stdout (yuck).
3. Handling failure cases should golang fail isn't really handled.

However, it may provide a starting point for others to improve on.

Build your app
--------------
go build -o myapp.exe main.go

Create your nw zip file and run it with node-webkit
--------------
go run build.go && nw.exe myapp.nw
