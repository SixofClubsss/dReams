package rpc

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/ybbus/jsonrpc/v3"
)

type wallet struct {
	UserPass   string
	idHash     string
	Rpc        string
	Address    string
	ClientKey  string
	Balance    string
	TokenBal   string
	TourneyBal string
	Height     string
	Connect    bool
	PokerOwner bool
	BetOwner   bool
}

var Wallet wallet
var logEntry = widget.NewMultiLineEntry()

func StringToInt(s string) int {
	if s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Println("String Conversion Error", err)
			return 0
		}
		return i
	}

	return 0
}

func addLog(t string) {
	logEntry.SetText(logEntry.Text + "\n\n" + t)
	logEntry.Refresh()
}

func SessionLog() *fyne.Container {
	logEntry.Disable()
	button := widget.NewButton("Save", func() {
		saveLog(logEntry.Text)
	})

	cont := container.NewMax(logEntry)
	vbox := container.NewVBox(layout.NewSpacer(), button)
	max := container.NewMax(cont, vbox)

	return max
}

func saveLog(data string) {
	f, err := os.Create("Log " + time.Now().Format(time.UnixDate))

	if err != nil {
		log.Println(err)
		return
	}

	defer f.Close()

	_, err = f.WriteString(data)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Log File Saved")
}

func DeroAddress(v interface{}) (address string) {
	switch val := v.(type) {
	case string:
		decd, _ := hex.DecodeString(val)
		p := new(crypto.Point)
		if err := p.DecodeCompressed(decd); err == nil {
			addr := rpc.NewAddressFromKeys(p)
			address = addr.String()
		} else {
			address = string(decd)
		}
	}

	return address
}

func SetWalletClient(addr, pass string) (jsonrpc.RPCClient, context.Context, context.CancelFunc) { /// user:pass auth
	client := jsonrpc.NewClientWithOpts(pre+addr+suff, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(pass)),
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)

	return client, ctx, cancel
}

func GetAddress() error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	var result *rpc.GetAddress_Result
	err := rpcClientW.CallFor(ctx, &result, "GetAddress")

	if err != nil {
		Wallet.Connect = false
		log.Println(err)
		return nil
	}

	address := len(result.Address)
	if address == 66 {
		Wallet.Connect = true
		log.Println("Wallet Connected")
		log.Println("Dero Address: " + result.Address)
		Wallet.Address = result.Address
		id := []byte(result.Address)
		hash := sha256.Sum256(id)
		Wallet.idHash = hex.EncodeToString(hash[:])
	} else {
		Wallet.Connect = false
	}

	return err
}

func GetBalance(wc bool) error { /// get wallet dero balance
	if wc {
		rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
		defer cancel()

		var result *rpc.GetBalance_Result
		err := rpcClientW.CallFor(ctx, &result, "GetBalance")
		if err != nil {
			Wallet.Connect = false
			log.Println(err)
			return nil
		}

		Wallet.Balance = fromAtomic(result.Unlocked_Balance)

		return err
	}
	return nil
}

func TokenBalance(scid string) (uint64, error) { /// get wallet token balance
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	var result *rpc.GetBalance_Result
	sc := crypto.HashHexToHash(scid)
	params := &rpc.GetBalance_Params{
		SCID: sc,
	}

	err := rpcClientW.CallFor(ctx, &result, "GetBalance", params)
	if err != nil {
		log.Println(err)
		return 0, nil
	}

	return result.Unlocked_Balance, err
}

func DreamsBalance(wc bool) { /// get wallet dReam balance
	if wc {
		bal, _ := TokenBalance(dReamsSCID)
		Wallet.TokenBal = fromAtomic(bal)
	}
}

func TourneyBalance(wc, t bool, scid string) { /// get tournament balance
	if wc && t {
		bal, _ := TokenBalance(scid)
		value := float64(bal)
		Wallet.TourneyBal = fmt.Sprintf("%.2f", value/100000)
	}
}

