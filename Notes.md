# Maintaining Information Availability in a High Churn Unstructured Peer-to-Peer System

^^ Title subject to change

## TODO

* Change find peers stuff to be small request, then if accepted then join the network with the group Digest stuff
* change leader election function to avoid reliance on high ID stuff
* information retrieval is a request service (broadcast to group the existence of this request) (add to group struct a series of requests that exist)
* 



## Keywords

We need consistency with our terminology.

* Node
  * Node and Peer have been used interchangeably up until this point. 
  * For our paper we shall only use the term Node due to it's greater usage in general academic papers relating to distributed and p2p systems.
* Health??
  * evaluates a group status, i.e. Too few or too many nodes in a group
* Client
  * in traditional web server sense
* Churn
  * Refers to high number of joining and leaving participants in the network.
* Leader?
  * A node is elected at a group Leader to perform group actions such as handling requestes to clients
* Suitability
  * Basically the availability of a node to take more data and how good it is at node stuff




## Background

what area are we in

what is a pure P2P network

what is high churn

what is Information Availability

## Problem that exists

Current academia limited in this pure environment due to industry favouring partially centralised approaches for performance.

$\not\exist$ systems for (PURE) unstructured networks, that can perform in high churn whilst attempting to maintain minimal network traffic.

(Not sure how to say this, but we need to keep in mind that reducing message complexity is paramount. Trivial solutions exist but with extreme message complexities make them entirely unsuitable and void).

Key point is availability, and this will be the measure of success not message complexity. 

## Our contribution

We take pre-existing approach of peer content groups, from JXTA paper, show it's limitations (centralised advertising "super-peers") and it's hierarchical approach.

We then propose an algorithm that works on an entirely unstructured network and test it to show it performs under high churn situations.  (highly based on JXTA paper but with key differences for unsturctured systems and efficiency).

## Structure of paper 





# Research

## Reliable Content Distribution Based on Peer Groups 

https://www.scirp.org/pdf/IJIDS_2014042910144755.pdf

### Key points:

* Peer content groups (PCG)
  * Self contained management systems
  * Responsibility lies with the group members to maintain the group
  * Super-peers for work advertisements
    * When a node is needed, an advert is put out onto a super peer
    * Nodes query the advertisement board when looking for work. This helps reduce the message complexity. 
    * This is NOT unstructured. It creates potential for single points of failure. 
      * The paper implies duplication of super-peers will protect, but ideologically insufficient as if super-peers go down the entire network breaks
  * Group interface is identical to a singular client interface, if there are changes to the group then the client querying the group will have no knowledge. The underlying network below groups should be TRANSPARENT to the user.
* Hierarchical Approach
  * May be a requirement from JXTA but not entirely sure.
  * gives you abstract groups that don't seem to serve so much purpose
  * seems just kinda stupid

### What we can take from it

Good idea about PCG and how they maintain their group health with heartbeats and leader elections.

Very bad structured implementation with super-peers and stuff like that.

Heirarchies are also stupid.

TLDR; good PCG everything else fucking stupid. :)

## Dynamic Model-Driven Replication System

https://ieeexplore.ieee.org/stamp/stamp.jsp?tp=&arnumber=1540492

### Key points

* No groups (must use Replica locator which is slow)
* At some given interval, look at network and see if the data should be replicated anymore (using replica locator)
* Neccesity to replicate defined dynamically based on network status
* Find best hosts for candidates if required and give them the data

### What we can take from it

Dynamic methodologies for replication

Our initial approach will just $r=3$ for simplicity. If $r$ could be modified dynamically based on live features of the network this may be desiriable. 

It uses factors such as node location also for finding best candidate nodes.

This also is entirely unstructured which is nice.

# Our Proposed Solution

## Key points

What we are building is an *Overlay Network*, much lke Chord or Kadmilia on structured peer-to-peer architectures, but in this case we present a persistent storage overlay on an unstructured architecure - this allow for more decentralisation. In an overlay network we interface with the underlying nodes and add the extra functionalities that allow for persistent storage on the network.

A Peer Content Group (PCG) implementation without any hierarchies.

Each PCG will maintain some block of data or a file (TBD...)

Each group is responsible for it's own "heath", i.e. the status of the group(the number of non-faulty nodes) is at the desired number $r$. ($r$ can be variable should be implement a dynamic replication system)

As groups should be in agreement of who is in the group (may require local group view and an agreement view as described in JXTA paper (first implementaiton may not need this as not dynamic and small)) it then means the leader election problem becomes trivial.

​	 Leader = group member of highest ID.



## Group Maintenance

As stated earlier, the responsibility lies with the group to ensure group health. More specifically, it is the group leader's responsibility to 'fix' the situation if the group is in an "unsafe" state.

The primary "unsafe" state is when there are too few members in the group. In this case, the leader will work to find new group members and recruit them. This ideology is used to help reduce message complexity, meaning nodes are only used when they are needed. 

If a leader wishes to add a member to the group it will consult it's known hosts list (contained in node) and ask for their availability. If they are the most suitable they will be requested to join a group.

Notion of node suitability comes from a calculation performed internally to each node, which can then be queried by any other node. The calculation is likely to include factors such as; storage availability; up-time; and network quality. (Note: this is highly unsafe given byzantine failure or malicious operation)

## Store

