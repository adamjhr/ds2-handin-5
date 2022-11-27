# ds2-handin-5

To start, run all the replicas: ``go run server/server.go -port=<PORT>``
Run at least 3 of these, with different ports.

Then the clients and frontends can be created. Each client is connected to a frontend. Is frontend is connected to all replicas. Each frontend is only connected to one client.

To create the frontend: ``go run frontend/frontend.go -port=<PORT> <REPLICA PORT> <REPLICA PORT> ...``
Replace <PORT> with the port of the frontend, and give all the ports of the replicas as arguments
Do this for each frontend.

To create the client: ``go run client/client.go -port=<PORT> -receiver=<FRONTEND PORT>``
Replace <PORT> with the port of the client, connect it to a frontend by giving the port of the frontend in as second flag.
