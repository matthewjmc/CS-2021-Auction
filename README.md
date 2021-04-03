
# CS-2021-Auction Project

There will be two parts for this branch. First, server-side caching and auctioning logic for the system while the other would be the database transactions and SQL statements which will be used to interact with the database.

=====================================================================================================================================================

These functions below are the core of the server's process which returns boolean and uint64 as reporting error codes.

#### MakeBidMain(u *UserHashTable, h *AuctionHashTable, uid uint64, targetid uint64, placeVal uint64).
: MakeBidMain is called when the "command" is "bid".
: code 0 . the bid has been made and updated properly, also returning "true".
: code 1 , the user has not been found within the system, also returning "false".
: code 2 , the auction has not been found within the system, also returning "false".
: code 3 , the auction has not been updated due to losing conditions, also returning "false".

#### CreateUserMain(h *UserHashTable, uid uint64, name string).
: createUserMain is called when the "command" is "User" or "user".
: code 0 , the user has not been found within the system, creating new user object while also returning "true".
: code 1 , the user has been found in the system, also returning "false".

#### CreateAuctionMain(U *UserHashTable, A *AuctionHashTable, uid uint64, aid uint64, initial uint64, step uint64, duration time.Duration, itemName string). 
: createAuctionMain is called when the "command" is "Auction" or "auction".
: code 0, auction has not been found within the system, creating new auction object while also returning "true".
: code 1 , auction has been found in the system, also returning "false".

=====================================================================================================================================================

While these functions can be used to initiate database server with relationships and tables for data storage. In this case, we are using MySQL from MariaDB to work as a data storage in case when main logic server failed.

#### DatabaseInit() is called to create the database named "auction_system" at your specified database server.

#### InsertAuctionToDB( auction Auction ) is used to insert auctions into the database. 
#### InsertUserToDB( user User ) is used to insert users into the database. 
#### InsertBidToDB( bid Bid  ) is used to insert bidding transactions into the database. 
#### UpdateAuctionInDB( auction Auction ) is called to update the bidded auction within the database.
#### UserFromDBToHash() and AuctionFromDBToHash() are used to retrieve records from the database in needed.

=====================================================================================================================================================

#### Code Optimization and Time Usage regarding Localhost:3306 testing.
This section composes of time analysis upon the business logic and database transaction part which will affect the system as a whole. Consideration of time-efficient and preciseness are a priority for the third milestone.

#### For this part, time usage for each functions are used to consider the possibility of actual launched server required resources. So, time is captured to analyze for each possible case.

- Time required to complete CreateUserMain() which composed of user struct creation and insertion into both the server cache and the locally hosted database is approximately 1-2 milliseconds while implementing it to be on the actual MariaDB database will require approximately 50-70 milliseconds with consideration of network latency and local computation time.

- Time required to complete CreateAuctionMain() which considers auction struct creation and inserting it into the server-cache and the database takes around 3 milliseconds while taking ~60 milliseconds to complete the database statement and inserting it into the database itself.

- Lastly, time required to complete MakeBidMain(). The provided time includes creating bid struct along with checking whether the newly placed bid would win the auction or not and either update the targeted auction for its new winner within the server cache and the database or ignore the updates.

    - For the bidding transaction, it takes around 1-3 milliseconds to create and insert it into the locally hosted database, while it instead takes around 57-63 milliseconds to apply it onto the actual running server.
    - For the auction updating section, 1.5-2.5 milliseconds is required to completely update the auction onto both the server-cache and the database while it takes 60-70 milliseconds to the actual database server.
    - If the auction is not being updated due to the losing condition, it takes around 500-600 microseconds(local) or 60 milliseconds(actual server) just to ignore both the query and server-cache updates.
    
    - In summary, for this section for the whole MakeBidMain() function, it takes within 2 to 5 milliseconds as its maximum (locally) or 62-65 milliseconds (to the server).





