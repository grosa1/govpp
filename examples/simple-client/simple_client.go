// Copyright (c) 2017 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// simple-client is an example VPP management application that exercises the
// govpp API on real-world use-cases.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"git.fd.io/govpp.git"
	"git.fd.io/govpp.git/adapter/socketclient"
	"git.fd.io/govpp.git/api"
	"git.fd.io/govpp.git/core"
	"git.fd.io/govpp.git/examples/binapi/interfaces"
	"git.fd.io/govpp.git/examples/binapi/ip"
	"git.fd.io/govpp.git/examples/binapi/vpe"
)

var (
	sockAddr = flag.String("sock", socketclient.DefaultSocketName, "Path to VPP binary API socket file")
)

func main() {
	flag.Parse()

	fmt.Println("Starting simple client example")

	// connect to VPP asynchronously
	conn, conev, err := govpp.AsyncConnect(*sockAddr, core.DefaultMaxReconnectAttempts, core.DefaultReconnectInterval)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	defer conn.Disconnect()

	// wait for Connected event
	select {
	case e := <-conev:
		if e.State != core.Connected {
			log.Fatalln("ERROR: connecting to VPP failed:", e.Error)
		}
	}

	// create an API channel that will be used in the examples
	ch, err := conn.NewAPIChannel()
	if err != nil {
		log.Fatalln("ERROR: creating channel failed:", err)
	}
	defer ch.Close()

	if err := ch.CheckCompatiblity(vpe.AllMessages()...); err != nil {
		log.Fatal(err)
	}

	vppVersion(ch)

	if err := ch.CheckCompatiblity(interfaces.AllMessages()...); err != nil {
		log.Fatal(err)
	}

	idx := createLoopback(ch)
	interfaceDump(ch)

	addIPAddress(ch, idx)
	ipAddressDump(ch, idx)
	interfaceNotifications(ch, idx)

	if len(Errors) > 0 {
		fmt.Printf("finished with %d errors\n", len(Errors))
		os.Exit(1)
	} else {
		fmt.Println("finished successfully")
	}
}

var Errors []error

func logError(err error, msg string) {
	fmt.Printf("ERROR: %s: %v\n", msg, err)
	Errors = append(Errors, err)
}

// vppVersion is the simplest API example - it retrieves VPP version.
func vppVersion(ch api.Channel) {
	fmt.Println("Retrieving version")

	req := &vpe.ShowVersion{}
	reply := &vpe.ShowVersionReply{}

	if err := ch.SendRequest(req).ReceiveReply(reply); err != nil {
		logError(err, "retrieving version")
		return
	}
	fmt.Printf("reply: %+v\n", reply)

	fmt.Printf("VPP version: %q\n", cleanString(reply.Version))
	fmt.Println("OK")
	fmt.Println()
}

// createLoopback sends request to create loopback interface.
func createLoopback(ch api.Channel) interfaces.InterfaceIndex {
	fmt.Println("Creating loopback interface")

	req := &interfaces.CreateLoopback{}
	reply := &interfaces.CreateLoopbackReply{}

	if err := ch.SendRequest(req).ReceiveReply(reply); err != nil {
		logError(err, "creating loopback interface")
		return 0
	}
	fmt.Printf("reply: %+v\n", reply)

	fmt.Printf("interface index: %v\n", reply.SwIfIndex)
	fmt.Println("OK")
	fmt.Println()

	return reply.SwIfIndex
}

// interfaceDump shows an example of multipart request (multiple replies are expected).
func interfaceDump(ch api.Channel) {
	fmt.Println("Dumping interfaces")

	n := 0
	reqCtx := ch.SendMultiRequest(&interfaces.SwInterfaceDump{})
	for {
		msg := &interfaces.SwInterfaceDetails{}
		stop, err := reqCtx.ReceiveReply(msg)
		if stop {
			break
		}
		if err != nil {
			logError(err, "dumping interfaces")
			return
		}
		n++
		fmt.Printf(" - interface #%d: %+v\n", n, msg)
	}

	fmt.Println("OK")
	fmt.Println()
}

