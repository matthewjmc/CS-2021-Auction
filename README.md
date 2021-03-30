
# CS-2021-Auction Project

type Data struct {
	command      string
	uid          uint64
	fullname     string
	aid          uint64
	itemname     string
	biddingValue uint64
	biddingStep  uint64
	duration     time.Duration
}

// Data struct is initiated to represent the requirements for each functions.

// current mainTimeline if-else if conditions are { User, Auction, bid, searchAuction, deleteAuction, searchUser, deleteUser }


Calling users -> create new user and insert it into the storage system , requiring userID and user's full name.
Calling Auctions -> create new auction into the system, requiring userID , auctionID, first initial bidding price and bidding step condition for each bids.
Calling bid -> create new bidding transaction onto a targeting auction, thus, it requires userID and their bidding price along with targeting auctionID.
Calling searchXXXXX -> will display the statement providing whether the XXXX asked for is within the system or not, requiring the object's ID ( auction/user )
Calling deleteXXXXX -> will delete the object XXXXX within the system.

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Three composing functions within the main processing function "makeTimeline()"
: these functions listed below were initially created to be called and running asynchronously from the main thread.


These functions below are the core of the server's process which returns boolean and uint64 as reporting error codes.

### MakeBidMain(u *UserHashTable, h *AuctionHashTable, uid uint64, targetid uint64, placeVal uint64) 
: MakeBidMain is called when the "command" is "bid".
: code 0 . the bid has been made and updated properly, also returning "true"
: code 1 , the user has not been found within the system, also returning "false"
: code 2 , the auction has not been found within the system, also returning "false"

### CreateUserMain(h *UserHashTable, uid uint64, name string) 
: createUserMain is called when the "command" is "User" or "user".
: code 0 , the user has not been found within the system, creating new user object while also returning "true"
: code 1 , the user has been found in the system, also returning "false"

### CreateAuctionMain(U *UserHashTable, A *AuctionHashTable, uid uint64, aid uint64, initial uint64, step uint64, duration time.Duration, itemName string) 
: createAuctionMain is called when the "command" is "Auction" or "auction".
: code 0, auction has not been found within the system, creating new auction object while also returning "true"
: code 1 , auction has been found in the system, also returning "false"

