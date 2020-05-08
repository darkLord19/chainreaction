package helpers

import (
	"log"

	"github.com/chainreaction/models"
)

// NotifyIndividual notifies individual player
func NotifyIndividual(p *models.Player, val interface{}) {
	err := p.WriteToWebsocket(val)
	if err != nil {
		log.Printf("error: %v", err)
		p.CleanupWs()
	}
}