func TourneyDeposit(bal uint64, name string) error {
	if bal > 0 {
		rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
		defer cancel()

		arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Deposit"}
		arg2 := rpc.Argument{Name: "name", DataType: "S", Value: name}
		args := rpc.Arguments{arg1, arg2}
		txid := rpc.Transfer_Result{}
		params := &rpc.SC_Invoke_Params{
			SC_ID:            TourneySCID,
			SC_RPC:           args,
			SC_TOKEN_Deposit: bal,
			Ringsize:         2,
		}

		err := rpcClientW.CallFor(ctx, &txid, "scinvoke", params)
		if err != nil {
			log.Println(err)
			return nil
		}

		log.Println("Tournament Deposit TX:", txid)
		addLog("Tournament Deposit TX: " + txid.TXID)

		return err
	}
	log.Println("No Tournament Chips")
	return nil
}

func GetHeight(wc bool) error { /// get wallet height
	if wc {
		rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
		defer cancel()

		var result *rpc.GetHeight_Result
		err := rpcClientW.CallFor(ctx, &result, "GetHeight")
		if err != nil {
			return nil
		}

		Wallet.Height = fmt.Sprint(result.Height)

		return err
	}
	return nil
}

func SitDown(name, av string) error { /// sit at holdero table
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	var hx string
	if av != "" {
		with := "_" + name + "_" + av
		hx = hex.EncodeToString([]byte(with))
	} else {
		out := "_" + name
		hx = hex.EncodeToString([]byte(out))
	}

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "PlayerEntry"}
	arg2 := rpc.Argument{Name: "address", DataType: "S", Value: Wallet.idHash + hx}
	args := rpc.Arguments{arg1, arg2}
	txid := rpc.Transfer_Result{}
	params := &rpc.SC_Invoke_Params{
		SC_ID:    Round.Contract,
		SC_RPC:   args,
		Ringsize: 2,
	}

	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Sit Down TX:", txid)
	addLog("Sit Down TX: " + txid.TXID)

	return err
}

func Leave() error { /// leave holdero table
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	checkoutId := StringToInt(Display.PlayerId)
	singleNameClear(checkoutId)
	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "PlayerLeave"}
	arg2 := rpc.Argument{Name: "id", DataType: "U", Value: checkoutId}
	args := rpc.Arguments{arg1, arg2}
	txid := rpc.Transfer_Result{}
	params := &rpc.SC_Invoke_Params{
		SC_ID:    Round.Contract,
		SC_RPC:   args,
		Ringsize: 2,
	}

	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Leave TX:", txid)
	addLog("Leave Down TX: " + txid.TXID)

	return err
}

func SetTable(seats int, bb, sb, ante uint64, chips, name, av string) error { /// set holdero
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	var hx string
	if av != "" {
		with := "_" + name + "_" + av
		hx = hex.EncodeToString([]byte(with))
	} else {
		out := "_" + name
		hx = hex.EncodeToString([]byte(out))
	}

	var args rpc.Arguments
	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "SetTable"}
	arg2 := rpc.Argument{Name: "seats", DataType: "U", Value: seats}
	arg3 := rpc.Argument{Name: "bigBlind", DataType: "U", Value: bb}
	arg4 := rpc.Argument{Name: "smallBlind", DataType: "U", Value: sb}
	arg5 := rpc.Argument{Name: "ante", DataType: "U", Value: ante}
	arg6 := rpc.Argument{Name: "address", DataType: "S", Value: Wallet.idHash + hx}

	if Round.Version < 110 {
		args = rpc.Arguments{arg1, arg2, arg3, arg4, arg5, arg6}
	} else if Round.Version == 110 {
		arg7 := rpc.Argument{Name: "chips", DataType: "S", Value: chips}
		args = rpc.Arguments{arg1, arg2, arg3, arg4, arg5, arg6, arg7}
	}
	txid := rpc.Transfer_Result{}

	t1 := rpc.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []rpc.Transfer{t1}

	fee, _ := GasEstimate(Round.Contract, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     Round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Set Table TX:", txid)
	addLog("Set Table TX: " + txid.TXID)

	return err
}

func GenerateKey() string {
	random, _ := rand.Prime(rand.Reader, 128)
	shasum := sha256.Sum256([]byte(random.String()))
	str := hex.EncodeToString(shasum[:])
	log.Println("Round Key: ", str)
	addLog("Round Key: " + str)

	return str
}

