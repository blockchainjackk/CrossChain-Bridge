package rpc

// GetAccountBalanceResult models the account data from the getbalance command.
type GetAccountBalanceResult struct {
	AccountName             string  `json:"accountname"`
	ImmatureCoinbaseRewards float64 `json:"immaturecoinbaserewards"`
	ImmatureStakeGeneration float64 `json:"immaturestakegeneration"`
	LockedByTickets         float64 `json:"lockedbytickets"`
	Spendable               float64 `json:"spendable"`
	Total                   float64 `json:"total"`
	Unconfirmed             float64 `json:"unconfirmed"`
	VotingAuthority         float64 `json:"votingauthority"`
}

type GetBalanceResult struct {
	Balances                     []GetAccountBalanceResult `json:"balances"`
	BlockHash                    string                    `json:"blockhash"`
	TotalImmatureCoinbaseRewards float64                   `json:"totalimmaturecoinbaserewards,omitempty"`
	TotalImmatureStakeGeneration float64                   `json:"totalimmaturestakegeneration,omitempty"`
	TotalLockedByTickets         float64                   `json:"totallockedbytickets,omitempty"`
	TotalSpendable               float64                   `json:"totalspendable,omitempty"`
	CumulativeTotal              float64                   `json:"cumulativetotal,omitempty"`
	TotalUnconfirmed             float64                   `json:"totalunconfirmed,omitempty"`
	TotalVotingAuthority         float64                   `json:"totalvotingauthority,omitempty"`
}
