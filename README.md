# LoadBalancer for Auction System

## Load Balance from CPUusage on backend servers 1 and 2

## Getstat
### Always run getstat file on both servers to continuously get the cpu usage and run getstat on load to retrieve usage and update the values onto redis with fixed keys

## Redis
### Implement redis to store auctionID and map it to specific address

## Main
### Main function will received every connections and send command to redis only then will it pass to reverse proxy to direct to the backend servers