func DealHand() error { /// holdero hand
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	Wallet.ClientKey = GenerateKey()
	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "DealHand"}
	arg2 := rpc.Argument{Name: "pcSeed", DataType: "H", Value: Wallet.ClientKey}
	args := rpc.Arguments{arg1, arg2}
	txid := rpc.Transfer_Result{}

	var amount uint64

	if Round.Pot == 0 {
		amount = Round.Ante + Round.SB

	} else if Round.Pot == Round.SB || Round.Pot == Round.Ante+Round.SB {
		amount = Round.Ante + Round.BB

	} else {
		amount = Round.Ante
	}

	t := []rpc.Transfer{}
	if Round.Asset {
		t1 := rpc.Transfer{
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      500,
			Burn:        0,
		}

		if Round.Tourney {
			t2 := rpc.Transfer{
				SCID:        crypto.HashHexToHash(TourneySCID),
				Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
				Burn:        amount,
			}
			t = append(t, t1, t2)
		} else {
			t2 := rpc.Transfer{
				SCID:        crypto.HashHexToHash(dReamsSCID),
				Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
				Burn:        amount,
			}
			t = append(t, t1, t2)
		}
	} else {
		t1 := rpc.Transfer{
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      500,
			Burn:        amount,
		}
		t = append(t, t1)
	}

	fee, _ := GasEstimate(Round.Contract, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     Round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	Display.Res = ""
	log.Println("Deal TX:", txid)
	addLog("Deal TX: " + txid.TXID)

	return err
}

