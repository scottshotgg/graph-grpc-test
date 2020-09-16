A <---> B <---> C

1. A joins the network (makes the network)
2. B joins the network (joins A)
3. C joins the network (joins B)

A and B can still address eachother and know about each other:

1. NetworkMap from A's perspective:
    A: {
        info: ...
        netMap: {}
    }

2. NetworkMap is exchanged and updated between nodes A and B:
    A: {
        info: ...
        netMap: {
            "B": {
                info: ...
                netMap: {}
            }
        }
    }

    B: {
        info: ...
        netMap: {
            "A": {
                info: ...
                netMap: {}
            }
        }
    }

3. C joins the network through B; NetworkMap is exchanged between B and C

    B: {
        info: ...
        netMap: {
            "A": {
                info: ...
                netMap: {}
            }
        }
    }

    C: {
        info: ...
        netMap: {
            "B": {
                info: ...
                netMap: {
                    "A": {
                        info: ...
                        netMap: {}
                    }
                }
            }
        }
    }

A client of C wants to send a message to A:
msg ---> C ---> B ---> A

Client sends msg to C:
{
    addr: "b.a",
    msg: "something here"
}

C will peel the address and send to B:
{
    addr: "a",
    msg: "something here"
}

B will recieve the message and send it to a:
{
    msg: "something here"
}