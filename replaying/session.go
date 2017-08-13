package replaying

import (
	"github.com/v2pro/koala/st"
	"time"
	"unsafe"
	"strings"
)

type ReplayingSession struct {
	st.Session `json:"-"`
	ReplayedOutboundTalkCollector chan ReplayedTalk `json:"-"`
	ReplayedRequestTime           int64
	ReplayedResponse              []byte
	ReplayedResponseTime          int64
	ReplayedOutboundTalks         []ReplayedTalk
}

func (replayingSession *ReplayingSession) Finish(response []byte) {
	replayingSession.ReplayedResponse = response
	replayingSession.ReplayedResponseTime = time.Now().UnixNano()
	done := false
	for !done {
		select {
		case replayedTalk := <- replayingSession.ReplayedOutboundTalkCollector:
			replayingSession.ReplayedOutboundTalks = append(replayingSession.ReplayedOutboundTalks, replayedTalk)
		default:
			done = true
		}
	}
}

func (replayingSession *ReplayingSession) MatchOutboundTalk(outboundRequest []byte) *st.Talk {
	unit := 16
	chunks := cutToChunks(outboundRequest, unit)
	keys := replayingSession.loadKeys()
	scores := make([]int, len(replayingSession.OutboundTalks))
	maxScore := 0
	maxScoreIndex := 0
	for _, chunk := range chunks {
		for j, key := range keys {
			if len(key) < len(chunk) {
				continue
			}
			keyAsString := *(*string)(unsafe.Pointer(&key))
			chunkAsString := *(*string)(unsafe.Pointer(&chunk))
			pos := strings.Index(keyAsString, chunkAsString)
			if pos >= 0 {
				keys[j] = key[pos:]
				scores[j]++
				if scores[j] > maxScore {
					maxScore = scores[j]
					maxScoreIndex = j
				}
			}
		}
	}
	if maxScore == 0 {
		return nil
	}
	return replayingSession.OutboundTalks[maxScoreIndex]

}

func (replayingSession *ReplayingSession) loadKeys() [][]byte {
	keys := make([][]byte, len(replayingSession.OutboundTalks))
	for i, entry := range replayingSession.OutboundTalks {
		keys[i] = entry.Request
	}
	return keys
}


func cutToChunks(key []byte, unit int) [][]byte {
	chunks := [][]byte{}
	chunkCount := len(key) / unit
	for i := 0; i < len(key) / unit; i++ {
		chunks = append(chunks, key[i * unit:(i + 1) * unit])
	}
	lastChunk := key[chunkCount * unit:]
	if len(lastChunk) > 0 {
		chunks = append(chunks, lastChunk)
	}
	return chunks
}