func Bet(amt string) error { /// holdero bet
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Bet"}
	args := rpc.Arguments{arg1}
	txid := rpc.Transfer_Result{}

	var t1 rpc.Transfer
	if Round.Asset {
		if Round.Tourney {
			t1 = rpc.Transfer{
				SCID:        crypto.HashHexToHash(TourneySCID),
				Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
				Burn:        ToAtomicOne(amt),
			}
		} else {
			t1 = rpc.Transfer{
				SCID:        crypto.HashHexToHash(dReamsSCID),
				Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
				Burn:        ToAtomicOne(amt),
			}
		}
	} else {
		t1 = rpc.Transfer{
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      0,
			Burn:        ToAtomicOne(amt),
		}
	}

	t := []rpc.Transfer{t1}
	fee, _ := GasEstimate(Round.Contract, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     Round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	Display.Res = ""
	Signal.PlacedBet = true
	log.Println("Bet TX:", txid)
	addLog("Bet TX: " + txid.TXID)

	return err
}

func Check() error { /// holdero check and fold
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Bet"}
	args := rpc.Arguments{arg1}
	txid := rpc.Transfer_Result{}

	var t1 rpc.Transfer
	if !Round.Asset {
		t1 = rpc.Transfer{
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      0,
			Burn:        0,
		}
	} else {
		if Round.Tourney {
			t1 = rpc.Transfer{
				SCID:        crypto.HashHexToHash(TourneySCID),
				Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
				Burn:        0,
			}
		} else {
			t1 = rpc.Transfer{
				SCID:        crypto.HashHexToHash(dReamsSCID),
				Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
				Burn:        0,
			}
		}
	}

	t := []rpc.Transfer{t1}
	fee, _ := GasEstimate(Round.Contract, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     Round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	Display.Res = ""
	log.Println("Check/Fold TX:", txid)
	addLog("Check/Fold TX: " + txid.TXID)

	return err
}

func PayOut(w string) error { /// holdero single winner
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Winner"}
	arg2 := rpc.Argument{Name: "whoWon", DataType: "S", Value: w}
	args := rpc.Arguments{arg1, arg2}
	txid := rpc.Transfer_Result{}

	t1 := rpc.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []rpc.Transfer{t1}
	fee, _ := GasEstimate(Round.Contract, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     Round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Winner TX:", txid)
	addLog("Winnner TX: " + txid.TXID)

	return err
}

func PayoutSplit(r ranker, f1, f2, f3, f4, f5, f6 bool) error { /// holdero split winners
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	ways := 0
	splitWinners := [6]string{"Zero", "Zero", "Zero", "Zero", "Zero", "Zero"}

	if r.p1HighCardArr[0] > 0 && !f1 {
		ways = 1
		splitWinners[0] = "Player1"
	}

	if r.p2HighCardArr[0] > 0 && !f2 {
		ways++
		splitWinners[1] = "Player2"
	}

	if r.p3HighCardArr[0] > 0 && !f3 {
		ways++
		splitWinners[2] = "Player3"
	}

	if r.p4HighCardArr[0] > 0 && !f4 {
		ways++
		splitWinners[3] = "Player4"
	}

	if r.p5HighCardArr[0] > 0 && !f5 {
		ways++
		splitWinners[4] = "Player5"
	}

	if r.p6HighCardArr[0] > 0 && !f6 {
		ways++
		splitWinners[5] = "Player6"
	}

	sort.Strings(splitWinners[:])

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "SplitWinner"}
	arg2 := rpc.Argument{Name: "div", DataType: "U", Value: ways}
	arg3 := rpc.Argument{Name: "split1", DataType: "S", Value: splitWinners[0]}
	arg4 := rpc.Argument{Name: "split2", DataType: "S", Value: splitWinners[1]}
	arg5 := rpc.Argument{Name: "split3", DataType: "S", Value: splitWinners[2]}
	arg6 := rpc.Argument{Name: "split4", DataType: "S", Value: splitWinners[3]}
	arg7 := rpc.Argument{Name: "split5", DataType: "S", Value: splitWinners[4]}
	arg8 := rpc.Argument{Name: "split6", DataType: "S", Value: splitWinners[5]}

	args := rpc.Arguments{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8}
	txid := rpc.Transfer_Result{}

	t1 := rpc.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []rpc.Transfer{t1}
	fee, _ := GasEstimate(Round.Contract, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     Round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Split Winner TX:", txid)
	addLog("Split Winner TX: " + txid.TXID)

	return err
}

func RevealKey(key string) error { /// holdero reveal
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "RevealKey"}
	arg2 := rpc.Argument{Name: "pcSeed", DataType: "H", Value: key}
	args := rpc.Arguments{arg1, arg2}
	txid := rpc.Transfer_Result{}

	params := &rpc.SC_Invoke_Params{
		SC_ID:    Round.Contract,
		SC_RPC:   args,
		Ringsize: 2,
	}

	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	Display.Res = ""
	log.Println("Reveal TX:", txid)
	addLog("Reveal TX: " + txid.TXID)

	return err
}

func CleanTable(amt uint64) error { /// shuffle and clean holdero
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "CleanTable"}
	arg2 := rpc.Argument{Name: "amount", DataType: "U", Value: amt}
	args := rpc.Arguments{arg1, arg2}
	txid := rpc.Transfer_Result{}

	t1 := rpc.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []rpc.Transfer{t1}
	fee, _ := GasEstimate(Round.Contract, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     Round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Clean Table TX:", txid)
	addLog("Clean Table TX: " + txid.TXID)

	return err
}

func TimeOut() error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "TimeOut"}
	args := rpc.Arguments{arg1}
	txid := rpc.Transfer_Result{}

	t1 := rpc.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []rpc.Transfer{t1}
	fee, _ := GasEstimate(Round.Contract, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     Round.Contract,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Timeout TX:", txid)
	addLog("Timeout TX: " + txid.TXID)

	return err
}

func ForceStat() error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "ForceStart"}
	args := rpc.Arguments{arg1}
	txid := rpc.Transfer_Result{}
	params := &rpc.SC_Invoke_Params{
		SC_ID:    Round.Contract,
		SC_RPC:   args,
		Ringsize: 2,
	}

	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Force Start TX:", txid)
	addLog("Force Start TX: " + txid.TXID)

	return err
}

type CardSpecs struct {
	Faces struct {
		Name string `json:"Name"`
		Url  string `json:"Url"`
	} `json:"Faces"`
	Backs struct {
		Name string `json:"Name"`
		Url  string `json:"Url"`
	} `json:"Backs"`
}

