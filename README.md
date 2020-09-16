# README

Given a network map from node B and a network map from node A, combine these to the shortest path from A to each node

// a --> b --> c --> d
//			   e --> f
// g --> f

- Take new network map
- Calculate shortest path from yourself to all other clients
- Compare those paths with the *cached* paths of your own network map
- If you find a shorter path, replace your path with the new path
- 