package main

import (
        "fmt"
        "math/rand"
        "time"
)

func main() {
        //strSantas := [...]string{"Anna", "John", "Nico", "Nicole", "Salim", "Susan" }
        strSantas := [...]string{"Anna", "Nina", "Laila", "John", "Nico", "Nicole", "Salim", "Susan", }
        var k int
        fmt.Printf("Arranging %d Santas ...\n\n", len(strSantas))
        l := len(strSantas) - 1

        rand.Seed(time.Now().UnixNano())
        rand.Shuffle(len(strSantas), func(i, j int) { strSantas[i], strSantas[j] = strSantas[j], strSantas[i] })

        // fmt.Println("Dear grown-up gift-exchangers:\n")
        for j := 0; j <= l; j++ {
                if j == l {
                        k = 0
                } else {
                        k = j + 1
                }
        fmt.Printf("Ho! Ho! Ho! %s, you will be making %s especially merry this year by putting a special gift in their stocking. \n", strSantas[j], strSantas[k])
        //fmt.Printf("'Tis the season! %s, you will be making %s especially merry this year with a Christmas gift!\n", strSantas[j], strSantas[k])
        }
} //main
