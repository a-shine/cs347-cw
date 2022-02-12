# CS347 coursework
Building a custom overlay network on top of the a-shine/butter framework (unstructured p2p application framework) focused on implementing group based persistent data management techniques to improve fault-tolerance (specifically information availability) on high churn networks.

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

To get the latest changes from teh butter framework you can run `go get github.com/a-shine/butter` (warning there may be breaking changes)