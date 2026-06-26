package modelschema

import "encoding/xml"

// mo phong cau truc the <show-config> của NETCONF
type GetConfigRequest struct {
	XMLName xml.Name `xml:"urn.ietf:params:xml:ns:netconf:base:1.0 get-config"`
	Source  Source   `xml:"source"`
	Filter  Filter   `xml:"filter"`
}

type Source struct {
	Running struct{} `xml:"running"`
}

type Filter struct {
	Type  string `xml:"type,attr"` //attr type="subtree"
	Inner string `xml:",innerxml"` //Inject RawXML(dataPayload) without encoding
}

// RPCReply struct contains all content from tag <rpc-reply>
type RPCReply struct {
	InnerXML string `xml:",innerxml"`
}

// --- EDITING FUNCTION ---
// mo phong cau truc the <edit-config> của NETCONF
type EditConfigRequest struct {
	XMLName xml.Name `xml:"urn.ietf:params:xml:ns:netconf:base:1.0 edit-config"`
	Target  Target   `xml:"target"`
	Config  Config   `xml:"config"`
}

type Target struct {
	Candidate struct{} `xml:"candidate"` //write drafts before commit
}

type Config struct {
	Inner string `xml:",innerxml"`
}

type CommitRequest struct {
	XMLName xml.Name `xml:"urn.ietf:params:xml:ns:netconf:base:1.0 commit"`
}

// //Marshal method interface for netconf.Exec
// type RawXML string
// func (r RawXML) MarshalMethod() string {
// 	return string(r)
// }