type TableSpecs struct {
	MaxBet float64 `json:"Maxbet"`
	MinBuy float64 `json:"Minbuy"`
	MaxBuy float64 `json:"Maxbuy"`
	Time   int     `json:"Time"`
}

func SharedDeckUrl(face, faceUrl, back, backUrl string) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	var cards string
	if face == "" || back == "" {
		cards = "nil"
	} else {
		cards = `{"Faces":{"Name":"` + face + `", "Url":"` + faceUrl + `"},"Backs":{"Name":"` + back + `", "Url":"` + backUrl + `"}}`
	}

	specs := "nil"
	// specs := `{"MaxBet":10,"MinBuy":10,"MaxBuy":20, "Time":120}`

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Deck"}
	arg2 := rpc.Argument{Name: "face", DataType: "S", Value: cards}
	arg3 := rpc.Argument{Name: "back", DataType: "S", Value: specs}
	args := rpc.Arguments{arg1, arg2, arg3}
	txid := rpc.Transfer_Result{}

	params := &rpc.SC_Invoke_Params{
		SC_ID:    Round.Contract,
		SC_RPC:   args,
		Ringsize: 2,
	}

	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Shared TX:", txid)
	addLog("Shared TX: " + txid.TXID)

	return err
}

func GetdReams(amt uint64) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "IssueChips"}
	args := rpc.Arguments{arg1}
	txid := rpc.Transfer_Result{}

	params := &rpc.SC_Invoke_Params{
		SC_ID:           BaccSCID,
		SC_RPC:          args,
		SC_DERO_Deposit: amt,
		Ringsize:        2,
	}

	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Get dReams", txid)
	addLog("Get dReams " + txid.TXID)

	return err
}

func TradedReams(amt uint64) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "ConvertChips"}
	args := rpc.Arguments{arg1}
	txid := rpc.Transfer_Result{}

	scid := crypto.HashHexToHash(dReamsSCID)
	t1 := rpc.Transfer{
		SCID:        scid,
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        amt,
	}

	t := []rpc.Transfer{t1}
	fee, _ := GasEstimate(BaccSCID, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     BaccSCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Trade dReams TX:", txid)
	addLog("Trade dReams TX: " + txid.TXID)

	return err
}

func ownerT3(o bool) (t *rpc.Transfer) {
	if o {
		t = &rpc.Transfer{
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      0,
		}
	} else {
		t = &rpc.Transfer{
			Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
			Amount:      300000,
		}
	}

	return
}

func UploadHolderoContract(d, w bool, pub int) error {
	if d && w {
		rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
		defer cancel()

		code, code_err := GetHoldero110Code(d, pub)
		if code_err != nil {
			log.Println(code_err)
			return nil
		}

		args := rpc.Arguments{}
		txid := rpc.Transfer_Result{}

		params := &rpc.Transfer_Params{
			Transfers: []rpc.Transfer{*ownerT3(Wallet.PokerOwner)},
			SC_Code:   code,
			SC_Value:  0,
			SC_RPC:    args,
			Ringsize:  2,
			Fees:      30000,
		}

		err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
		if err != nil {
			log.Println(err)
			return nil
		}

		log.Println("Upload", txid)
		addLog("Upload " + txid.TXID)

		return err
	}

	return nil
}

func BaccBet(amt, w string) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "PlayBaccarat"}
	arg2 := rpc.Argument{Name: "betOn", DataType: "S", Value: w}
	args := rpc.Arguments{arg1, arg2}
	txid := rpc.Transfer_Result{}

	scid := crypto.HashHexToHash(dReamsSCID)
	t1 := rpc.Transfer{
		SCID:        scid,
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        ToAtomicOne(amt),
	}

	t := []rpc.Transfer{t1}
	fee, _ := GasEstimate(BaccSCID, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     BaccSCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	Bacc.Last = txid.TXID
	Bacc.Notified = false
	if w == "player" {
		log.Println("Baccarat Player TX:", txid)
		addLog("Baccarat Player TX: " + txid.TXID)
	} else if w == "banker" {
		log.Println("Baccarat Banker TX:", txid)
		addLog("Baccarat Banker TX: " + txid.TXID)
	} else {
		log.Println("Baccarat Tie TX:", txid)
		addLog("Baccarat Tie TX: " + txid.TXID)
	}

	Bacc.CHeight = StringToInt(Wallet.Height)

	return err
}

