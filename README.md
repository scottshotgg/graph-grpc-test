# README

This repo mainly served as a test for a node network utilizing an underlying graph that is distributed upon change with a relative shortest-path computation. Also playing with some graphs for some underlying ideas

Right now the example located in `cmd/exchange` and `docker-compose.yml` creates a three node network where each node starts with varying intial peer connections as below:
- node1: _no initial peers_
- node2: _node1_
- node3: _node1 and node2_

On start the node will attempt to replicate the rest of the network map. This involves exchanging graphs with each of the given intial peers to discover the rest of the nodes. On exchange, each node will send back its full graph which may introduce previously unknown nodes, which are then exchanged with as well allowing the node to recursively discover previously unknown peers. On each discovery, the node adds verticies for each new peer and an accompanying arc weight based on latency. This process goes on until the point where discovery has not produced any new results, which means the entire _known_ network has been replicated. Ideally, the node has effectively created a _complete graph_ by establishing a connection to each and every other node (taking into account that every _other_ node will **also** have a connection to each and every node, which this _node_ will be aware of), which it can then use to recalculate the reduced graph (shortest path to every other node) and optimize the flow of traffic.




------
------


Given a network map from node B and a network map from node A, combine these to the shortest path from A to each node

// a --> b --> c --> d
//			   e --> f
// g --> f

- Take new network map
- Calculate shortest path from yourself to all other clients
- Compare those paths with the *cached* paths of your own network map
- If you find a shorter path, replace your path with the new path