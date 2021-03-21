
# CS-2021-Auction Project

type Data struct {
	command      string
	uid          uint64
	fullname     string
	aid          uint64
	itemname     string
	biddingValue uint64
	biddingStep  uint64
}

// Data struct is initiated to represent the requirements for each functions.

// current mainTimeline if-elseif conditions are { User, Auction, bid, searchAuction, deleteAuction, searchUser, deleteUser }


Calling users -> create new user and insert it into the storage system , requiring userID and user's full name.
Calling Auctions -> create new auction into the system, requiring userID , auctionID, first initial bidding price and bidding step condition for each bids.
Calling bid -> create new bidding transaction onto a targeting auction, thus, it requires userID and their bidding price along with targeting auctionID.
Calling searchXXXXX -> will display the statement providing whether the XXXX asked for is within the system or not, requiring the object's ID ( auction/user )
Calling deleteXXXXX -> will delete the object XXXXX within the system.

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Three composing functions within the main processing function "makeTimeline()"
: these functions listed below were initially created to be called and running asynchronously from the main thread.


These functions below are the core of the server's process.

### makeBidMain(u *UserHashTable, h *AuctionHashTable, report_price chan uint64, report_log chan string, uid uint64, targetid uint64, placeVal uint64)
: makeBidMain is called when the "command" is "bid".
: if some return value is required, Auction which has already been updated and stored at the same place should be the type .

### createUserMain(u *UserHashTable, report chan User, report_log chan string, uid uint64, name string)
: createUserMain is called when the "command" is "User" or "user".
: if some return value is required, User should be the type.

### createAuctionMain(u *UserHashTable, A *AuctionHashTable, auction chan Auction, report_log chan string, uid uint64, aid uint64, initial uint64, step uint64)
: createAuctionMain is called when the "command" is "Auction" or "auction".
: if some return value is required, Auction should be the type.

Note: report_xxx are channels used to return the result from their asynchronous processes back to the main process. 
      When those functions are called within the main thread, removal of report_xxx chan are required.