func PredictHigher(scid, name string) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	a := uint64(Predict.Amount)

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Predict"}
	arg2 := rpc.Argument{Name: "pre", DataType: "U", Value: 1}
	arg3 := rpc.Argument{Name: "name", DataType: "S", Value: name}
	args := rpc.Arguments{arg1, arg2, arg3}
	txid := rpc.Transfer_Result{}

	p := &rpc.SC_Invoke_Params{
		SC_ID:           scid,
		SC_RPC:          args,
		SC_DERO_Deposit: a,
		Ringsize:        2,
	}

	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", p)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Prediction TX:", txid)
	addLog("Prediction TX: " + txid.TXID)

	return err
}

func PredictLower(scid, name string) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	a := uint64(Predict.Amount)

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Predict"}
	arg2 := rpc.Argument{Name: "pre", DataType: "U", Value: 0}
	arg3 := rpc.Argument{Name: "name", DataType: "S", Value: name}
	args := rpc.Arguments{arg1, arg2, arg3}
	txid := rpc.Transfer_Result{}

	p := &rpc.SC_Invoke_Params{
		SC_ID:           scid,
		SC_RPC:          args,
		SC_DERO_Deposit: a,
		Ringsize:        2,
	}

	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", p)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Prediction TX:", txid)
	addLog("Prediction TX: " + txid.TXID)

	return err
}

func NameChange(scid, name string) error { /// change leaderboard name
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "NameChange"}
	arg2 := rpc.Argument{Name: "name", DataType: "S", Value: name}
	args := rpc.Arguments{arg1, arg2}
	txid := rpc.Transfer_Result{}

	p := &rpc.SC_Invoke_Params{
		SC_ID:           scid,
		SC_RPC:          args,
		SC_DERO_Deposit: 10000,
		Ringsize:        2,
	}

	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", p)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Name Change TX:", txid)
	addLog("Name Change TX: " + txid.TXID)

	return err
}

func RemoveAddress(scid, name string) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Remove"}
	arg2 := rpc.Argument{Name: "name", DataType: "S", Value: name}
	args := rpc.Arguments{arg1, arg2}
	txid := rpc.Transfer_Result{}

	p := &rpc.SC_Invoke_Params{
		SC_ID:           scid,
		SC_RPC:          args,
		SC_DERO_Deposit: 10000,
		Ringsize:        2,
	}

	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", p)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Remove TX:", txid)
	addLog("Remove TX: " + txid.TXID)

	return err
}

func PickTeam(scid, multi, n string, a uint64, pick int) error { /// pick sports team
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	var amt uint64

	switch multi {
	case "1x":
		amt = a
	case "3x":
		amt = a * 3
	case "5x":
		amt = a * 5
	default:
		amt = a
	}

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Book"}
	arg2 := rpc.Argument{Name: "n", DataType: "S", Value: n}
	arg3 := rpc.Argument{Name: "pre", DataType: "U", Value: pick}
	args := rpc.Arguments{arg1, arg2, arg3}
	txid := rpc.Transfer_Result{}

	t1 := rpc.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        amt,
	}

	t := []rpc.Transfer{t1}
	fee, _ := GasEstimate(scid, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Pick TX:", txid)
	addLog("Pick TX: " + txid.TXID)

	return err
}

func SetSports(end int, amt, dep uint64, scid, league, game, feed string) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "S_start"}
	arg2 := rpc.Argument{Name: "end", DataType: "U", Value: end}
	arg3 := rpc.Argument{Name: "amt", DataType: "U", Value: amt}
	arg4 := rpc.Argument{Name: "league", DataType: "S", Value: league}
	arg5 := rpc.Argument{Name: "game", DataType: "S", Value: game}
	arg6 := rpc.Argument{Name: "feed", DataType: "S", Value: feed}
	args := rpc.Arguments{arg1, arg2, arg3, arg4, arg5, arg6}
	txid := rpc.Transfer_Result{}

	params := &rpc.SC_Invoke_Params{
		SC_ID:           scid,
		SC_RPC:          args,
		SC_DERO_Deposit: dep,
		Ringsize:        2,
	}

	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Set Sports TX:", txid)
	addLog("Set Sports TX: " + txid.TXID)

	return err
}

