How to run our program 101

1. Launch three different nodes in three different terminals. A unique integer id and a unique port
must be specified for each node. We recommend the following:
- go run node.go -name 1 -port 5400
- go run node.go -name 2 -port 5401
- go run node.go -name 3 -port 5402

2. The nodes should be set up in a peer-to-peer architecture, by connecting them to each other. This
can be done by typing "connect " + the port you wish to connect to. If you followed the example above,
you should type the following for each node.
Node 1:
- connect 5401
- connect 5402

Node 2:
- connect 5400
- connect 5402

Node 3:
- connect 5400
- connect 5401

3. You can request for a node to gain access to the crittical section by typing "request". The program
uses the Ricart-Agrawala algorithm to determine who gains access. When a node has entered the critical,
it will simply print "I'm in the critical section". Afterwards the program will sleep for 10 seconds to
simulate work being done, and then it will exit the critical section by printing "Exiting the critical 
section". There are many print statements during the execution of the program that should help you follow
whats happening.  