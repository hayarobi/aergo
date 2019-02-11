package mempool

import (
	"github.com/aergoio/aergo-actor/actor"
	"github.com/aergoio/aergo/message"
	"github.com/aergoio/aergo/types"
)

type TxVerifier struct {
	mp *MemPool
}

func NewTxVerifier(p *MemPool) *TxVerifier {
	return &TxVerifier{mp: p}
}

//Receive actor message
func (s *TxVerifier) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *message.MemPoolPut:
		tx := msg.Tx
		var err error
		if s.mp.exist(tx.GetHash()) != nil {
			err = types.ErrTxAlreadyInMempool
		} else {
			tx := types.NewTransaction(tx)
			err = s.mp.verifyTx(tx)
			if err == nil {
				err = s.mp.put(tx)
			}
		}
		context.Respond(&message.MemPoolPutRsp{Err: checkToBlame(err), Sender: msg.Sender})
	}
}

func checkToBlame(err error) error {
	if err == nil {
		return nil
	}
	switch rawerr := err.(type) {
	case message.BlamableError:
		return rawerr
	default:
		blameerr, found := faultMap[rawerr]
		if found {
			return blameerr
		}
		return err
	}
}

var (
	faultMap map[error]error
)

func init() {
	faultMap = make(map[error]error)
	//Severe: verification failures
	faultMap[types.ErrSignNotMatch] = message.NewBlamableWrapper(message.Severe, types.ErrSignNotMatch)
	faultMap[types.ErrTxHasInvalidHash] = message.NewBlamableWrapper(message.Severe, types.ErrTxHasInvalidHash)
	faultMap[types.ErrTxInvalidChainIdHash] = message.NewBlamableWrapper(message.Severe, types.ErrTxInvalidChainIdHash)
	faultMap[types.ErrTxInvalidType] = message.NewBlamableWrapper(message.Severe, types.ErrTxInvalidType)

	//Big: violate formats
	faultMap[types.ErrCouldNotRecoverPubKey] = message.NewBlamableWrapper(message.Big, types.ErrCouldNotRecoverPubKey)
	faultMap[types.ErrSignNotMatch] = message.NewBlamableWrapper(message.Big, types.ErrSignNotMatch)
	faultMap[types.ErrTxFormatInvalid] = message.NewBlamableWrapper(message.Big, types.ErrTxFormatInvalid)
	faultMap[types.ErrTxInvalidAccount] = message.NewBlamableWrapper(message.Big, types.ErrTxInvalidAccount)
	faultMap[types.ErrTxInvalidAmount] = message.NewBlamableWrapper(message.Big, types.ErrTxInvalidAmount)
	faultMap[types.ErrTxInvalidPrice] = message.NewBlamableWrapper(message.Big, types.ErrTxInvalidPrice)
	faultMap[types.ErrTxInvalidPayload] = message.NewBlamableWrapper(message.Big, types.ErrTxInvalidPayload)
	faultMap[types.ErrTxInvalidRecipient] = message.NewBlamableWrapper(message.Big, types.ErrTxInvalidRecipient)
	faultMap[types.ErrTxInvalidSize] = message.NewBlamableWrapper(message.Big, types.ErrTxInvalidSize)

	//Normal: violate rules
	faultMap[types.ErrExceedAmount] = message.NewBlamableWrapper(message.Normal, types.ErrExceedAmount)
	faultMap[types.ErrInsufficientBalance] = message.NewBlamableWrapper(message.Normal, types.ErrInsufficientBalance)
	faultMap[types.ErrLessTimeHasPassed] = message.NewBlamableWrapper(message.Normal, types.ErrLessTimeHasPassed)
	faultMap[types.ErrMustStakeBeforeUnstake] = message.NewBlamableWrapper(message.Normal, types.ErrMustStakeBeforeUnstake)
	faultMap[types.ErrMustStakeBeforeVote] = message.NewBlamableWrapper(message.Normal, types.ErrMustStakeBeforeVote)
	faultMap[types.ErrNameNotFound] = message.NewBlamableWrapper(message.Normal, types.ErrNameNotFound)
	faultMap[types.ErrSameNonceAlreadyInMempool] = message.NewBlamableWrapper(message.Normal, types.ErrSameNonceAlreadyInMempool)
	faultMap[types.ErrShouldUnlockAccount] = message.NewBlamableWrapper(message.Normal, types.ErrShouldUnlockAccount)
	faultMap[types.ErrTooSmallAmount] = message.NewBlamableWrapper(message.Normal, types.ErrTooSmallAmount)
	faultMap[types.ErrTxNotFound] = message.NewBlamableWrapper(message.Normal, types.ErrTxNotFound)
	faultMap[types.ErrWrongAddressOrPassWord] = message.NewBlamableWrapper(message.Normal, types.ErrWrongAddressOrPassWord)

	//Small: concurrent receive or etc
	faultMap[types.ErrTxAlreadyInMempool] = message.NewBlamableWrapper(message.Small, types.ErrTxAlreadyInMempool)
	faultMap[types.ErrTxNonceTooLow] = message.NewBlamableWrapper(message.Small, types.ErrTxNonceTooLow)
	faultMap[types.ErrTxNonceToohigh] = message.NewBlamableWrapper(message.Small, types.ErrTxNonceToohigh)
}