func SetPrediction(end int, amt, dep uint64, scid, predict, feed string) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "P_start"}
	arg2 := rpc.Argument{Name: "end", DataType: "U", Value: end}
	arg3 := rpc.Argument{Name: "amt", DataType: "U", Value: amt}
	arg4 := rpc.Argument{Name: "predict", DataType: "S", Value: predict}
	arg5 := rpc.Argument{Name: "feed", DataType: "S", Value: feed}
	args := rpc.Arguments{arg1, arg2, arg3, arg4, arg5}
	txid := rpc.Transfer_Result{}

	params := &rpc.SC_Invoke_Params{
		SC_ID:           scid,
		SC_RPC:          args,
		SC_DERO_Deposit: dep,
		Ringsize:        2,
	}

	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Set Prediction TX:", txid)
	addLog("Set Prediction TX: " + txid.TXID)

	return err
}

func PostPrediction(scid string, price int) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Post"}
	arg2 := rpc.Argument{Name: "price", DataType: "U", Value: price}
	args := rpc.Arguments{arg1, arg2}
	txid := rpc.Transfer_Result{}

	params := &rpc.SC_Invoke_Params{
		SC_ID:           scid,
		SC_RPC:          args,
		SC_DERO_Deposit: 0,
		Ringsize:        2,
	}

	err := rpcClientW.CallFor(ctx, &txid, "scinvoke", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Post TX:", txid)
	addLog("Post TX: " + txid.TXID)

	return err
}

func EndSports(scid, num, team string) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "S_end"}
	arg2 := rpc.Argument{Name: "n", DataType: "S", Value: num}
	arg3 := rpc.Argument{Name: "team", DataType: "S", Value: team}
	args := rpc.Arguments{arg1, arg2, arg3}
	txid := rpc.Transfer_Result{}

	t := []rpc.Transfer{}
	fee, _ := GasEstimate(scid, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("End Sports TX:", txid)
	addLog("End Sports TX: " + txid.TXID)

	return err
}

func EndPredition(scid string, price int) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "P_end"}
	arg2 := rpc.Argument{Name: "price", DataType: "U", Value: price}
	args := rpc.Arguments{arg1, arg2, arg2}
	txid := rpc.Transfer_Result{}

	t := []rpc.Transfer{}
	fee, _ := GasEstimate(scid, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_Value:  0,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("End Predition TX:", txid)
	addLog("End Prediction TX: " + txid.TXID)

	return err
}

func UploadBetContract(d, w, c bool, pub int) error {
	if d && w {
		rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
		defer cancel()

		var code string
		var code_err error

		if c {
			code, code_err = GetPredictCode(d, pub)
			if code_err != nil {
				log.Println(code_err)
				return nil
			}
		} else {
			code, code_err = GetSportsCode(d, pub)
			if code_err != nil {
				log.Println(code_err)
				return nil
			}
		}

		args := rpc.Arguments{}
		txid := rpc.Transfer_Result{}

		params := &rpc.Transfer_Params{
			Transfers: []rpc.Transfer{*ownerT3(Wallet.BetOwner)},
			SC_Code:   code,
			SC_Value:  0,
			SC_RPC:    args,
			Ringsize:  2,
			Fees:      11000,
		}

		err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
		if err != nil {
			log.Println(err)
			return nil
		}

		log.Println("Upload", txid)
		addLog("Upload " + txid.TXID)

		return err
	}

	return nil
}

