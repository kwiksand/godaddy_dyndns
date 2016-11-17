# GoDaddy Dynamic DNS
A butchered verssion of prashantv/godaddy_dyndns to be used as an include library

## Installing

[Install Go](https://golang.org/doc/install) and run:
```
go get -v -u -f github.com/kwiksand/godaddy_dyndns
```

## Usage

To use the updater:
```

import (

    ...

    "github.com/kwiksand/godaddy_dyndns"
)
    var goDaddyDNS = godaddy_dyndns.New()

    var rootDomain = "example.com"
    var subDomain = "testhost"

    goDaddyDNS.SetKey("insert API Key")
    goDaddyDNS.SetSecret("insert API Secret")

    if ip != "" {
		// Use specified IP
        publicIP = *ip
    } else {
        pubIP, err := goDaddyDNS.GetPublicIP()
        publicIP = pubIP
        if err != nil {
            log.Fatalf("GetPublicIP failed: %v", err)
        }
    }

    log.Printf("About to set Host to point at: " + publicIP)

    currentIP, err := goDaddyDNS.GetDNS(*rootDomain, *subDomain)
    if err != nil {
        // we're inserting
        log.Printf("Create DNS record for %v", publicIP)
        if err := goDaddyDNS.InsertDNS(publicIP, *rootDomain, *subDomain); err != nil {
            log.Fatalf("InsertDNS failed: %v", err)
        }
    } else {
        // we're updating

        if currentIP == publicIP {
            log.Printf("Nothing to update (publicIP = DNS = %v)", publicIP)
            return
        }

        log.Printf("Update DNS from %v to %v", currentIP, publicIP)
        if err := goDaddyDNS.UpdateDNS(publicIP, *rootDomain, *subDomain); err != nil {
            log.Fatalf("UpdateDNS failed: %v", err)
        }
    }

    log.Printf("Update successful")
```

