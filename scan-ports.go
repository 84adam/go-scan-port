// Based on code from Black Hat Go: Go Programming for Hackers and Pentesters
// by Tom Steele, Chris Patten, and Dan Kottmann
// ch.2, pp. 30-31: Multichannel Communication

// Scans: `scanme.nmap.org`
// WARNING: Only scan servers under your own control, or which you have permission to scan

package main

import (
    "fmt"
    "net"
    "sort"
)

// for a given range of ports, check whether each is open or closed
// if closed, write a 0 into the result channel
// if open, write the port number into the results channel
func worker(ports, results chan int) {
    for p := range ports {
        address := fmt.Sprintf("scanme.nmap.org:%d", p)
        conn, err := net.Dial("tcp", address)
        if err != nil {
            fmt.Printf("[CLOSED] port %d\n", p)
            results <- 0
            continue
        }
        conn.Close()
        fmt.Printf("[OPEN] port %d\n", p)
        results <- p
    }
}

func main() {
    // create a buffered channel with a capacity of 100
    ports := make(chan int, 100)
    // create an unbuffered channel in which to store/send/retrieve results
    results := make(chan int)
    // define and initialize a slice of ints in which to store the open ports only
    var openports []int
    
    // for each port in the current set of ports start a `worker` goroutine
    // ... and check whether it is open or closed
    for i := 0; i < cap(ports); i++ {
        go worker(ports, results)
    }
    
    // closure to add ports 1 to 1024 into the queue to be checked
    go func() {
        for i := 1; i <= 1024; i++ {
            ports <- i
        }
    }()
    
    // read from the results channel one at a time
    // ... and if it is open (number present) then add it to the openports int slice
    for i := 0; i < 1024; i++ {
        port := <-results
        if port != 0 {
            openports = append(openports, port)
        }
    }

    // close the ports channel and the results channel as all the corresponding work is done now
    close(ports)
    close(results)
    
    // sort the openports slice
    sort.Ints(openports)
    
    // print out the results
    for _, port := range openports {
        fmt.Printf("[OPEN] port %d\n", port)
    }

    fmt.Println("DONE")
}
