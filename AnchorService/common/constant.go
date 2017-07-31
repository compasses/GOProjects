package common

import "time"

const (
	// Limits for keeping inputs from flooding our execution
	INMSGQUEUE_HIGH = 1000
	INMSGQUEUE_MED  = 500
	INMSGQUEUE_LOW  = 100

	DBSTATE_REQUEST_LIM_HIGH = 200
	DBSTATE_REQUEST_LIM_MED  = 50

	// Replay
	INTERNAL_REPLAY = 1
	NETWORK_REPLAY  = 2
	TIME_TEST       = 4 // Checks the time_stamp;  Don't put actual hashes into the map with this.
	REVEAL_REPLAY   = 8 // Checks for Reveal Entry Replays ... No duplicate Entries within our 4 hours!

	ADDRESS_LENGTH = 32 // Length of an Address or a Hash or Public Key
	// length of a Private Key
	SIGNATURE_LENGTH     = 64    // Length of a signature
	MAX_TRANSACTION_SIZE = 10240 // 10K like everything else?
	// Not sure if we need a minimum amount.  Set at 1 Factoshi

	// Database
	//==================
	// Limit on size of keys, since Maps in Go can't handle variable length keys.

	// Wallet
	//==================
	// Holds the root seeds for address generation
	// Holds the latest generated seed for each root seed.

	// Block
	//==================
	MARKER                  = 0x00                       // Byte used to mark minute boundries in Factoid blocks
	TRANSACTION_PRIOR_LIMIT = int64(12 * 60 * 60 * 1000) // Transactions prior to 12hrs before a block are invalid
	TRANSACTION_POST_LIMIT  = int64(12 * 60 * 60 * 1000) // Transactions after 12hrs following a block are invalid

	//Entry Credit Blocks (For now, everyone gets the same cap)
	EC_CAP = 5 //Number of ECBlocks we start with.
	//Administrative Block Cap for AB messages

	//Limits and Sizes
	//==================
	//Maximum size for Entry External IDs and the Data
	HASH_LENGTH = int(32) //Length of a Hash
	//Length of a signature
	//Prphan mem pool size
	//Transaction mem pool size
	//Block mem bool size
	//MY Process List size

	//Max number of entry credits per entry
	//Max number of entry credits per chain

	COMMIT_TIME_WINDOW = time.Duration(12) //Time windows for commit chain and commit entry +/- 12 hours

	//NETWORK constants
	//==================
	VERSION_0               = byte(0)
	FACTOMD_VERSION         = 4002005
	MAIN_NETWORK_ID  uint32 = 0xFA92E5A2
	TEST_NETWORK_ID  uint32 = 0xFA92E5A3
	LOCAL_NETWORK_ID uint32 = 0xFA92E5A4
	MaxBlocksPerMsg         = 500
)

var ZERO_HASH = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
