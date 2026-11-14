package main

import (
        "flag"
        "fmt"
        "math/rand"
        "time"
)

func main() {
        var strSantas []string
        s := flag.Bool("stockings", true, "a bool indicating whether to run for stocking stuffers")
        flag.Parse()
        fmt.Println("stockings:", *s)

        if *s {
                strSantas = []string{"Name1", "Name2"}               
        } else {
                strSantas = []string{"Name1", "Name2"}
        }           
        var k int
        fmt.Printf("Arranging %d Santas ...\n\n", len(strSantas))
        l := len(strSantas) - 1

        rand.Seed(time.Now().UnixNano())
        rand.Shuffle(len(strSantas), func(i, j int) { strSantas[i], strSantas[j] = strSantas[j], strSantas[i] })

        if *s {
                fmt.Println("Dear stocking stuffers:\n")
        } else {    
                fmt.Println("Dear grown-up gift-exchangers:\n")
        }           
        for j := 0; j <= l; j++ {
                if j == l {
                        k = 0 
                } else {    
                        k = j + 1
                }           
                if *s {
                        fmt.Printf("Ho! Ho! Ho! %s, you will be making %s especially merry this year by putting a special gift in their stocking. \n", strSantas[j], strSantas[k])
                } else {
                        fmt.Printf("'Tis the season! %s, you will be making %s especially merry this year with a Christmas gift!\n", strSantas[j], strSantas[k])
                }
        }           
} //main 