The node wishing to store a piece of data will create a new group and instantaneously fall into the state as it's only member and therefore leader attempting to resolve the unsafe state the group is in.

## Retrieve

When a client(in the traditional web service sense) wishes to make a request, it must join the network. The existence of this node is then propagated through butter nodes inbuilt known-Hosts functionality. 

As a node in the system it does not automatically join the group containing the data which it is requesting. This is for the purposes of limiting group size to the optimal level, which helps reduce message complexity.

To actually retrieve the data, the node simply performs a breadth first search through the network querying nodes if they contain some data it is looking for. ==(Alex please advise on the details of this section)==

When a retrieval request is received by a group member (note: this can be any group member). It will broadcast this request to every node in the group. Thus ensuring that the request is not lost in the case of node failure. 

The leader of the group will serve the requested data back to the client. When that is complete it will then inform all group members that it has correctly served the request and all nodes can forget about it. 

In the case of  leader node failure, the next leader can then pickup the request and begin re-sending it. (May need some way to inform the client of data send restart ==please advise==)



## Stuff we are NOT doing

We are not implementing a system to automatically remove nodes from a group and replace them with a new one.

Not using locatino to help sort it.

Not doing NAT traversal systems (I think? just not worth it potentially)

# Server Behaviours Required (Endpoints)

==People please help me fill these out==

## Heartbeat

 /heartbeat

##### Description:

###### Data-in:

###### Data-out:

---

## Get Suitability

 /suitability

##### Description:

###### Data-in:

###### Data-out:

---

## Request to Join Group

 /joinrequest

##### Description:

###### Data-in:

###### Data-out:

---

## Request data generally??

---

## Group Request 

inform members of the new request 

---

## Request Complete

inform members a request has been completed and they can now forget it. 

# Questions to keep Nick Brain Happy

Woop woop it's empty.

Keywords: data redundancy, data consistency/availability, unstructured peer-to-peer, decentralised persistent data

Background The main idea is to look at maintaining data consistently (i.e. no loss of data) on a distributed system with an unstructured peer-to-peer architecture.

The main approach to maintaining consistent data is to have some form of data replication. Each data block can be replicated to a certain data replication degree $n$.

Finding efficient ways of handling data redundancy in such a way that data is persistent while not gaining an unmanageable overhead of storage (replicated data) and computations (network traffic to maintain the copied data) is a non trivial problem.

The resulting system is rendered more fault-tolerant as it is decentralised (no central point of failure by its very nature) and maintains persistent data despite node failure.

Why (pure) unstructured peer-to-peer architecture? Lots of seemingly decentralised platforms are based on on structured peer-to-peer architecture which ar arguably less fault-tolerant. They claim to be decentralised but yet have some form of hierarchy within the resulting system. Maintaining teh structure can centralised and there is some node inter-dependability. While this is more robust that a non-distributed system, a more powerful architecture would be an unstructured peer-to-peer architecture, where every peer has as much impact on the network as another. Managing persistent data on a pure unstructured system becomes more of a challenge as the common approach of DHT is not available to us. (DHTs require a strcurured peer to peer network and some form of centralised index of bootstrap nodes.)

Interesting projects in the space

libp2p/IPFS Gnutella Kademlia/Chord PlanetP

Basically you can think on a group abstraction level

a group can be in 3 states

wanting more members
repetedly call out for new members
godilocks state - just the right amount of people in the group
too many people in group
cut members from the participant list
group leader is determined on a per request basis

try maximise groups with diverse IP address - maximise across subnetworks

Vhange blocks to groups and add a participants field


**testing thoughts**
/*
Testing
* Other thoughts:
  * What are we doing with nodes that leave the network and join again? Currently we take them as fresh (not holding any data) when they rejoin. Is this what we decide?
  * If nodes are "reinitialised" on rejoining, then we can simulate this easily just with new nodes that contain no data. (only slight difference that I can think of is the underlying butter network known host list churn but that's fine?)

* General Notes:
  * Don't have to implement all proper functionalilty to every node, such as the ones for just making requests
    they can be seen as nodes outside the network that don't do anythign properly just make requests.
  * How do we want to go about adding data to the network? does each node start with some random amount of data? then don't add any later on?
  * Could make nodes that don't add any, and they're just good helpful bois.

* Global settings:
  * average lifetime of nodes + additional randomized lifespan + auto kill.
  * % of deaths be graceful or ungraceful.
  * have a delayed communication response to check if nodes are being removed from PCGs if they are only slow not dead?
  * Rate of requests?
  * Number of nodes on the netowrk?
  * Number of nodes making requests? ()
  * vary rate of data addition or max data added?
  * vary mean and s.d. of node capacities

* Functions:
  * Kill myself function:
    * Uses config settings to generate pseudo random lifetime and death style
    * will kill the node when time is up!

  * Add data to network
    * Nodes generate random strings of data and send it to other nodes
    * Add UUID (and the data expected for validation purposes) to global array for the sake of retrieve

  * Retreive data
    * The main function for evaluating availability
    * requests data it knows should be on the network
    * report success rate
    * performed at vraying request rates?

  * Network evaluation
    * Everytimme a tcp request is made add to some global counter
    * Used to model network traffic
    * can be used to compare heartbeat to gossip should we choose to do that.

  * Generic Evaluation:
    * work out % availability under different environments/churn rate/scenarios
      */
