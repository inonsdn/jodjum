// scratch/genvapid.go  — run once with `go run`, then delete
package main

import (
	"fmt"

	webpush "github.com/SherClockHolmes/webpush-go"
)

func main() {
	priv, pub, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		panic(err)
	}
	fmt.Println("VAPID_PRIVATE_KEY=" + priv)
	fmt.Println("VAPID_PUBLIC_KEY=" + pub)
}
