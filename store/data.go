package store

import (
	"fmt"
	"github.com/3hajk/ethereum-cli/models"
	"sort"
)

type TransactionMap map[string]*models.Transaction
type TransactionDataIndex map[uint64][]string
type TransactionAddressIndex map[string][]*models.Transaction
type TransactionAddressCount map[string]int
type TransactionData struct {
	ready      bool
	processing bool
	wc         chan *models.Transaction
	Done       chan struct{}
	tm         TransactionMap
	di         TransactionDataIndex
	ai         TransactionAddressIndex
	ac         TransactionAddressCount
}

func InitData() *TransactionData {
	tm := make(TransactionMap)
	di := make(TransactionDataIndex)
	ai := make(TransactionAddressIndex)
	ac := make(TransactionAddressCount)
	w := make(chan *models.Transaction)
	d := make(chan struct{})
	td := &TransactionData{
		false,
		false,
		w,
		d,
		tm,
		di,
		ai,
		ac,
	}
	go td.run()
	return td
}

func (td *TransactionData) run() {
	for {
		select {
		case tx := <-td.wc:
			td.tm[tx.Hash] = tx
			td.ai[tx.To] = append(td.ai[tx.To], tx)
			td.ac[tx.To] += 1
		case <-td.Done:
			return
		}
	}
}

func (td *TransactionData) Add(tx *models.Transaction) {
	td.wc <- tx
}

func (td *TransactionData) Set(time uint64, tx *models.Transaction) {
	td.tm[tx.Hash] = tx
	//td.di[time]=append(td.di[time],tx.Hash)
	td.ai[tx.To] = append(td.ai[tx.To], tx)
	td.ac[tx.To] += 1
}

func (td *TransactionData) SetProcessing() {
	td.processing = true
}

func (td *TransactionData) IsProcessing() bool {
	return td.processing
}

func (td *TransactionData) SetReady() {
	td.ready = true
	td.processing = false
}

func (td *TransactionData) DataIsReady() bool {
	return td.ready
}

func (td *TransactionData) GetTmByAddress(address string) ([]*models.Transaction, error) {
	txs, ok := td.ai[address]
	if ok != true {
		return nil, fmt.Errorf("tx by %s not found", address)
	}
	return txs, nil
}

func (td *TransactionData) GetTopAddress(top int) map[string]int {
	r := make(map[string]int, 0)
	if !td.ready {
		return r
	}
	for _, res := range sortedKeys(td.ac)[:top] {
		r[res] = td.ac[res]
	}
	return r
}

func (td *TransactionData) GetTm() *TransactionMap {
	return &td.tm
}

func (td *TransactionData) GetDi() *TransactionDataIndex {
	return &td.di
}

func (td *TransactionData) GetAi() *TransactionAddressIndex {
	return &td.ai
}

func (td *TransactionData) Count() int {
	return len(td.tm)
}

func (td *TransactionData) Clean() {
	td.ready = false
	td.processing = false
	for k := range td.tm {
		delete(td.tm, k)
	}
	for k := range td.di {
		delete(td.di, k)
	}
	for k := range td.ai {
		delete(td.ai, k)
		delete(td.ac, k)
	}
}
func (td *TransactionData) GetAc() *TransactionAddressCount {
	return &td.ac
}

type sortedMap struct {
	m map[string]int
	s []string
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[string]int) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}