func SetHeaders(name, desc, icon, scid string) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "SetSCIDHeaders"}
	arg2 := rpc.Argument{Name: "name", DataType: "S", Value: name}
	arg3 := rpc.Argument{Name: "descr", DataType: "S", Value: desc}
	arg4 := rpc.Argument{Name: "icon", DataType: "S", Value: icon}
	arg5 := rpc.Argument{Name: "scid", DataType: "S", Value: scid}
	args := rpc.Arguments{arg1, arg2, arg3, arg4, arg5}
	txid := rpc.Transfer_Result{}

	t1 := rpc.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        200,
	}

	t := []rpc.Transfer{t1}
	fee, _ := GasEstimate(GnomonSCID, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_Value:  0,
		SC_ID:     GnomonSCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}
	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("Set Headers TX:", txid)
	addLog("Set Headers TX: " + txid.TXID)

	return err
}

func ClaimNfa(scid string) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "ClaimOwnership"}
	args := rpc.Arguments{arg1, arg1}
	txid := rpc.Transfer_Result{}

	nfa_sc := crypto.HashHexToHash(scid)
	t1 := rpc.Transfer{
		SCID:        nfa_sc,
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        1,
	}

	t := []rpc.Transfer{t1}
	fee, _ := GasEstimate(scid, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Println("Claim TX:", txid)
	addLog("Claim TX: " + txid.TXID)

	return err
}

func NfaBidBuy(scid, bidor string, amt uint64) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: bidor}
	args := rpc.Arguments{arg1}
	txid := rpc.Transfer_Result{}

	t1 := rpc.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        amt,
	}

	t := []rpc.Transfer{t1}
	fee, _ := GasEstimate(scid, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	if bidor == "Bid" {
		log.Println("Bid TX:", txid)
		addLog("Bid TX: " + txid.TXID)
	} else {
		log.Println("Buy TX:", txid)
		addLog("Buy TX: " + txid.TXID)
	}

	return err
}

func NfaSetListing(scid, list, char string, dur, amt, perc uint64) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: "Start"}
	arg2 := rpc.Argument{Name: "listType", DataType: "S", Value: strings.ToLower(list)}
	arg3 := rpc.Argument{Name: "duration", DataType: "U", Value: dur}
	arg4 := rpc.Argument{Name: "startPrice", DataType: "U", Value: amt}
	arg5 := rpc.Argument{Name: "charityDonateAddr", DataType: "S", Value: char}
	arg6 := rpc.Argument{Name: "charityDonatePerc", DataType: "U", Value: perc}
	args := rpc.Arguments{arg1, arg2, arg3, arg4, arg5, arg6}
	txid := rpc.Transfer_Result{}

	asset_scid := crypto.HashHexToHash(scid)
	t1 := rpc.Transfer{
		SCID:        asset_scid,
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        1,
	}

	/// dReams
	t2 := rpc.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      5000,
		Burn:        0,
	}

	/// artificer
	t3 := rpc.Transfer{
		Destination: "dero1qy0khp9s9yw2h0eu20xmy9lth3zp5cacmx3rwt6k45l568d2mmcf6qgcsevzx",
		Amount:      5000,
		Burn:        0,
	}

	t := []rpc.Transfer{t1, t2, t3}
	fee, _ := GasEstimate(scid, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}
	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("NFA List TX:", txid)
	addLog("NFA List TX: " + txid.TXID)

	return err
}

func NfaCancelClose(scid, c string) error {
	rpcClientW, ctx, cancel := SetWalletClient(Wallet.Rpc, Wallet.UserPass)
	defer cancel()

	arg1 := rpc.Argument{Name: "entrypoint", DataType: "S", Value: c}
	args := rpc.Arguments{arg1}
	txid := rpc.Transfer_Result{}

	t1 := rpc.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        0,
	}

	t := []rpc.Transfer{t1}
	fee, _ := GasEstimate(scid, args, t)
	params := &rpc.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	err := rpcClientW.CallFor(ctx, &txid, "transfer", params)
	if err != nil {
		log.Println(err)
		return nil
	}

	if c == "CloseListing" {
		log.Println("Close Listing TX:", txid)
		addLog("Close Listing TX: " + txid.TXID)
	} else {
		log.Println("Cancel Listing TX:", txid)
		addLog("Cancel Listing TX: " + txid.TXID)
	}

	return err
}
