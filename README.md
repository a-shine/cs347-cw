# CS347 coursework

## Background

Building a custom overlay on top of the [a-shine/butter](https://github.com/a-shine/butter) framework (unstructured p2p application framework) focused on implementing group based persistent data management techniques to improve fault-tolerance (specifically information availability) on high churn networks.

## Try demo
To try a demo of the PCG persistent storage management overlay in action:
1. Make sure to have Go (1.17) installed
2. Clone the repository
3. Run the demo blog/wiki cli program with:
    ```bash
    go run demo.go
    ```
4. Now start several nodes in different terminal instances, try adding and retrieving information from the various nodes. Watch as information persists beyond the existence of any single node.