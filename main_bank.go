package main

import (
	"fmt"
	"sync"
)

type Account struct {
	balance int
	mu      sync.Mutex
}

func (a *Account) Deposit(amount int, c chan bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += amount
	c <- true
}

func (a *Account) Withdraw(amount int, c chan bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.balance >= amount {
		a.balance -= amount
		c <- true
	} else {
		c <- false
	}
}

func (a *Account) Transfer(amount int, to *Account, c chan bool) {
	a.Withdraw(amount, c)
	to.Deposit(amount, c)
}

func main() {
	// initialize two accounts with a balance of 1000
	account1 := &Account{balance: 1000}
	account2 := &Account{balance: 1000}

	// create a channel to receive transfer results
	results := make(chan bool)

	// launch 1 million goroutines to transfer funds between the accounts
	for i := 0; i < 2000000; i++ {
		// create a separate channel for each transfer
		c := make(chan bool)

		// launch a goroutine to perform the transfer
		go func() {
			// transfer $100 from account 1 to account 2
			account1.Transfer(100, account2, c)
		}()

		// wait for the transfer to complete and collect its result
		go func() {
			<-c
			results <- true
		}()

		// print a status update every 1000 transfers
		if i%1000 == 0 {
			fmt.Printf("Completed %d transfers\n", i)
		}
	}

	// wait for all transfers to complete and collect the results
	successCount := 0
	failureCount := 0
	for i := 0; i < 2000000; i++ {
		if <-results {
			successCount++
		} else {
			failureCount++
		}
	}

	// print the final balances and transfer stats
	fmt.Printf("Final account 1 balance: $%d\n", account1.balance)
	fmt.Printf("Final account 2 balance: $%d\n", account2.balance)
	fmt.Printf("Successful transfers: %d\n", successCount)
	fmt.Printf("Failed transfers: %d\n", failureCount)
}