// addIPAddress sends request to add IP address to interface.
func addIPAddress(ch api.Channel, index interfaces.InterfaceIndex) {
	fmt.Printf("Adding IP address to interface to interface index %d\n", index)

	req := &interfaces.SwInterfaceAddDelAddress{
		SwIfIndex: index,
		IsAdd:     true,
		Prefix: interfaces.AddressWithPrefix{
			Address: interfaces.Address{
				Af: interfaces.ADDRESS_IP4,
				Un: interfaces.AddressUnionIP4(interfaces.IP4Address{10, 10, 0, uint8(index)}),
			},
			Len: 32,
		},
	}
	reply := &interfaces.SwInterfaceAddDelAddressReply{}

	if err := ch.SendRequest(req).ReceiveReply(reply); err != nil {
		logError(err, "adding IP address to interface")
		return
	}
	fmt.Printf("reply: %+v\n", reply)

	fmt.Println("OK")
	fmt.Println()
}

func ipAddressDump(ch api.Channel, index interfaces.InterfaceIndex) {
	fmt.Printf("Dumping IP addresses for interface index %d\n", index)

	req := &ip.IPAddressDump{
		SwIfIndex: ip.InterfaceIndex(index),
	}
	reqCtx := ch.SendMultiRequest(req)

	for {
		msg := &ip.IPAddressDetails{}
		stop, err := reqCtx.ReceiveReply(msg)
		if err != nil {
			logError(err, "dumping IP addresses")
			return
		}
		if stop {
			break
		}
		fmt.Printf(" - ip address: %+v\n", msg)
	}

	fmt.Println("OK")
	fmt.Println()
}

// interfaceNotifications shows the usage of notification API. Note that for notifications,
// you are supposed to create your own Go channel with your preferred buffer size. If the channel's
// buffer is full, the notifications will not be delivered into it.
func interfaceNotifications(ch api.Channel, index interfaces.InterfaceIndex) {
	fmt.Printf("Subscribing to notificaiton events for interface index %d\n", index)

	notifChan := make(chan api.Message, 100)

	// subscribe for specific notification message
	sub, err := ch.SubscribeNotification(notifChan, &interfaces.SwInterfaceEvent{})
	if err != nil {
		logError(err, "subscribing to interface events")
		return
	}

	// enable interface events in VPP
	err = ch.SendRequest(&interfaces.WantInterfaceEvents{
		PID:           uint32(os.Getpid()),
		EnableDisable: 1,
	}).ReceiveReply(&interfaces.WantInterfaceEventsReply{})
	if err != nil {
		logError(err, "enabling interface events")
		return
	}

	// receive notifications
	go func() {
		for notif := range notifChan {
			fmt.Printf("incoming event: %+v\n", notif.(*interfaces.SwInterfaceEvent))
		}
	}()

	// generate some events in VPP
	err = ch.SendRequest(&interfaces.SwInterfaceSetFlags{
		SwIfIndex: index,
		Flags:     interfaces.IF_STATUS_API_FLAG_ADMIN_UP,
	}).ReceiveReply(&interfaces.SwInterfaceSetFlagsReply{})
	if err != nil {
		logError(err, "setting interface flags")
		return
	}
	err = ch.SendRequest(&interfaces.SwInterfaceSetFlags{
		SwIfIndex: index,
		Flags:     0,
	}).ReceiveReply(&interfaces.SwInterfaceSetFlagsReply{})
	if err != nil {
		logError(err, "setting interface flags")
		return
	}

	// disable interface events in VPP
	err = ch.SendRequest(&interfaces.WantInterfaceEvents{
		PID:           uint32(os.Getpid()),
		EnableDisable: 0,
	}).ReceiveReply(&interfaces.WantInterfaceEventsReply{})
	if err != nil {
		logError(err, "setting interface flags")
		return
	}

	// unsubscribe from delivery of the notifications
	err = sub.Unsubscribe()
	if err != nil {
		logError(err, "unsubscribing from interface events")
		return
	}

	fmt.Println("OK")
	fmt.Println()
}

func cleanString(str string) string {
	return strings.Split(str, "\x00")[0]
}
