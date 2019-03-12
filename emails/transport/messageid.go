package transport

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	hostname    string
	getHostname sync.Once
)

// makeMsgID creates a new, globally unique message ID, useable as
// a Message-ID as per RFC822/RFC2822.
func makeMsgID() string {
	getHostname.Do(func() {
		var err error
		if hostname, err = os.Hostname(); err != nil {
			logrus.Infof("ERROR get hostname: %v", err)
			hostname = "localhost"
		}
	})
	now := time.Now()
	return fmt.Sprintf("<%d.%d.%d@%s>", now.Unix(), now.UnixNano(), rand.Int63(), hostname)
}
